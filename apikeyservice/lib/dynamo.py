from aws_cdk import aws_dynamodb


class DynamoDBTable(aws_dynamodb.Table):
    def __init__(self, scope, construct_id, **kwargs):
        super().__init__(
            scope,
            construct_id,
            **kwargs,
            point_in_time_recovery=True,
        )