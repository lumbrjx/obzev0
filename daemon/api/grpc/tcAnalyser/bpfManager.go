package tcanalyser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"
	"github.com/vishvananda/netlink"
)

type Event struct {
	SrcIP     uint32
	DstIP     uint32
	SrcPort   uint16
	DstPort   uint16
	Protocol  uint8
	Direction uint8
	TcpFlags  uint8
}

func bpfLoader(interf string) {
	ifaceName := interf

	spec, err := loadBpfSpec("tc.o")
	if err != nil {
		log.Fatalf("loading eBPF spec: %v", err)
	}

	var objs struct {
		TcIngress *ebpf.Program `ebpf:"tc_ingress"`
		TcEgress  *ebpf.Program `ebpf:"tc_egress"`
		Events    *ebpf.Map     `ebpf:"events"`
	}
	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.TcIngress.Close()
	defer objs.TcEgress.Close()
	defer objs.Events.Close()

	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		log.Fatalf("getting interface: %v", err)
	}

	// Ensure clsact qdisc exists
	qdiscs, err := netlink.QdiscList(link)
	if err != nil {
		log.Fatalf("listing qdiscs: %v", err)
	}
	clsactExists := false
	for _, qdisc := range qdiscs {
		if _, ok := qdisc.(*netlink.GenericQdisc); ok &&
			qdisc.Type() == "clsact" {
			clsactExists = true
			break
		}
	}
	if !clsactExists {
		qdisc := &netlink.GenericQdisc{
			QdiscAttrs: netlink.QdiscAttrs{
				LinkIndex: link.Attrs().Index,
				Handle:    netlink.MakeHandle(0xffff, 0),
				Parent:    netlink.HANDLE_CLSACT,
			},
			QdiscType: "clsact",
		}
		if err := netlink.QdiscAdd(qdisc); err != nil {
			log.Fatalf("adding clsact qdisc: %v", err)
		}
		fmt.Println("Added clsact qdisc")
	} else {
		fmt.Println("clsact qdisc already exists")
	}

	ingressFilter := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: link.Attrs().Index,
			Parent:    netlink.HANDLE_MIN_INGRESS,
			Handle:    netlink.MakeHandle(0, 1),
			Protocol:  syscall.ETH_P_ALL,
		},
		Fd:           objs.TcIngress.FD(),
		Name:         "tc_ingress",
		DirectAction: true,
	}
	if err := netlink.FilterAdd(ingressFilter); err != nil {
		log.Fatalf("adding ingress filter: %v", err)
	}
	fmt.Println("Added ingress filter")

	egressFilter := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: link.Attrs().Index,
			Parent:    netlink.HANDLE_MIN_EGRESS,
			Handle:    netlink.MakeHandle(0, 1),
			Protocol:  syscall.ETH_P_ALL,
		},
		Fd:           objs.TcEgress.FD(),
		Name:         "tc_egress",
		DirectAction: true,
	}
	if err := netlink.FilterAdd(egressFilter); err != nil {
		log.Fatalf("adding egress filter: %v", err)
	}
	fmt.Println("Added egress filter")

	// perf reader
	reader, err := perf.NewReader(objs.Events, os.Getpagesize())
	if err != nil {
		log.Fatalf("creating perf reader: %v", err)
	}
	defer reader.Close()

	// read events
	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				if err == perf.ErrClosed {
					return
				}
				log.Printf("reading from perf event reader: %s", err)
				continue
			}

			if record.LostSamples != 0 {
				log.Printf(
					"perf event ring buffer full, dropped %d samples",
					record.LostSamples,
				)
				continue
			}

			var event Event
			if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
				log.Printf("parsing perf event: %s", err)
				continue
			}
			printPacket(event)

		}
	}()

	fmt.Printf("eBPF programs attached to interface %s\n", ifaceName)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Received interrupt, cleaning up...")

	// Clean up filters
	filters, err := netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
	if err != nil {
		log.Printf("error listing ingress filters: %v", err)
	} else {
		for _, filter := range filters {
			if bpfFilter, ok := filter.(*netlink.BpfFilter); ok && bpfFilter.Name == "tc_ingress" {
				if err := netlink.FilterDel(bpfFilter); err != nil {
					log.Printf("error removing ingress filter: %v", err)
				} else {
					fmt.Println("Removed ingress filter")
				}
			}
		}
	}

	filterss, err := netlink.FilterList(link, netlink.HANDLE_MIN_EGRESS)
	if err != nil {
		log.Printf("error listing egress filters: %v", err)
	} else {
		for _, filter := range filterss {
			if bpfFilter, ok := filter.(*netlink.BpfFilter); ok && bpfFilter.Name == "tc_egress" {
				if err := netlink.FilterDel(bpfFilter); err != nil {
					log.Printf("error removing egress filter: %v", err)
				} else {
					fmt.Println("Removed egress filter")
				}
			}
		}
	}

	qdisc := &netlink.GenericQdisc{
		QdiscAttrs: netlink.QdiscAttrs{
			LinkIndex: link.Attrs().Index,
			Handle:    netlink.MakeHandle(0xffff, 0),
			Parent:    netlink.HANDLE_CLSACT,
		},
		QdiscType: "clsact",
	}
	if err := netlink.QdiscDel(qdisc); err != nil {
		log.Fatalf("deleting clsact qdisc: %v", err)
	}
	fmt.Println("Deleted clsact qdisc")
}
