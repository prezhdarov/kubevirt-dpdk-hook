package domain

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"

	vmschema "kubevirt.io/api/core/v1"

	domainschema "kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
)

type NetworkConfiguratorOptions struct {
	IstioProxyInjectionEnabled bool
	UseVirtioTransitional      bool
}

type OVSNetworkConfigurator struct {
	vmiSpecIface []*vmschema.Interface
	hugePageSize string
	options      NetworkConfiguratorOptions
}

const (
	// PasstPluginName passt binding plugin name should be registered to Kubevirt through Kubevirt CR
	OVSPluginName = "ovs"
	// PasstLogFilePath passt log file path Kubevirt consume and record
	OVSLogFilePath = "/var/run/kubevirt/ovs.log"

	OVSSocketDir = "/var/run/kubevirt/vh"

	OVSVHostUserDirectory = "vhostuser"
)

func NewOVSNetworkConfigurator(ifaces []vmschema.Interface, networks []vmschema.Network, memory vmschema.Memory, opts NetworkConfiguratorOptions) (*OVSNetworkConfigurator, error) {

	var (
		newIfaces []*vmschema.Interface
		pageSize  string
	)

	for _, iface := range ifaces {
		if iface.Binding == nil || iface.Binding != nil && iface.Binding.Name != OVSPluginName {
			return nil, fmt.Errorf("interface %q is not set with ovs network binding plugin", iface.Name)
		}
		newIfaces = append(newIfaces, &iface)

	}

	if memory.Hugepages != nil {
		pageSize = memory.Hugepages.PageSize
	}

	return &OVSNetworkConfigurator{
		vmiSpecIface: newIfaces,
		hugePageSize: pageSize,
		options:      opts,
	}, nil
}

func (p OVSNetworkConfigurator) Mutate(domainSpec *domainschema.DomainSpec) (*domainschema.DomainSpec, error) {
	const (
		sharedMemoryBackingAccessMode = "shared"
		memfdMemoryBackingSourceType  = "memfd"
	)

	domainSpecCopy := domainSpec.DeepCopy()

	if domainSpecCopy.MemoryBacking != nil &&
		domainSpecCopy.MemoryBacking.Access != nil &&
		domainSpecCopy.MemoryBacking.Access.Mode != sharedMemoryBackingAccessMode {
		return nil, fmt.Errorf("memory backing access mode must be 'shared'; cannot override existing mode: %q",
			domainSpec.MemoryBacking.Access.Mode)
	}

	if domainSpecCopy.MemoryBacking == nil {
		domainSpecCopy.MemoryBacking = &domainschema.MemoryBacking{
			Access: &domainschema.MemoryBackingAccess{
				Mode: sharedMemoryBackingAccessMode,
			},
			Source: &domainschema.MemoryBackingSource{
				Type: memfdMemoryBackingSourceType,
			},
		}
	}

	if p.hugePageSize != "" {

		ugePage, err := hugepageFromVMI(p.hugePageSize)
		if err != nil {
			return nil, err
		}
		domainSpecCopy.MemoryBacking.HugePages.HugePage = append(domainSpecCopy.MemoryBacking.HugePages.HugePage, *ugePage)
	}

	return domainSpecCopy, nil
}

func lookupIfaceByAliasName(ifaces []domainschema.Interface, name string) *domainschema.Interface {
	for i, iface := range ifaces {
		if iface.Alias != nil && iface.Alias.GetName() == name {
			return &ifaces[i]
		}
	}

	return nil
}

func (p OVSNetworkConfigurator) mutateInterface(ifaceName string) (*domainschema.Interface, error) {
	return nil, nil
}

