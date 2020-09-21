#!/usr/bin/env python3

from aws_cdk import core

from cognito_mfa_flow.cognito_mfa_flow_stack import(
    BuildPipelineStack,
    CognitoMfaFlowStack,
    DeployPipelineStack,
)


app = core.App()

build_stack = BuildPipelineStack(
    scope=app,
    id="cognito-mfa-build-pipeline"
)
cognito_stack = CognitoMfaFlowStack(
    scope=app,
    id="cognito-mfa-flow",
    artifact_bucket=build_stack.artifact_bucket,
)
deploy_stack = DeployPipelineStack(
    scope=app,
    id="cognito-mfa-deploy-pipeline",
    artifact_bucket=build_stack.artifact_bucket,
    backend_fn=cognito_stack.backend_fn,
    api=cognito_stack.api,
    static_website_bucket=cognito_stack.static_website_bucket,
)
app.synth()
