package boxdns

import (
	"fmt"

	tun "github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/control"
	logger2 "github.com/sagernet/sing/common/logger"
)

// Always-on and independent of the box lifecycle; on Windows it also drives system DNS, elsewhere HandleSystemDNS is a no-op.
var DnsManagerInstance *DnsManager

type DnsManager struct {
	Monitor tun.DefaultInterfaceMonitor
	// Kept current by Monitor on every network change.
	Finder  control.InterfaceFinder
	lastIfc *control.Interface
}

func init() {
	logger := logger2.NOP()
	updMonitor, err := tun.NewNetworkUpdateMonitor(logger)
	if err != nil {
		fmt.Println("Could not create NetworkUpdateMonitor")
		return
	}
	finder := control.NewDefaultInterfaceFinder()
	monitor, err := tun.NewDefaultInterfaceMonitor(updMonitor, logger, tun.DefaultInterfaceMonitorOptions{
		InterfaceFinder: finder,
	})
	if err != nil {
		fmt.Println("Could not create DefaultInterfaceMonitor")
		return
	}
	DnsManagerInstance = &DnsManager{Monitor: monitor, Finder: finder}
	monitor.RegisterCallback(DnsManagerInstance.HandleSystemDNS)
	if err = updMonitor.Start(); err != nil {
		fmt.Println("Could not start updMonitor")
		return
	}
	if err = monitor.Start(); err != nil {
		fmt.Println("Could not start monitor")
		return
	}
}

// nil when the monitor is unavailable; TUN and loopback are excluded, so the result is safe to bind egress to while throne-tun is up.
func DefaultInterface() *control.Interface {
	if DnsManagerInstance == nil || DnsManagerInstance.Monitor == nil {
		return nil
	}
	return DnsManagerInstance.Monitor.DefaultInterface()
}
