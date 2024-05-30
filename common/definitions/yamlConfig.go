package definitions

type Config struct {
	Delays DelaysConfig `yaml:"delays"`
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

type DelaysConfig struct {
	ReqDelay int32 `yaml:"reqDelay"`
	ResDelay int32 `yaml:"resDelay"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
}
type ClientConfig struct {
	Port string `yaml:"port"`
}
