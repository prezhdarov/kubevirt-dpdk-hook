package hook

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"

	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"
	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"
	"libvirt.org/go/libvirtxml"
)

const (
	OVSSocketDir = "/var/run/kubevirt"
)

func (s v1Alpha2Server) OnDefineDomain(ctx context.Context, params *hooksV1alpha2.OnDefineDomainParams) (*hooksV1alpha2.OnDefineDomainResult, error) {
	log.Log.Info(onDefineDomainLoggingMessage)

	vmiSpec := virtv1.VirtualMachineInstance{}
	if err := json.Unmarshal(params.GetVmi(), &vmiSpec); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal given VMI spec: %s due %v", params.GetVmi(), err)
	}

	domainSpec := libvirtxml.Domain{}
	if err := xml.Unmarshal(params.GetDomainXML(), &domainSpec); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal given Domain spec: %s %s", err, string(params.GetDomainXML()))
	}

	if havePCIControllers(domainSpec.Devices.Controllers) == false && domainSpec.Devices.Emulator == "" {
		return &hooksV1alpha2.OnDefineDomainResult{
			DomainXML: params.GetDomainXML(),
		}, nil
	}

	if err := configMemory(vmiSpec.Spec.Domain.Memory, domainSpec.MemoryBacking); err != nil {
		return nil, fmt.Errorf("Failed to configure memory backing with hugepages and shared access: %s", err)
	}

	// We don't care about the first version of the XML, do we?
	newDomainXML, err := xml.Marshal(domainSpec)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal new Domain spec: %s %+v", err, domainSpec)
	}

	return &hooksV1alpha2.OnDefineDomainResult{
		DomainXML: newDomainXML,
	}, nil
}

/*
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
						Path: filepath.Join(OVSSocketDir, fmt.Sprintf("%s.sock", strings.TrimPrefix(iface.Target.Dev, "tap"))),
						Mode: "server",
					},
				},
			}

			iface.Target.Managed = "yes"

			iface.Driver = &libvirtxml.DomainInterfaceDriver{
				Name:        "vhost",
				Queues:      2,
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
			//domainSpec.MemoryBacking.MemoryLocked = &libvirtxml.DomainMemoryLocked{}
			domainSpec.MemoryBacking.MemoryAccess = &libvirtxml.DomainMemoryAccess{Mode: "shared"}
		}

		newDomainXML, err := xml.Marshal(domainSpec)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal new Domain spec: %s %+v", err, domainSpec)
		}

		log.Log.Infof("new domain xml: %s", string(newDomainXML))

		return newDomainXML, nil
	}
*/
func havePCIControllers(controllers []libvirtxml.DomainController) bool {
	for _, device := range controllers {

		if device.Type == "pci" {
			return true
		}
	}
	return false
}
