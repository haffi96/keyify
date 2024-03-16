from constructs import Construct
from aws_cdk import aws_lambda


class CodeSourceConstruct(Construct):
    def __init__(self, scope: Construct, construct_id: str) -> None:
        super().__init__(scope, construct_id)

        self.construct_id = construct_id

    def source(self) -> aws_lambda.Code:
        # if os.environ["CDK_ACCOUNT"] == "000000000000":
        #     return aws_lambda.Code.from_bucket(
        #         aws_s3.Bucket.from_bucket_name(
        #             self,
        #             "hot-reloading-lambda-bucket",
        #             "hot-reload",
        #         ),
        #         f"dist/{self.construct_id}/function.zip",
        #     )
        # else:
        return aws_lambda.Code.from_asset(
            path=f"dist/{self.construct_id}/function.zip",
        )
