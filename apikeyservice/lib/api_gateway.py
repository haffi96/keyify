from aws_cdk import aws_apigateway
from aws_cdk import RemovalPolicy
from aws_cdk.aws_logs import LogGroup, RetentionDays
from aws_cdk.aws_lambda import Function


from constructs import Construct


class ApiGatewayStack(aws_apigateway.RestApi):
    def __init__(self, scope: Construct, construct_id: str, stage: str) -> None:
        super().__init__(
            scope,
            f"{construct_id}ApiGateway",
            rest_api_name=f"{construct_id}ApiGateway{stage}",
            # cloud_watch_role=True,
            # deploy_options=aws_apigateway.StageOptions(
            #     access_log_destination=aws_apigateway.LogGroupLogDestination(
            #         LogGroup(
            #             self,
            #             f"ApiGatewayLogGroup{stage}",
            #             log_group_name=f"ApiGatewayLogGroup{stage}",
            #             removal_policy=RemovalPolicy.DESTROY,
            #             retention=RetentionDays.ONE_MONTH,
            #         )
            #     ),
            # ),
        )

    def add_lambda_integration(
        self, path: str, method: str, lambda_function: Function
    ):
        resource = self.root.resource_for_path(path)
        resource.add_method(
            method,
            aws_apigateway.LambdaIntegration(lambda_function),
        )
