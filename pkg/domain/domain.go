package domain

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	vmschema "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"

	"kubevirt.io/kubevirt/pkg/hooks"
	domainschema "kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"

	"kubevirt.io/kubevirt/pkg/network/vmispec"
)

type NetworkConfiguratorOptions struct {
	IstioProxyInjectionEnabled bool
	UseVirtioTransitional      bool
}

type OVSNetworkConfigurator struct {
	vmiSpecIface *vmschema.Interface
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

func NewOVSNetworkConfigurator(ifaces []vmschema.Interface, networks []vmschema.Network, opts NetworkConfiguratorOptions) (*OVSNetworkConfigurator, error) {
	network := vmispec.LookupPodNetwork(networks)
	if network == nil {
		return nil, fmt.Errorf("pod network not found: %d, %s, %+v", len(networks), networks[0].Name, networks[0].Multus)
	}
	iface := vmispec.LookupInterfaceByName(ifaces, network.Name)
	if iface == nil {
		return nil, fmt.Errorf("no interface found")
	}
	if iface.Binding == nil || iface.Binding != nil && iface.Binding.Name != OVSPluginName {
		return nil, fmt.Errorf("interface %q is not set with Passt network binding plugin", network.Name)
	}

	return &OVSNetworkConfigurator{
		vmiSpecIface: iface,
		options:      opts,
	}, nil
}

func (p OVSNetworkConfigurator) Mutate(domainSpec *domainschema.DomainSpec) (*domainschema.DomainSpec, error) {
	//generatedIface, err := p.generateInterface()
	//if err != nil {
	//	return nil, fmt.Errorf("failed to generate domain interface spec: %v", err)
	//}

	/*
		socketPath := filepath.Join(hooks.HookSocketsSharedDirectory, OVSVHostUserDirectory)

		exists, err := exists(socketPath)
		if err != nil {
			log.Log.Warningf("Could not check if directory exists: %s", socketPath)
		}

		if !exists {
			if err := os.Mkdir(socketPath, os.ModePerm); err != nil {
				log.Log.Warningf("Could not create directory: %s", socketPath)
			}
		}
	*/

	domainSpecCopy := domainSpec.DeepCopy()
	if iface := lookupIfaceByAliasName(domainSpecCopy.Devices.Interfaces, p.vmiSpecIface.Name); iface != nil {

		filePath := filepath.Join(hooks.HookSocketsSharedDirectory, OVSVHostUserDirectory)
		//socketPath := filepath.Join(OVSSocketDir, fmt.Sprintf("%s.sock", p.vmiSpecIface.Name))
		//	*iface = *generatedIface

		exists, err := exists(filePath)
		if err != nil {
			log.Log.Warningf("Could not check if directory exists: %s", filePath)
		}

		if !exists {
			if err := os.Mkdir(filePath, os.ModePerm); err != nil {
				log.Log.Warningf("Could not create directory: %s", filePath)
			}
		}

		socketPath := filepath.Join(filePath, p.vmiSpecIface.Name)
		os.OpenFile(socketPath, os.O_RDONLY|os.O_CREATE, 0666)

		log.Log.Infof("ovs interface is NOT added to domain spec successfully: %+v", iface)
	} else {
		//	domainSpecCopy.Devices.Interfaces = append(domainSpecCopy.Devices.Interfaces, *generatedIface)
	}

	//

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
