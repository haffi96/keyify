from aws_cdk import aws_lambda
from aws_cdk import aws_iam
from constructs import Construct
from aws_cdk import BundlingOptions


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
                path=f"apikeyservice/lambdas/{construct_id}",
                bundling=BundlingOptions(
                    image=aws_lambda.Runtime.PROVIDED_AL2.bundling_image,
                    environment={
                        "GOOS": "linux",
                        "GOARCH": "amd64",
                        "CGO_ENABLED": "0",
                    },
                    user="root",
                    command=[
                        "bash",
                        "-c",
                        "go build -tags lambda.norpc -o /asset-output/bootstrap main.go && cd /asset-output && zip function.zip bootstrap",
                    ],
                ),
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
