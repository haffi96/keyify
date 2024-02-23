from aws_cdk import aws_lambda
from aws_cdk import aws_iam
from constructs import Construct


class GenericGoLambdaFunction(aws_lambda.DockerImageFunction):
    def __init__(self, scope: Construct, construct_id: str, description: str | None) -> None:
        super().__init__(
            scope,
            f"{construct_id}Lambda",
            code=aws_lambda.DockerImageCode.from_image_asset(
                f"apikeyservice/lambdas/{construct_id}",
                asset_name=f"{construct_id.lower()}-lambda",
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
