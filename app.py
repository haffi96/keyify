#!/usr/bin/env python3
import aws_cdk as cdk

from iac.stack import ApikeyserviceStack
from settings import settings

app = cdk.App()

prod_stage = cdk.Stage(
    app,
    "Prod",
    env=cdk.Environment(
        account=settings.account,
        region=settings.region,
    ),
    stage_name="Prod",
)

dev_stage = cdk.Stage(
    app,
    "Dev",
    env=cdk.Environment(
        account=settings.account,
        region=settings.region,
    ),
    stage_name="Dev",
)

ApikeyserviceStack(
    prod_stage,
    "ApiKeyServiceStack",
    stage=prod_stage.stage_name,
)

ApikeyserviceStack(
    dev_stage,
    "ApiKeyServiceStack",
    stage=dev_stage.stage_name,
)

app.synth()
