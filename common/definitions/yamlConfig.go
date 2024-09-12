package definitions

type LatencySvcConfig struct {
	Enabled  bool   `yaml:"enabled"`
	ReqDelay int    `yaml:"reqDelay"`
	ResDelay int    `yaml:"resDelay"`
	Server   string `yaml:"server"`
	Client   string `yaml:"client"`
}

type TcAnalyserSvcConfig struct {
	Enabled  bool   `yaml:"enabled"`
	NetIFace string `yaml:"netIFace"`
}

type PacketManipulationSvcConfig struct {
	Enabled         bool   `yaml:"enabled"`
	Server          string `yaml:"server"`
	Client          string `yaml:"client"`
	DropRate        string `yaml:"dropRate"`
	CorruptRate     string `yaml:"corruptRate"`
	DurationSeconds int    `yaml:"durationSeconds"`
}

type Config struct {
	ServerAddr                  string                      `yaml:"serverAddr"`
	LatencySvcConfig            LatencySvcConfig            `yaml:"latencySvcConfig"`
	TcAnalyserSvcConfig         TcAnalyserSvcConfig         `yaml:"tcAnalyserSvcConfig"`
	PacketManipulationSvcConfig PacketManipulationSvcConfig `yaml:"packetManipulationSvcConfig"`
}
