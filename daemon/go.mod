module obzev0/daemon

go 1.22.5

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
	github.com/prometheus/client_golang v1.19.1
	golang.org/x/net v0.27.0
	google.golang.org/grpc v1.66.0
	gopkg.in/yaml.v2 v2.4.0
	obzev0/common v0.0.0
)

replace obzev0/common => ../common

require (
	github.com/cilium/ebpf v0.16.0
	github.com/sirupsen/logrus v1.9.3
	github.com/vishvananda/netlink v1.3.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240709173604-40e1e62336c5 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
