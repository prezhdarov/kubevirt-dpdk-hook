package hook

import (
	"encoding/json"
	"fmt"

	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"
)

func runOnDefineDomain(vmiJSON []byte, domainXML []byte) ([]byte, error) {
	//if _, err := exec.LookPath(onDefineDomainBin); err != nil {
	//	return nil, fmt.Errorf("Failed in finding %s in $PATH due %v", onDefineDomainBin, err)
	//}

	vmiSpec := virtv1.VirtualMachineInstance{}
	if err := json.Unmarshal(vmiJSON, &vmiSpec); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal given VMI spec: %s due %v", vmiJSON, err)
	}

	log.Log.Infof("vmi json: %s", string(vmiJSON))
	log.Log.Infof("domain xml: %s", string(domainXML))

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
