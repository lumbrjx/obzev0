<h1 align="center">obzev0</h1>

<p align="center">
 obzev0: Chaos Engineering Platform
obzev0 is a chaos engineering platform designed to help you test the resilience of your systems by simulating real-world failures. It allows you to define and execute chaos experiments to uncover weaknesses in your infrastructure and applications.
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
