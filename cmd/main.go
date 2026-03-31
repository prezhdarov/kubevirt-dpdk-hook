package main

import (
	"net"
	"os"
	"path/filepath"

	srv "github.com/prezhdarov/kubevirt-dpdk-hook/pkg/server"

	"google.golang.org/grpc"

	"kubevirt.io/kubevirt/pkg/hooks"
	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"

	"kubevirt.io/client-go/log"
)

const hookSocket = "ovs.sock"

func main() {

	log.InitializeLogging("ovs-plugin")

	socketPath := filepath.Join(hooks.HookSocketsSharedDirectory, hookSocket)
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to initialized socket on path: %s", socket)
		log.Log.Error("Check whether given directory exists and socket name is not already taken by other file")
		os.Exit(1)
	}
	defer os.Remove(socketPath)

	server := grpc.NewServer([]grpc.ServerOption{}...)
	hooksInfo.RegisterInfoServer(server, srv.InfoServer{Version: "v1alpha2"})
	hooksV1alpha2.RegisterCallbacksServer(server, srv.V1alpha2Server{})

	log.Log.Infof("Starting hook server exposing 'info' and '%s' services on socket %q", socketPath, "v1alpha2")
	server.Serve(socket)
}
