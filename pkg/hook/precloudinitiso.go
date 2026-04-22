package hook

import (
	"context"

	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"
)

func (s v1Alpha2Server) PreCloudInitIso(_ context.Context, params *hooksV1alpha2.PreCloudInitIsoParams) (*hooksV1alpha2.PreCloudInitIsoResult, error) {
	return &hooksV1alpha2.PreCloudInitIsoResult{
		CloudInitData: params.GetCloudInitData(),
	}, nil
}

/*
func (s v1Alpha2Server) PreCloudInitIso(_ context.Context, params *hooksV1alpha2.PreCloudInitIsoParams) (*hooksV1alpha2.PreCloudInitIsoResult, error) {
	log.Log.Info(preCloudInitIsoLoggingMessage)
	cloudInitData, err := runPreCloudInitIso(params.GetVmi(), params.GetCloudInitData())
	if err != nil {
		log.Log.Reason(err).Error("Failed ProCloudInitIso")
		return nil, err
	}
	return &hooksV1alpha2.PreCloudInitIsoResult{
		CloudInitData: cloudInitData,
	}, nil
}

func runPreCloudInitIso(vmiJSON []byte, cloudInitDataJSON []byte) ([]byte, error) {
	// Check binary exists
	//if _, err := exec.LookPath(preCloudInitIsoBin); err != nil {
	//	return nil, fmt.Errorf("Failed in finding %s in $PATH: %v", preCloudInitIsoBin, err)
	//}

	// Validate params before calling hook script
	vmiSpec := virtv1.VirtualMachineInstance{}
	if err := json.Unmarshal(vmiJSON, &vmiSpec); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal given VMI spec: %s due %v", vmiJSON, err)
	}

	cloudInitData := cloudinit.CloudInitData{}
	err := json.Unmarshal(cloudInitDataJSON, &cloudInitData)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal given CloudInitData: %s due %v", cloudInitDataJSON, err)
	}

	//args := append([]string{},
	//	"--vmi", string(vmiJSON),
	//	"--cloud-init", string(cloudInitDataJSON))

	//log.Log.Infof("Executing %s", preCloudInitIsoBin)
	//command := exec.Command(preCloudInitIsoBin, args...)
	//if reader, err := command.StderrPipe(); err != nil {
	//	log.Log.Reason(err).Infof("Could not pipe stderr")
	//} else {
	//	go logStderr(reader, "cloudInitData")
	//}
	return cloudInitDataJSON, err
}
*/
