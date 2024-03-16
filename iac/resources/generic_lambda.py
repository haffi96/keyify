from aws_cdk import aws_lambda
from aws_cdk import aws_iam
from constructs import Construct


class GenericGoLambdaFunction(aws_lambda.Function):
    def __init__(
        self,
        scope: Construct,
        construct_id: str,
        stage: str,
        code_source: aws_lambda.Code,
        description: str | None,
    ) -> None:
        super().__init__(
            scope,
            f"{construct_id}Lambda",
            runtime=aws_lambda.Runtime.PROVIDED_AL2,
            architecture=aws_lambda.Architecture.ARM_64,
            handler="bootstrap",
            code=code_source,
            initial_policy=[
                aws_iam.PolicyStatement(
                    actions=["dynamodb:*"],
                    resources=["*"],
                    effect=aws_iam.Effect.ALLOW,
                )
            ],
            function_name=f"{construct_id}Lambda{stage}",
            description=description,
        )
