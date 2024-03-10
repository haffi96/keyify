import os


class CdkSettings:
    def __init__(self):
        self.region = os.getenv("CDK_REGION", "eu-west-2")
        self.account = os.getenv("CDK_ACCOUNT", "905418485198")


settings = CdkSettings()
