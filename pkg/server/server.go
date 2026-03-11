package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prezhdarov/kubevirt-dpdk-hook/pkg/callback"
	"github.com/prezhdarov/kubevirt-dpdk-hook/pkg/domain"
	vmschema "kubevirt.io/api/core/v1"

	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"
)

type InfoServer struct {
	Version string
}

func (s InfoServer) Info(_ context.Context, _ *hooksInfo.InfoParams) (*hooksInfo.InfoResult, error) {
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
		},
	}, nil
}

type V1alpha2Server struct{}

func (s V1alpha2Server) OnDefineDomain(_ context.Context, params *hooksV1alpha2.OnDefineDomainParams) (*hooksV1alpha2.OnDefineDomainResult, error) {
	vmi := &vmschema.VirtualMachineInstance{}
	if err := json.Unmarshal(params.GetVmi(), vmi); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VMI: %v", err)
	}

	useVirtioTransitional := vmi.Spec.Domain.Devices.UseVirtioTransitional != nil && *vmi.Spec.Domain.Devices.UseVirtioTransitional

	const istioInjectAnnotation = "sidecar.istio.io/inject"
	istioProxyInjectionEnabled := false
	if val, ok := vmi.GetAnnotations()[istioInjectAnnotation]; ok {
		istioProxyInjectionEnabled = strings.ToLower(val) == "true"
	}

	opts := domain.NetworkConfiguratorOptions{
		UseVirtioTransitional:      useVirtioTransitional,
		IstioProxyInjectionEnabled: istioProxyInjectionEnabled,
	}

	ovsConfigurator, err := domain.NewOVSNetworkConfigurator(vmi.Spec.Domain.Devices.Interfaces, vmi.Spec.Networks, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create passt configurator: %v", err)
	}

	newDomainXML, err := callback.OnDefineDomain(params.GetDomainXML(), ovsConfigurator)
	if err != nil {
		return nil, err
	}

	return &hooksV1alpha2.OnDefineDomainResult{
		DomainXML: newDomainXML,
	}, nil
}

func (s V1alpha2Server) PreCloudInitIso(_ context.Context, params *hooksV1alpha2.PreCloudInitIsoParams) (*hooksV1alpha2.PreCloudInitIsoResult, error) {
	return &hooksV1alpha2.PreCloudInitIsoResult{
		CloudInitData: params.GetCloudInitData(),
	}, nil
}
