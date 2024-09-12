module obzev0/cli

go 1.22.5

require (
	google.golang.org/grpc v1.66.2
	gopkg.in/yaml.v2 v2.4.0
	obzev0/common v0.0.0-00010101000000-000000000000
)

replace obzev0/common => ../../common

require (
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
