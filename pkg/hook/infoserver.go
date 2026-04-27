package hook

import (
	"context"

	"kubevirt.io/client-go/log"
	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
)

type infoServer struct {
	Version string
}

func (s infoServer) Info(ctx context.Context, params *hooksInfo.InfoParams) (*hooksInfo.InfoResult, error) {
	log.Log.Info("Info method has been called")

	return &hooksInfo.InfoResult{
		Name: "ovs",
		Versions: []string{
			s.Version,
		},
		HookPoints: []*hooksInfo.HookPoint{
			{
				Name:     hooksInfo.OnDefineDomainHookPointName,
				Priority: 0,
			},
			{
				Name:     hooksInfo.ShutdownHookPointName,
				Priority: 0,
			},
		},
	}, nil
}
