package hook

import (
	"context"

	"kubevirt.io/client-go/log"
	hooksV1alpha3 "kubevirt.io/kubevirt/pkg/hooks/v1alpha3"
)

func (s v1Alpha3Server) OnDefineDomain(_ context.Context, params *hooksV1alpha3.OnDefineDomainParams) (*hooksV1alpha3.OnDefineDomainResult, error) {
	log.Log.Info(onDefineDomainLoggingMessage)
	newDomainXML, err := runOnDefineDomain(params.GetVmi(), params.GetDomainXML())
	if err != nil {
		log.Log.Reason(err).Error("Failed OnDefineDomain")
		return nil, err
	}
	return &hooksV1alpha3.OnDefineDomainResult{
		DomainXML: newDomainXML,
	}, nil
}

func (s v1Alpha3Server) PreCloudInitIso(_ context.Context, params *hooksV1alpha3.PreCloudInitIsoParams) (*hooksV1alpha3.PreCloudInitIsoResult, error) {
	log.Log.Info(preCloudInitIsoLoggingMessage)
	cloudInitData, err := runPreCloudInitIso(params.GetVmi(), params.GetCloudInitData())
	if err != nil {
		log.Log.Reason(err).Error("Failed ProCloudInitIso")
		return nil, err
	}
	return &hooksV1alpha3.PreCloudInitIsoResult{
		CloudInitData: cloudInitData,
	}, nil
}

func (s v1Alpha3Server) Shutdown(_ context.Context, _ *hooksV1alpha3.ShutdownParams) (*hooksV1alpha3.ShutdownResult, error) {
	log.Log.Info(onShutdownMessage)
	s.done <- struct{}{}
	return &hooksV1alpha3.ShutdownResult{}, nil
}
