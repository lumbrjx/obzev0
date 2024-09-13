<h1 align="center">obzev0 : Chaos Engineering Tool</h1>

<p align="center">

Obzev0 is a chaos engineering tool designed to help you test the resilience of your systems by simulating real-world failures. It allows you to define and execute chaos experiments to uncover weaknesses in your infrastructure and applications.
</p>

<p align="center">
  <img src="./assets/tn.jpg" />
</p>

## Tech Stack

- **Language:** Go (for the controller and gRPC server), C (for the eBPF program).
- **Frameworks/Libraries:** Kubernetes client-go (for interacting with the Kubernetes API), gRPC (for communication between components), eBPF (for kernel space programming).
- **Monitoring:** Prometheus and Grafana for monitoring and visualization.
## Architecture

Obzev0 is built using a microservices architecture, with the following components:

- **Controller:** Responsible for watching Custom Resource Definitions (CRDs) representing chaos scenarios and dispatching work to the DaemonSet.
- **DaemonSet:** Runs on every node in the Kubernetes cluster and acts as a gRPC server. It executes the chaos scenarios and communicates with the eBPF program in the kernel space.
- **eBPF Program:** Written in C, it runs in the kernel space and is responsible for monitoring and manipulating network traffic for performance monitoring.

## Features

### Latency Injection
- Simulate network latency by adding configurable delays to incoming and outgoing requests
- Customize delay duration for each service to test specific latency thresholds

### Packet Manipulation
- Introduce packet loss with configurable drop rates
- Corrupt specific packets to test handling of incomplete or malformed data
- Apply packet manipulation rules for defined periods

### Traffic Monitoring with eBPF
- Track network traffic on specific interfaces using eBPF programs
- Collect detailed metrics on packet flow, dropped packets, and corrupted packets

### Self-cleaning TCP Proxy for Testing
- Route traffic through a transparent TCP proxy to apply latency and packet manipulation
- Configure different target services for testing under various network conditions

### Dynamic Configuration
- Manage chaos experiments through a simple gRPC interface
- Enable or disable individual services based on testing needs

### Metrics Exposure to Prometheus
- Expose critical chaos engineering metrics to Prometheus for monitoring and alerting
- View and analyze metrics such as latency impact, packet loss, and traffic anomalies over time

## Installation

### Prerequisites
- Go 1.17+
- Helm and Kubernetes (for cloud-native environments)
- Prometheus (for collecting chaos engineering metrics)

### Install with Helm

1. Add the Helm repository:
```bash
helm repo add obzev0 https://lumbrjx.github.io/obzev0/chart
helm repo update
```

2. Install Obzev0:
```bash
helm install obzev0 obzev0/obzev0
```

### Install obzevMini (for non-Kubernetes environments)

Download the binary:
```bash
curl -o obzevMini https://raw.githubusercontent.com/lumbrjx/obzev0/main/cmd/cli/obzev0mini
chmod +x obzevMini
```

### Build from Source

1. Clone the repository:
```bash
git clone https://github.com/lumbrjx/obzev0.git
cd obzev0
```

2. For Helm and Kubernetes:
   - Build the DaemonSet:
     ```bash
     make build-daemon TAG=<tag>
     ```
   - Build the controller:
     ```bash
     make build-controller TAG=<tag>
     ```
   - Package Helm chart:
     ```bash
     make package-chart
     ```
Make sure to give your docker image the "latest" tag and push it to a container registry.

3. For obzevMini:
   - Build the binary:
     ```bash
     make build-cli
     ```

## Usage

### Helm

1. Install the chart:
```bash
helm install <release-name> obzev0/obzev0
```

2. Grant controller permissions:
```bash
kubectl create clusterrolebinding permissive-binding \
  --clusterrole=cluster-admin \
  --serviceaccount=controller-system:controller-controller-manager
```

3. Label nodes:
   - Control plane:
     ```bash
     kubectl label nodes <node-name> node-role.kubernetes.io/control-plane=""
     ```
   - Worker nodes:
     ```bash
     kubectl label nodes <node-name> node-role.kubernetes.io/worker=""
     ```

4. Apply Custom Resource:

The chart folder already contains a Custom Resource and you still can define your own one. here's the [CR file](https://github.com/lumbrjx/obzev0/blob/main/chart/templates/obzev0resource.yaml) and apply the changes by running:

```bash
kubectl apply -f path/to/CR/file
```

### run locally with kind

after clonning the repo you can run the following command and have everything done for you:

```bash
make all TAG=latest
```
it will use a simple dummy express server you can access it [here](https://github.com/lumbrjx/obzev0/blob/main/expressdep.yaml) and the [service](https://github.com/lumbrjx/obzev0/blob/main/expresssvc.yaml)


### obzevMini

1. Launch the Daemon Docker container:
```bash
docker run -p 50051:50051 lumbrjx/obzev0-grpc-daemon:latest
```

2. Initialize configuration:
```bash
obzevMini init -Addr=127.0.0.1:50051 -dist=path/to/destination/folder/
```

3. Apply configuration:
```bash
obzevMini apply -path=path/to/configuration/file
```

## Contributing

We welcome contributions to Obzev0! Whether you're fixing bugs, adding features, or improving documentation, your contributions are valuable.

For questions or discussions, contact the project maintainers at kaibitayeb01@gmail.com.
