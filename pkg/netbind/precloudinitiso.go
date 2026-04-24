package netbind

import (
	"context"

	hooksV1alpha3 "kubevirt.io/kubevirt/pkg/hooks/v1alpha3"
)

func (s v1Alpha3Server) PreCloudInitIso(_ context.Context, params *hooksV1alpha3.PreCloudInitIsoParams) (*hooksV1alpha3.PreCloudInitIsoResult, error) {
	return &hooksV1alpha3.PreCloudInitIsoResult{
		CloudInitData: params.GetCloudInitData(),
	}, nil
}
