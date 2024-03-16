from aws_cdk import aws_dynamodb, Stack
from constructs import Construct
from iac.resources.dynamo import DynamoDBTable
from iac.resources.api_gateway import ApiGatewayStack
from iac.resources.router import ApiKeyServiceRouter


class ApikeyserviceStack(Stack):
    def __init__(
        self, scope: Construct, construct_id: str, stage: str, **kwargs
    ) -> None:
        super().__init__(scope, construct_id, **kwargs)

        DynamoDBTable(
            self,
            f"ApiKeyTable{stage}",
            table_name=f"ApiKeyTable{stage}",
            partition_key=aws_dynamodb.Attribute(
                name="pk", type=aws_dynamodb.AttributeType.STRING
            ),
            sort_key=aws_dynamodb.Attribute(
                name="sk", type=aws_dynamodb.AttributeType.STRING
            ),
        )

        api_gateway = ApiGatewayStack(self, "ApiKeyService", stage=stage)

        ApiKeyServiceRouter(
            self, "ApiRouter", stage=stage, api_gateway=api_gateway
        )
