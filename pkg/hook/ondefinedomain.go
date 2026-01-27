package hook

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"
	"libvirt.org/go/libvirtxml"
)

func runOnDefineDomain(vmiJSON []byte, domainXML []byte) ([]byte, error) {
	//if _, err := exec.LookPath(onDefineDomainBin); err != nil {
	//	return nil, fmt.Errorf("Failed in finding %s in $PATH due %v", onDefineDomainBin, err)
	//}

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

		log.Log.Infof("Interface %s of type %s", iface.XMLName.Local, iface.Model.Type)

	}

	//log.Log.Infof("vmi json: %s", string(vmiJSON))
	//log.Log.Infof("domain xml: %s", string(domainXML))

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
