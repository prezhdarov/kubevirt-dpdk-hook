package main

import (
	netbind "github.com/prezhdarov/kubevirt-dpdk-hook/pkg/netbind"
)

func main() {
	netbind.NetBind("v1alpha2")
}
