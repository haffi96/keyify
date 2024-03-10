package cfg

type Settings struct {
	AwsRegion string
}

var Config = Settings{
	AwsRegion: "eu-west-2",
}
