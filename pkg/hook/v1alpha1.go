package hook

import (
	"context"

	"kubevirt.io/client-go/log"
	hooksV1alpha1 "kubevirt.io/kubevirt/pkg/hooks/v1alpha1"
)

func (s v1Alpha1Server) OnDefineDomain(ctx context.Context, params *hooksV1alpha1.OnDefineDomainParams) (*hooksV1alpha1.OnDefineDomainResult, error) {
	log.Log.Info(onDefineDomainLoggingMessage)
	newDomainXML, err := runOnDefineDomain(params.GetVmi(), params.GetDomainXML())
	if err != nil {
		log.Log.Reason(err).Error("Failed OnDefineDomain")
		return nil, err
	}
	return &hooksV1alpha1.OnDefineDomainResult{
		DomainXML: newDomainXML,
	}, nil
}
