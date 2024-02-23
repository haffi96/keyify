import aws_cdk as core
import aws_cdk.assertions as assertions

from apikeyservice.apikeyservice_stack import ApikeyserviceStack

# example tests. To run these tests, uncomment this file along with the example
# resource in apikeyservice/apikeyservice_stack.py
def test_sqs_queue_created():
    app = core.App()
    stack = ApikeyserviceStack(app, "apikeyservice", "pytest")
    _ = assertions.Template.from_stack(stack)

#     template.has_resource_properties("AWS::SQS::Queue", {
#         "VisibilityTimeout": 300
#     })
