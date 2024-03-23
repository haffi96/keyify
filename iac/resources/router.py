from constructs import Construct

from iac.resources.api_gateway import ApiGatewayStack
from iac.resources.code_source import CodeSourceConstruct
from iac.resources.generic_lambda import GenericGoLambdaFunction


class ApiKeyServiceRouter(Construct):
    def __init__(
        self,
        scope: Construct,
        construct_id: str,
        stage: str,
        api_gateway: ApiGatewayStack,
    ) -> None:
        super().__init__(scope, construct_id)

        self.stage = stage

        for lambda_id, route, method, description in [
            ("TestGo", "/test-go", "GET", "test go on aws"),
            ("CreateRootKey", "/rootKey", "POST", "create a root key in db"),
            ("GetApiKey", "/key", "GET", "get an api key in db"),
            ("CreateApiKey", "/key", "POST", "create an api key in db"),
            ("VerifyApiKey", "/verifyKey", "POST", "verify an api key in db"),
            ("CreateApi", "/api", "POST", "create an api in db"),
            ("ListApiKeys", "/keys", "GET", "list all api keys for an api"),
        ]:
            lambda_function = GenericGoLambdaFunction(
                self,
                lambda_id,
                stage=self.stage,
                code_source=CodeSourceConstruct(self, lambda_id).source(),
                description=f"Lambda to {description}",
            )

            api_gateway.add_lambda_integration(route, method, lambda_function)
