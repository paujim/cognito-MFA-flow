import os
from aws_cdk import (
    core,
    aws_iam as iam,
    aws_cognito as cognito,
    aws_apigateway as apigateway,
    aws_lambda as _lambda,
    aws_cognito as cognito,
    aws_apigateway as apigateway,
)


class CognitoMfaFlowStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

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
        self.pool = pool

        client = pool.add_client(
            id="customer-app-client",
            auth_flows=cognito.AuthFlow(
                user_password=True,
                refresh_token=True),
        )

        self.client = client


class MFAAPIStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, pool: cognito.UserPool, client: cognito.UserPoolClient, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

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
            code=_lambda.Code.from_asset(
                os.path.join("lambda", "main.zip"))
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
        apigateway.LambdaRestApi(
            scope=self,
            id="mfa-api",
            handler=backend,
            endpoint_types=[apigateway.EndpointType.REGIONAL]
        )
