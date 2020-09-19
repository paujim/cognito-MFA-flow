#!/usr/bin/env python3

from aws_cdk import core

from cognito_mfa_flow.cognito_mfa_flow_stack import(
    BuildPipelineStack,
    CognitoMfaFlowStack,
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
app.synth()
