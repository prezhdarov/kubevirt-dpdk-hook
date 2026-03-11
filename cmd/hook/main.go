package main

import "github.com/prezhdarov/kubevirt-dpdk-hook/pkg/hook"

func main() {

	var version = "v1alpha2"

	hook.Hook(version)

}
