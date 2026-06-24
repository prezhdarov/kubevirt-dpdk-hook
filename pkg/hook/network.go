package hook

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	vmschema "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"
	"libvirt.org/go/libvirtxml"
)

const dpdkAnnotationRoot = "dpdk.kubevirt.io"

func configNetwork(interfaces []vmschema.Interface, domainInterfaces []libvirtxml.DomainInterface, driverSettings *libvirtxml.DomainInterfaceDriver) error {

	for _, iface := range interfaces {

		if iface := lookupIfaceByAliasName(domainInterfaces, iface.Name); iface != nil {
			iface.Target.Managed = "yes"

			iface.Source = &libvirtxml.DomainInterfaceSource{
				VHostUser: &libvirtxml.DomainInterfaceSourceVHostUser{
					Chardev: &libvirtxml.DomainChardevSource{
						UNIX: &libvirtxml.DomainChardevSourceUNIX{
							Path: filepath.Join(OVSSocketDir, fmt.Sprintf("%s.sock", strings.TrimPrefix(iface.Target.Dev, "tap"))),
							Mode: "server",
						},
					},
				},
			}

			iface.Driver = driverSettings

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

func getDPDKDriverSettings(annotations map[string]string) *libvirtxml.DomainInterfaceDriver {

	driverSettings := &libvirtxml.DomainInterfaceDriver{
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
			UFO:      "on",
			MrgRXBuf: "on",
		},
	}

	for annotation, value := range annotations {

		switch annotation {
		case dpdkAnnotationRoot + "/queues":
			setDriverUInt(annotation, value, &driverSettings.Queues)
		case dpdkAnnotationRoot + "/rx-queue-size":
			setDriverUInt(annotation, value, &driverSettings.RXQueueSize)
		case dpdkAnnotationRoot + "/tx-queue-size":
			setDriverUInt(annotation, value, &driverSettings.TXQueueSize)
		case dpdkAnnotationRoot + "/guest-csum":
			setDriverString(value, &driverSettings.Guest.CSum)
		case dpdkAnnotationRoot + "/guest-tso4":
			setDriverString(value, &driverSettings.Guest.TSO4)
		case dpdkAnnotationRoot + "/guest-tso6":
			setDriverString(value, &driverSettings.Guest.TSO6)
		case dpdkAnnotationRoot + "/guest-ecn":
			setDriverString(value, &driverSettings.Guest.ECN)
		case dpdkAnnotationRoot + "/host-csum":
			setDriverString(value, &driverSettings.Host.CSum)
		case dpdkAnnotationRoot + "/host-gso":
			setDriverString(value, &driverSettings.Host.GSO)
		case dpdkAnnotationRoot + "/host-tso4":
			setDriverString(value, &driverSettings.Host.TSO4)
		case dpdkAnnotationRoot + "/host-tso6":
			setDriverString(value, &driverSettings.Host.TSO6)
		case dpdkAnnotationRoot + "/host-ecn":
			setDriverString(value, &driverSettings.Host.ECN)
		case dpdkAnnotationRoot + "/host-ufo":
			setDriverString(value, &driverSettings.Host.UFO)
		case dpdkAnnotationRoot + "/host-mrg-rxbuf":
			setDriverString(value, &driverSettings.Host.MrgRXBuf)
		}

	}

	return driverSettings
}

func setDriverUInt(annotation, value string, driverValue *uint) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Log.Infof("%s: %s cannot convert to integer, skipping.. ", annotation, value)
	} else {
		*driverValue = uint(intValue)
	}
}

func setDriverString(annotationValue string, driverValue *string) {
	if (annotationValue == "on" || annotationValue == "off") && annotationValue != *driverValue {
		*driverValue = annotationValue
	}
}
