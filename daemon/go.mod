module obzev0/daemon

go 1.22.5

require (
	github.com/go-kit/log v0.2.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
	github.com/prometheus/client_golang v1.19.1
	golang.org/x/net v0.27.0
	google.golang.org/grpc v1.65.0
	gopkg.in/yaml.v2 v2.4.0
	obzev0/common v0.0.0
)

replace obzev0/common => ../common

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240709173604-40e1e62336c5 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
