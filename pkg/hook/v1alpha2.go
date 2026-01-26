package hook

import (
	"context"

	"kubevirt.io/client-go/log"
	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"
)

func (s v1Alpha2Server) OnDefineDomain(ctx context.Context, params *hooksV1alpha2.OnDefineDomainParams) (*hooksV1alpha2.OnDefineDomainResult, error) {
	log.Log.Info(onDefineDomainLoggingMessage)
	newDomainXML, err := runOnDefineDomain(params.GetVmi(), params.GetDomainXML())
	if err != nil {
		log.Log.Reason(err).Error("Failed OnDefineDomain")
		return nil, err
	}
	return &hooksV1alpha2.OnDefineDomainResult{
		DomainXML: newDomainXML,
	}, nil
}

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
