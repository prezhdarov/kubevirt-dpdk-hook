package hook

import (
	"fmt"
	"path/filepath"
	"strings"

	vmschema "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"
	"libvirt.org/go/libvirtxml"
)

func configNetwork(interfaces []vmschema.Interface, domainInterfaces []libvirtxml.DomainInterface) error {

	for _, iface := range interfaces {

		if iface := lookupIfaceByAliasName(domainInterfaces, iface.Name); iface != nil {
			iface.Target.Managed = "yes"

			iface.Source = &libvirtxml.DomainInterfaceSource{
				VHostUser: &libvirtxml.DomainChardevSource{
					UNIX: &libvirtxml.DomainChardevSourceUNIX{
						Path: filepath.Join(OVSSocketDir, fmt.Sprintf("%s.sock", strings.TrimPrefix(iface.Target.Dev, "tap"))),
						Mode: "server",
					},
				},
			}

			iface.Driver = &libvirtxml.DomainInterfaceDriver{
				Name:        "vhost",
				Queues:      4,
				RXQueueSize: 1024,
				TXQueueSize: 1024,
				Guest: &libvirtxml.DomainInterfaceDriverGuest{
					CSum: "on",
					TSO4: "on",
					TSO6: "on",
					ECN:  "on",
				},
				Host: &libvirtxml.DomainInterfaceDriverHost{
					CSum:     "on",
					GSO:      "on",
					TSO4:     "on",
					TSO6:     "on",
					ECN:      "on",
					MrgRXBuf: "on",
				},
			}

			//iface.Model.Type = "virtio"

			//iface.MTU.Size = 9000

			log.Log.Infof("Mutated into %s", iface.Alias.Name)
		}
	}

	return nil
}

func lookupIfaceByAliasName(ifaces []libvirtxml.DomainInterface, name string) *libvirtxml.DomainInterface {
	for i, iface := range ifaces {
		if iface.Alias != nil && iface.Alias.Name == fmt.Sprintf("ua-%s", name) {
			log.Log.Infof("Found interface %s", name)
			return &ifaces[i]
		}
	}

	return nil
}
