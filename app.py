#!/usr/bin/env python3
import aws_cdk as cdk

from iac.stack import ApikeyserviceStack
from settings import settings

app = cdk.App()

prod_stage = cdk.Stage(
    app,
    "prod-env",
    env=cdk.Environment(
        account=settings.account,
        region=settings.region,
    ),
    stage_name="prod",
)

dev_stage = cdk.Stage(
    app,
    "dev-env",
    env=cdk.Environment(
        account=settings.account,
        region=settings.region,
    ),
    stage_name="dev",
)

ApikeyserviceStack(
    prod_stage,
    "ApiKeyServiceStack",
    stage="Prod",
)

ApikeyserviceStack(
    dev_stage,
    "ApiKeyServiceStack",
    stage="Dev",
)

app.synth()
