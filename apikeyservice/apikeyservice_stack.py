from aws_cdk import (
    aws_dynamodb,
    Stack,
)
from constructs import Construct
from apikeyservice.lib.lambda_stack import GenericGoLambdaFunction
from apikeyservice.lib.dynamo import DynamoDBTable
from apikeyservice.lib.api_gateway import ApiGatewayStack


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

        test_go_lambda = GenericGoLambdaFunction(
            self,
            "TestGo",
            stage=stage,
            description="Lambda to test go on aws",
        )
        get_api_key_lambda = GenericGoLambdaFunction(
            self,
            "GetApiKey",
            stage=stage,
            description="Lambda to get an api key in db",
        )
        create_api_key_lambda = GenericGoLambdaFunction(
            self,
            "CreateApiKey",
            stage=stage,
            description="Lambda to create an api key in db",
        )

        verify_api_key_lambda = GenericGoLambdaFunction(
            self,
            "VerifyApiKey",
            stage=stage,
            description="Lambda to verify an api key in db",
        )

        api_gateway = ApiGatewayStack(
            self,
            "ApiKeyService",
            stage,
        )

        api_gateway.add_lambda_integration("/test-go", "GET", test_go_lambda)
        api_gateway.add_lambda_integration("/key", "GET", get_api_key_lambda)
        api_gateway.add_lambda_integration(
            "/key", "POST", create_api_key_lambda
        )
        api_gateway.add_lambda_integration(
            "/key", "POST", verify_api_key_lambda
        )
