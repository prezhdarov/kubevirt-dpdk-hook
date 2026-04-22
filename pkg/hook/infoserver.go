package hook

import (
	"context"
	"errors"
	"os/exec"

	"kubevirt.io/client-go/log"
	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
)

type infoServer struct {
	Version string
}

func (s infoServer) Info(ctx context.Context, params *hooksInfo.InfoParams) (*hooksInfo.InfoResult, error) {
	log.Log.Info("Info method has been called")
	supportedHookPoints := map[string]string{
		hooksInfo.OnDefineDomainHookPointName:  onDefineDomainBin,
		hooksInfo.PreCloudInitIsoHookPointName: preCloudInitIsoBin,
	}
	var hookPoints = []*hooksInfo.HookPoint{}

	// Shutdown fixes proper termination of Sidecars. It isn't related to
	// user's binaries nor scripts.
	if s.Version != "v1alpha1" && s.Version != "v1alpha2" {
		hookPoints = append(hookPoints, &hooksInfo.HookPoint{
			Name:     hooksInfo.ShutdownHookPointName,
			Priority: 0,
		})
	}

	for hookPointName, binName := range supportedHookPoints {
		if _, err := exec.LookPath(binName); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				log.Log.Infof("Info: %s has not been found", binName)
			}
			continue
		}

		hookPoints = append(hookPoints, &hooksInfo.HookPoint{
			Name:     hookPointName,
			Priority: 0,
		})
	}

	return &hooksInfo.InfoResult{
		Name: "shim",
		Versions: []string{
			s.Version,
		},
		HookPoints: hookPoints,
	}, nil
}
