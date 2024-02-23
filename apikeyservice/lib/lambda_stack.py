from aws_cdk import aws_lambda
from aws_cdk import aws_iam
from constructs import Construct


class GenericGoLambdaFunction(aws_lambda.Function):
    def __init__(
        self, scope: Construct, construct_id: str, description: str | None
    ) -> None:
        super().__init__(
            scope,
            f"{construct_id}Lambda",
            runtime=aws_lambda.Runtime.PROVIDED_AL2,
            handler="bootstrap",
            code=aws_lambda.Code.from_asset(
                path=f"dist/{construct_id}/function.zip",
            ),
            initial_policy=[
                aws_iam.PolicyStatement(
                    actions=["dynamodb:*"],
                    resources=["*"],
                    effect=aws_iam.Effect.ALLOW,
                )
            ],
            description=description,
        )
