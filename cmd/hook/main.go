package main

import (
	hook "github.com/prezhdarov/kubevirt-dpdk-hook/pkg/hook"
)

func main() {
	hook.Hook("v1alpha2")
}
