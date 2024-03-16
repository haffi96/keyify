package cfg

type Settings struct {
	AwsRegion   string
	ApiKeyTable string
}

var Config = Settings{
	AwsRegion:   "eu-west-2",
	ApiKeyTable: "ApiKeyTableDev",
}