/*
	func (p OVSNetworkConfigurator) generateInterface() (*domainschema.Interface, error) {
		var pciAddress *domainschema.Address
		if p.vmiSpecIface.PciAddress != "" {
			var err error
			pciAddress, err = device.NewPciAddressField(p.vmiSpecIface.PciAddress)
			if err != nil {
				return nil, err
			}
		}

		var ifaceModel string
		if p.vmiSpecIface.Model == "" {
			ifaceModel = vmschema.VirtIO
		} else {
			ifaceModel = p.vmiSpecIface.Model
		}

		var ifaceModelType string
		if ifaceModel == vmschema.VirtIO {
			if p.options.UseVirtioTransitional {
				ifaceModelType = "virtio-transitional"
			} else {
				ifaceModelType = "virtio-non-transitional"
			}
		} else {
			ifaceModelType = p.vmiSpecIface.Model
		}
		model := &domainschema.Model{Type: ifaceModelType}

		var mac *domainschema.MAC
		if p.vmiSpecIface.MacAddress != "" {
			mac = &domainschema.MAC{MAC: p.vmiSpecIface.MacAddress}
		}

		var acpi *domainschema.ACPI
		if p.vmiSpecIface.ACPIIndex > 0 {
			acpi = &domainschema.ACPI{Index: uint(p.vmiSpecIface.ACPIIndex)}
		}

		const (
			ifaceTypeUser     = "user"
			ifaceBackendPasst = "passt"
		)
		return &domainschema.Interface{
			Alias:       domainschema.NewUserDefinedAlias(p.vmiSpecIface.Name),
			Model:       model,
			Address:     pciAddress,
			MAC:         mac,
			ACPI:        acpi,
			Type:        ifaceTypeUser,
			Source:      domainschema.InterfaceSource{Device: namescheme.PrimaryPodInterfaceName},
			Backend:     &domainschema.InterfaceBackend{Type: ifaceBackendPasst, LogFile: OVSLogFilePath},
			PortForward: p.generatePortForward(),
		}, nil
	}

	func (p OVSNetworkConfigurator) generatePortForward() []domainschema.InterfacePortForward {
		var tcpPortsRange, udpPortsRange []domainschema.InterfacePortForwardRange

		if p.options.IstioProxyInjectionEnabled {
			for _, port := range istio.ReservedPorts() {
				tcpPortsRange = append(tcpPortsRange, domainschema.InterfacePortForwardRange{Start: uint(port), Exclude: "yes"})
			}
		}

		const (
			protoTCP = "tcp"
			protoUDP = "udp"
		)

		for _, port := range p.vmiSpecIface.Ports {
			if strings.EqualFold(port.Protocol, protoTCP) || port.Protocol == "" {
				tcpPortsRange = append(tcpPortsRange, domainschema.InterfacePortForwardRange{Start: uint(port.Port)})
			} else if strings.EqualFold(port.Protocol, protoUDP) {
				udpPortsRange = append(udpPortsRange, domainschema.InterfacePortForwardRange{Start: uint(port.Port)})
			} else {
				log.Log.Errorf("protocol %s is not supported by passt", port.Protocol)
			}
		}

		var portsFwd []domainschema.InterfacePortForward
		if len(udpPortsRange) == 0 && len(tcpPortsRange) == 0 {
			portsFwd = append(portsFwd, domainschema.InterfacePortForward{Proto: protoTCP})
			portsFwd = append(portsFwd, domainschema.InterfacePortForward{Proto: protoUDP})
		}
		if len(tcpPortsRange) > 0 {
			portsFwd = append(portsFwd, domainschema.InterfacePortForward{Proto: protoTCP, Ranges: tcpPortsRange})
		}
		if len(udpPortsRange) > 0 {
			portsFwd = append(portsFwd, domainschema.InterfacePortForward{Proto: protoUDP, Ranges: udpPortsRange})
		}

		return portsFwd
	}
*/

func hugepageFromVMI(pagesize string) (*domainschema.HugePage, error) {

	var pagesizeRegex = regexp.MustCompile(`^(\d+)([A-Za-z]+)$`)

	pagesizeMatch := pagesizeRegex.FindStringSubmatch(pagesize)
	if len(pagesizeMatch) != 3 {
		return &domainschema.HugePage{}, fmt.Errorf("invalid pagesize: %s", pagesize)
	}

	//size, err := strconv.ParseUint(pagesizeMatch[1], 10, 64)
	//if err != nil {
	//	return &domainschema.HugePage{}, err
	//}

	return &domainschema.HugePage{
		//Size: uint(size),
		Size: pagesizeMatch[1],
		Unit: pagesizeMatch[2] + "B",
	}, nil
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
