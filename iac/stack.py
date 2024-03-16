from aws_cdk import aws_dynamodb, Stack
from constructs import Construct
from iac.resources.generic_lambda import GenericGoLambdaFunction
from iac.resources.dynamo import DynamoDBTable
from iac.resources.api_gateway import ApiGatewayStack


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
        create_root_key_lambda = GenericGoLambdaFunction(
            self,
            "CreateRootKey",
            stage=stage,
            description="Lambda to create a root key in db",
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
        api_gateway.add_lambda_integration(
            "/rootKey", "PUT", create_root_key_lambda
        )
        api_gateway.add_lambda_integration("/key", "GET", get_api_key_lambda)
        api_gateway.add_lambda_integration(
            "/key", "POST", create_api_key_lambda
        )
        api_gateway.add_lambda_integration(
            "/verifyKey", "POST", verify_api_key_lambda
        )
