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
		log.Log.Infof("Seems hook is running with preliminary domain spec (emulator: %s), skipping", domainSpec.Devices.Emulator)
		return &hooksV1alpha2.OnDefineDomainResult{
			DomainXML: params.GetDomainXML(),
		}, nil
	}

	if err := configMemory(vmiSpec.Spec.Domain.Memory, domainSpec.MemoryBacking); err != nil {
		return nil, fmt.Errorf("Failed to configure memory backing with hugepages and shared access: %s", err)
	}

	if err := configNetwork(vmiSpec.Spec.Domain.Devices.Interfaces, domainSpec.Devices.Interfaces); err != nil {
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

func havePCIControllers(controllers []libvirtxml.DomainController) bool {
	for _, device := range controllers {

		if device.Type == "pci" {
			return true
		}
	}
	return false
}
