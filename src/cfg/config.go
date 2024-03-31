package cfg

type Settings struct {
	AwsRegion     string
	ApiKeyTable   string
	RootKeyPrefix string
}

var Config = Settings{
	AwsRegion:     "eu-west-2",
	ApiKeyTable:   "ApiKeyTableDev",
	RootKeyPrefix: "keyify_",
}
