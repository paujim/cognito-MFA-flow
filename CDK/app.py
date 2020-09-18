#!/usr/bin/env python3

from aws_cdk import core

from cognito_mfa_flow.cognito_mfa_flow_stack import(
    CognitoMfaFlowStack,
)


app = core.App()
cognito_stack = CognitoMfaFlowStack(
    scope=app,
    id="cognito-mfa-flow",
)
app.synth()
