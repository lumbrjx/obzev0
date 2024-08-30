TARGET = tc.o
INTERFACE = enp1s0 
US_DIR = user_space 
CFLAGS = -I/usr/include -I/usr/include/x86_64-linux-gnu -I/usr/include/x86_64-linux-gnu/bits -I/usr/include/x86_64-linux-gnu/sys -I/usr/include/bpf

.PHONY: install-deps
install-deps:
	sudo apt update
	sudo apt install -y clang llvm libelf-dev linux-headers-$$(uname -r) build-essential
	cd $(US_DIR) && go mod tidy

# Compile the eBPF program
$(TARGET): main.c
	clang $(CFLAGS) -O2 -g -target bpf -c bpf/tc.c -o $(TARGET) 

# Load the eBPF program manually
.PHONY: load
load: $(TARGET)
	sudo tc qdisc add dev $(INTERFACE) clsact
	sudo tc filter add dev $(INTERFACE) ingress bpf da obj $(TARGET) sec tc
	sudo tc filter add dev $(INTERFACE) egress bpf da obj $(TARGET) sec tc

# View bpf_printk output
.PHONY: view-manual
view:
	sudo cat /sys/kernel/debug/tracing/trace_pipe

.PHONY: view-tcp-manual
view-tcp:
	sudo cat /sys/kernel/debug/tracing/trace_pipe | grep TCP

.PHONY: view-udp-manual
view-udp:
	sudo cat /sys/kernel/debug/tracing/trace_pipe | grep UDP


# build user space program
.PHONY: build-US
build-US:
	cd user_space && go build -o tc_US tc.go

# start user space program
.PHONY: start-US
start-US:
	sudo ./user_space/tc_US $(INTERFACE)

# Remove the filters and qdisc when done manually
.PHONY: clean
clean:
	sudo tc filter del dev $(INTERFACE) ingress
	sudo tc filter del dev $(INTERFACE) egress
	sudo tc qdisc del dev $(INTERFACE) clsact
	rm -f $(TARGET)
	rm -f user_space/tc_US

# All
.PHONY: all
all: install-deps $(TARGET) build-US start-US
