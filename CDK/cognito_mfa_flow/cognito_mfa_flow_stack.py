import os
from aws_cdk import (
    core,
    aws_iam as iam,
    aws_cognito as cognito,
    aws_apigateway as apigateway,
    aws_lambda as _lambda,
    aws_cognito as cognito,
    aws_apigateway as apigateway,
    aws_certificatemanager as acm,
    aws_route53 as route53,
    aws_route53_targets as targets,
    aws_secretsmanager as secretsmanager,
    aws_codebuild as codebuild,
    aws_codedeploy as codedeploy,
    aws_codecommit as codecommit,
    aws_codepipeline as codepipeline,
    aws_codepipeline_actions as codepipeline_actions,
    aws_s3 as s3,
    custom_resources as cr,
    aws_s3_notifications as s3n,
    aws_cloudfront as cloudfront,
    aws_cloudfront_origins as origins,
)

SECRET_GITHUB_ID = "github/oAuthToken"
SECRET_GITHUB_JSON_FIELD = "oAuthToken"


class BuildPipelineStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        artifact_bucket = s3.Bucket(
            scope=self,
            id="s3-artifact",
            versioned=True,
        )
        self.artifact_bucket = artifact_bucket

        oauth_token = core.SecretValue.secrets_manager(
            secret_id=SECRET_GITHUB_ID,
            json_field=SECRET_GITHUB_JSON_FIELD,
        )

        # Codepipeline
        build_pipeline = codepipeline.Pipeline(
            scope=self,
            id="lambda-pipeline",
            restart_execution_on_update=True,
        )

        source_output = codepipeline.Artifact()
        build_pipeline.add_stage(
            stage_name="Source",
            actions=[
                codepipeline_actions.GitHubSourceAction(
                    oauth_token=oauth_token,
                    action_name="GitHub",
                    owner="paujim",
                    repo="cognito-MFA-flow",
                    output=source_output,
                )]
        )

        build_specs = {
            "version": "0.2",
            "env": {
                "variables": {
                    "GO111MODULE": "on",
                }
            },
            "phases": {
                "install": {
                    "commands": [
                        "cd Server",
                        "go get .",
                    ]
                },
                "pre_build": {
                    "commands": [
                        "go test .",  # Run all tests included with our application
                    ]
                },
                "build": {
                    "commands": [
                        "go build -o main",  # Build the go application
                        "zip main.zip main",  # Zip the go application
                        "cd ..",
                        "zip -r Client.zip Client"
                    ]
                }
            },
            "artifacts": {
                # "base-directory": "Server",
                "files": [
                    "Server/main.zip",
                    "Client.zip",
                ],
            }
        }
        build_output = codepipeline.Artifact()
        build_pipeline.add_stage(
            stage_name="Build",
            actions=[codepipeline_actions.CodeBuildAction(
                action_name="CodeBuild",
                project=codebuild.Project(
                    scope=self,
                    id="codebuild-build",
                    build_spec=codebuild.BuildSpec.from_object(build_specs),
                ),
                input=source_output,
                outputs=[build_output]
            )]
        )

        build_pipeline.add_stage(
            stage_name="Upload",
            actions=[
                codepipeline_actions.S3DeployAction(
                    bucket=artifact_bucket,
                    input=build_output,
                    action_name="S3Upload",
                    extract=True,
                    # object_key="Server/main.zip",
                )]
        )


class CognitoMfaFlowStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, artifact_bucket: s3.Bucket, **kwargs) -> None:
        super().__init__(scope, id,  **kwargs)

        pool = cognito.UserPool(
            scope=self,
            id="user-pool",
            mfa=cognito.Mfa.OPTIONAL,
            mfa_second_factor=cognito.MfaSecondFactor(otp=True, sms=True),
            password_policy=cognito.PasswordPolicy(
                min_length=12,
                require_lowercase=True,
                require_uppercase=False,
                require_digits=False,
                require_symbols=False,
            )
        )

        client = pool.add_client(
            id="customer-app-client",
            auth_flows=cognito.AuthFlow(
                user_password=True,
                refresh_token=True),
        )

        backend = _lambda.Function(
            scope=self,
            id="api-function",
            runtime=_lambda.Runtime.GO_1_X,
            handler="main",
            memory_size=500,
            timeout=core.Duration.seconds(10),
            environment={
                "USER_POOL_ID": pool.user_pool_id,
                "CLIENT_ID": client.user_pool_client_id,
            },
            code=_lambda.Code.from_bucket(
                bucket=artifact_bucket,
                key="Server/main.zip",
            ),
        )
        backend.add_to_role_policy(
            statement=iam.PolicyStatement(
                actions=[
                    "cognito-idp:RespondToAuthChallenge",
                    "cognito-idp:InitiateAuth",
                    "cognito-idp:SetUserMFAPreference",
                    "cognito-idp:AssociateSoftwareToken",
                    "cognito-idp:VerifySoftwareToken"
                ],
                resources=[pool.user_pool_arn]))

        api = apigateway.LambdaRestApi(
            scope=self,
            id="mfa-api",
            handler=backend,
            endpoint_types=[apigateway.EndpointType.REGIONAL],
            default_cors_preflight_options=apigateway.CorsOptions(
                allow_origins=["*"])
        )
        self.api = api

        self.backend_fn = backend

        static_website_bucket = s3.Bucket(
            scope=self,
            id="static-website-bucket",
        )
        self.static_website_bucket = static_website_bucket

        distribution = cloudfront.CloudFrontWebDistribution(
            scope=self,
            id="static-website-distribution",
            default_root_object="index.html",
            origin_configs=[
                cloudfront.SourceConfiguration(
                    s3_origin_source=cloudfront.S3OriginConfig(
                        s3_bucket_source=static_website_bucket,
                        origin_access_identity=cloudfront.OriginAccessIdentity(
                            scope=self,
                            id="origin-access-identity",
                        )
                    ),
                    behaviors=[cloudfront.Behavior(is_default_behavior=True)]
                )
            ],
        )


class DeployPipelineStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, artifact_bucket: s3.Bucket, static_website_bucket: s3.Bucket, backend_fn: _lambda.Function, api: apigateway.LambdaRestApi, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        fn = _lambda.Function(
            scope=self,
            id="source-update-function",
            runtime=_lambda.Runtime.PYTHON_3_8,
            handler="index.handler",
            # memory_size=500,
            timeout=core.Duration.seconds(10),
            code=_lambda.Code.from_asset(
                os.path.join("lambdas", "updateSource"))
        )
        fn.add_to_role_policy(
            statement=iam.PolicyStatement(
                actions=[
                    "lambda:UpdateFunctionCode",
                ],
                resources=[
                    backend_fn.function_arn,
                ]))
        fn.add_to_role_policy(
            statement=iam.PolicyStatement(
                actions=[
                    "s3:GetObject",
                ],
                resources=[
                    artifact_bucket.bucket_arn + "/Server/main.zip",
                ]))

        # Codepipeline
        deploy_pipeline = codepipeline.Pipeline(
            scope=self,
            id="deploy-pipeline",
            restart_execution_on_update=True,
        )

        lambda_source_output = codepipeline.Artifact()
        client_source_output = codepipeline.Artifact()
        deploy_pipeline.add_stage(
            stage_name="Source",
            actions=[
                codepipeline_actions.S3SourceAction(
                    action_name="LambdaSource",
                    bucket=artifact_bucket,
                    bucket_key="Server/main.zip",
                    output=lambda_source_output,
                ),
                codepipeline_actions.S3SourceAction(
                    action_name="ClientSource",
                    bucket=artifact_bucket,
                    bucket_key="Client.zip",
                    output=client_source_output,
                )]
        )

        build_specs = {
            "version": "0.2",
            "env": {
                "variables": {
                    "REACT_APP_AUTH_URL": api.url,
                }
            },
            "phases": {
                "install": {
                    "runtime-versions": {
                        "nodejs": "10",
                    },
                    "commands": [
                        "npm install -g yarn",
                    ]
                },
                "build": {
                    "commands": [
                        "cd Client",
                        "npm install",
                        "yarn build",
                    ]
                }
            },
            "artifacts": {
                "base-directory": "Client/build",
                "files": [
                    "**/*",
                ],
            }
        }
        client_build_output = codepipeline.Artifact()
        deploy_pipeline.add_stage(
            stage_name="Build",
            actions=[
                codepipeline_actions.CodeBuildAction(
                    action_name="ClientBuild",
                    project=codebuild.Project(
                        scope=self,
                        id="codebuild-client",
                        build_spec=codebuild.BuildSpec.from_object(
                            build_specs),
                    ),
                    input=client_source_output,
                    outputs=[client_build_output]
                )]
        )

        deploy_pipeline.add_stage(
            stage_name="Deploy",
            actions=[
                codepipeline_actions.LambdaInvokeAction(
                    lambda_=fn,
                    inputs=[lambda_source_output],
                    action_name="UpdateSource",
                    user_parameters={
                        "functionName": backend_fn.function_name,
                        "sourceBucket": artifact_bucket.bucket_name,
                        "sourceKey": "Server/main.zip",
                    }
                ),
                codepipeline_actions.S3DeployAction(
                    bucket=static_website_bucket,
                    input=client_build_output,
                    action_name="DeployClient",
                    extract=True,
                    # object_key="Server/main.zip",
                ),
            ]
        )
