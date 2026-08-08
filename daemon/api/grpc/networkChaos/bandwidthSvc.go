package networkchaos

import (
	"fmt"
	"sync/atomic"

	"github.com/vishvananda/netlink"
)

// activeBandwidthIface stores the interface name where a TBF qdisc was added.
var activeBandwidthIface atomic.Value

// startBandwidth installs a TBF qdisc on the named interface to cap outbound traffic.
func startBandwidth(iface string, rateKbps uint64, burstKb uint32) error {
	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("link %q not found: %w", iface, err)
	}

	// Remove any previous qdisc first so we don't stack them.
	_ = stopBandwidth(iface)

	burst := uint32(burstKb) * 1024
	if burst == 0 {
		burst = 32 * 1024
	}
	rateBytes := rateKbps * 128 // kbps → bytes/s

	qdisc := &netlink.Tbf{
		QdiscAttrs: netlink.QdiscAttrs{
			LinkIndex: link.Attrs().Index,
			Handle:    netlink.MakeHandle(1, 0),
			Parent:    netlink.HANDLE_ROOT,
		},
		Rate:   rateBytes,
		Limit:  uint32(rateBytes),
		Buffer: burst,
	}
	if err := netlink.QdiscAdd(qdisc); err != nil {
		return fmt.Errorf("QdiscAdd on %q: %w", iface, err)
	}

	activeBandwidthIface.Store(iface)
	return nil
}

func stopBandwidth(iface string) error {
	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("link %q not found: %w", iface, err)
	}
	qdiscs, err := netlink.QdiscList(link)
	if err != nil {
		return err
	}
	for _, q := range qdiscs {
		if q.Attrs().Parent == netlink.HANDLE_ROOT {
			if err := netlink.QdiscDel(q); err != nil {
				return fmt.Errorf("QdiscDel on %q: %w", iface, err)
			}
		}
	}
	return nil
}
