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
