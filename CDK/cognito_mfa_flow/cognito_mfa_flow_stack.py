from aws_cdk import (
    core,
    aws_iam as iam,
    aws_cognito as cognito,
    aws_apigateway as apigateway,
    aws_lambda as _lambda,
    aws_cognito as cognito

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

        client = pool.add_client(
            id="customer-app-client",
            auth_flows=cognito.AuthFlow(
                user_password=True,
                refresh_token=True),
        )
