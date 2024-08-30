package tcanalyser

import (
	"fmt"
	"net"

	"github.com/cilium/ebpf"
)

func intToIP(ip uint32) net.IP {
	return net.IPv4(byte(ip), byte(ip>>8), byte(ip>>16), byte(ip>>24))
}

func loadBpfSpec(path string) (*ebpf.CollectionSpec, error) {
	spec, err := ebpf.LoadCollectionSpec(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load BPF spec: %v", err)
	}
	if eventsMap, ok := spec.Maps["events"]; ok {
		eventsMap.Type = ebpf.PerfEventArray
	}

	return spec, nil
}
func tcpFlagsToString(flags uint8) string {
	flagNames := []struct {
		mask uint8
		name string
	}{
		{0x01, "FIN"},
		{0x02, "SYN"},
		{0x04, "RST"},
		{0x08, "PSH"},
		{0x10, "ACK"},
		{0x20, "URG"},
		{0x40, "ECE"},
		{0x80, "CWR"},
	}

	var result []string
	for _, f := range flagNames {
		if flags&f.mask != 0 {
			result = append(result, f.name)
		}
	}
	if len(result) == 0 {
		return "NONE"
	}
	return fmt.Sprintf("0x%x (%s)", flags, fmt.Sprintf("%s", result))
}
func printPacket(event Event) {
	direction := "Ingress"
	if event.Direction == 1 {
		direction = "Egress"
	}
	flags := ""
	protoc := "?"
	s_msg := "%s %s: src=%s:%d -> dst=%s:%d | proto=%d | flags=%s\n"
	switch event.Protocol {
	case 6: // TCP
		flags = tcpFlagsToString(event.TcpFlags)
		protoc = "TCP"
	case 17: // UDP
		protoc = "UDP"
		s_msg = "%s %s: src=%s:%d -> dst=%s:%d | proto=%d\n"
	}
	fmt.Printf(s_msg,
		direction,
		protoc,
		intToIP(event.SrcIP), event.SrcPort,
		intToIP(event.DstIP), event.DstPort,
		event.Protocol,
		flags)

}
