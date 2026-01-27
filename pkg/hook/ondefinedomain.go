package hook

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"

	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"
	"libvirt.org/go/libvirtxml"
)

func localAddresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return
	}

	for _, i := range ifaces {
		//addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		}
		log.Log.Infof("%s => %s", i.Name, i.HardwareAddr)
	}
}

func runOnDefineDomain(vmiJSON []byte, domainXML []byte) ([]byte, error) {
	//if _, err := exec.LookPath(onDefineDomainBin); err != nil {
	//	return nil, fmt.Errorf("Failed in finding %s in $PATH due %v", onDefineDomainBin, err)
	//}

	log.Log.Infof("vmi json: %s", string(vmiJSON))

	log.Log.Infof("domain xml: %s", string(domainXML))

	vmiSpec := virtv1.VirtualMachineInstance{}
	if err := json.Unmarshal(vmiJSON, &vmiSpec); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal given VMI spec: %s due %v", vmiJSON, err)
	}

	domainSpec := libvirtxml.Domain{}
	if err := xml.Unmarshal(domainXML, &domainSpec); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal given Domain spec: %s %s", err, string(domainXML))
	}

	annotations := vmiSpec.GetAnnotations()

	for _, iface := range domainSpec.Devices.Interfaces {
		//if iface.Source.Type == "vhostuser" {
		//	vhostInterfaces = append(vhostInterfaces, iface)

		fmt.Printf("\n🔗 vHost User Interface Found:\n")
		//	fmt.Printf("   Socket Path: %s\n", iface.Source.Path)
		//	fmt.Printf("   Mode: %s\n", iface.Source.Mode)

		if iface.MAC != nil {
			fmt.Printf("   MAC: %s\n", iface.MAC.Address)
		}

		if iface.MTU != nil {
			fmt.Printf("   MTU: %d\n", iface.MTU.Size)
		}

		if iface.Model != nil {
			fmt.Printf("   Model: %s\n", iface.Model.Type)
		}

		if iface.Target != nil {
			fmt.Printf("   Target: %s (managed: %s\n", iface.Target.Dev, iface.Target.Managed)
		}

		if iface.Alias != nil {
			fmt.Printf("   Alias: %s\n", iface.Alias.Name)
		}

		if iface.Driver != nil {
			fmt.Printf("   Driver: %s\n", iface.Driver.Name)
			fmt.Printf("   Queues: %d\n", iface.Driver.Queues)
		}

		if iface.Address != nil {
			fmt.Printf("   PCI: %s:%s:%s.%s\n",
				iface.Address.PCI.Domain,
				iface.Address.PCI.Bus,
				iface.Address.PCI.Slot,
				iface.Address.PCI.Function)
		}

		localAddresses()
		//}
	}

	log.Log.Infof("VMI annotations", annotations)

	//args := append([]string{},
	//	"--vmi", string(vmiJSON),
	//	"--domain", string(domainXML))

	//log.Log.Infof("Executing %s", onDefineDomainBin)
	//command := exec.Command(onDefineDomainBin, args...)
	//if reader, err := command.StderrPipe(); err != nil {
	//	log.Log.Reason(err).Infof("Could not pipe stderr")
	//} else {
	//	go logStderr(reader, "onDefineDomain")
	//}
	//return command.Output()

	return domainXML, nil
}
