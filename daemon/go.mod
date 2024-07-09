module obzev0/daemon

go 1.22.5

require (
	golang.org/x/net v0.27.0
	google.golang.org/grpc v1.65.0
	gopkg.in/yaml.v2 v2.4.0
	obzev0/common v0.0.0
)

replace obzev0/common => ../common

require (
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
