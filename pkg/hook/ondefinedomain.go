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

	var newInterfaces []libvirtxml.DomainInterface

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

	log.Log.Info("VMI annotations detected")
	for annotation, value := range annotations {
		log.Log.Info(fmt.Sprintf("%s: %s", annotation, value))
	}
	// We don't care about the first version of the XML, do we?
	if havePCIControllers(domainSpec.Devices.Controllers) == false && domainSpec.Devices.Emulator == "" {
		return domainXML, nil
	}

	for _, iface := range domainSpec.Devices.Interfaces {

		iface.Source = &libvirtxml.DomainInterfaceSource{
			VHostUser: &libvirtxml.DomainChardevSource{
				UNIX: &libvirtxml.DomainChardevSourceUNIX{
					Path: "/tmp/vhost-user1.sock",
					Mode: "client", // or "server"
				},
			},
		}

		iface.Target.Managed = "yes"

		newInterfaces = append(newInterfaces, iface)
	}

	domainSpec.Devices.Interfaces = newInterfaces

	if vmiSpec.Spec.Domain.Memory != nil &&
		vmiSpec.Spec.Domain.Memory.Hugepages != nil &&
		vmiSpec.Spec.Domain.Memory.Hugepages.PageSize != "" {

		ugePage, err := hugepageFromVMI(vmiSpec.Spec.Domain.Memory.Hugepages.PageSize)
		if err != nil {
			return nil, err
		}

		domainSpec.MemoryBacking.MemoryHugePages.Hugepages = append(domainSpec.MemoryBacking.MemoryHugePages.Hugepages, ugePage)
		domainSpec.MemoryBacking.MemoryLocked = &libvirtxml.DomainMemoryLocked{}
	}

	newDomainXML, err := xml.Marshal(domainSpec)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal new Domain spec: %s %+v", err, domainSpec)
	}

	log.Log.Infof("new domain xml: %s", string(newDomainXML))

	return newDomainXML, nil
}
