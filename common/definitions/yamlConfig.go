package definitions

// Internal config types used by the latency daemon service.

type DelaysConfig struct {
	ReqDelay int32
	ResDelay int32
}

type ServerConfig struct {
	Port string
}

type ClientConfig struct {
	Port string
}

// LatencyInternalConfig is used by LaunchTcp inside the daemon; it is separate
// from the YAML-facing LatencySvcConfig to avoid coupling the daemon to the CLI schema.
type LatencyInternalConfig struct {
	Delays DelaysConfig
	Server ServerConfig
	Client ClientConfig
}

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
