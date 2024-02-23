#!/usr/bin/env python3
import os

import aws_cdk as cdk

from apikeyservice.apikeyservice_stack import ApikeyserviceStack


app = cdk.App()

prod_stage = cdk.Stage(
    app,
    "prod-env",
    env=cdk.Environment(
        account="905418485198",
        region="eu-west-2",
    ),
)

dev_stage = cdk.Stage(
    app,
    "dev-env",
    env=cdk.Environment(
        account="905418485198",
        region="eu-west-2",
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
