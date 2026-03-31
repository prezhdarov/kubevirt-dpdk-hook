package callback

import (
	"encoding/xml"
	"fmt"

	"libvirt.org/go/libvirtxml"
)

// TODO: move to Kubevirt domain API package
const libvirtDomainQemuSchema = "http://libvirt.org/schemas/domain/qemu/1.0"

type DomainSpecMutator interface {
	Mutate(*libvirtxml.Domain) error
}

func OnDefineDomain(domainXML []byte, domSpecMutator DomainSpecMutator) ([]byte, error) {
	//domainSpec := &domainschema.DomainSpec{
	//	// Unmarshalling domain spec makes the XML namespace attribute empty.
	//	// Some domain parameters requires namespace to be defined.
	//	// e.g: https://libvirt.org/drvqemu.html#pass-through-of-arbitrary-qemu-commands
	//	XmlNS: libvirtDomainQemuSchema,
	//}

	domainSpec := &libvirtxml.Domain{}

	if err := xml.Unmarshal(domainXML, domainSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal given domain spec: %v", err)
	}

	err := domSpecMutator.Mutate(domainSpec)
	if err != nil {
		return nil, err
	}

	updatedDomainSpecXML, err := xml.Marshal(domainSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated domain spec: %v", err)
	}

	return updatedDomainSpecXML, nil
}
