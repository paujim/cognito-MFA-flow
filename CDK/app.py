#!/usr/bin/env python3

from aws_cdk import core

from cognito_mfa_flow.cognito_mfa_flow_stack import(
    CognitoMfaFlowStack,
    MFAAPIStack,
)


app = core.App()
cognito_stack = CognitoMfaFlowStack(app, "cognito-mfa-flow")
MFAAPIStack(app, "cognito-mfa-apigateway", cognito_stack.pool, cognito_stack.client)
app.synth()
