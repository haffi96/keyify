#!/usr/bin/env python3
import aws_cdk as cdk

from apikeyservice.apikeyservice_stack import ApikeyserviceStack
from settings import settings

app = cdk.App()

prod_stage = cdk.Stage(
    app,
    "prod-env",
    env=cdk.Environment(
        account=settings.account,
        region=settings.region,
    ),
)

dev_stage = cdk.Stage(
    app,
    "dev-env",
    env=cdk.Environment(
        account=settings.account,
        region=settings.region,
    ),
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
