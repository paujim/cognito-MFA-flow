#!/usr/bin/env python3

from aws_cdk import core

from cognito_mfa_flow.cognito_mfa_flow_stack import CognitoMfaFlowStack


app = core.App()
CognitoMfaFlowStack(app, "cognito-mfa-flow")

app.synth()
