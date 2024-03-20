from aws_cdk import aws_apigateway
from aws_cdk.aws_lambda import Function


from constructs import Construct


class ApiGatewayStack(aws_apigateway.RestApi):
    def __init__(self, scope: Construct, construct_id: str, stage: str) -> None:
        super().__init__(
            scope,
            construct_id,
            rest_api_name=f"{construct_id}ApiGateway{stage}",
            deploy=True,
            cloud_watch_role=True,
            deploy_options=aws_apigateway.StageOptions(
                stage_name=stage.lower(),
                logging_level=aws_apigateway.MethodLoggingLevel.INFO,
                data_trace_enabled=True,
                metrics_enabled=True,
                tracing_enabled=True,
            ),
        )

    def add_lambda_integration(
        self, path: str, method: str, lambda_function: Function
    ):
        resource = self.root.resource_for_path(path)
        resource.add_method(
            method,
            aws_apigateway.LambdaIntegration(lambda_function),
        )
