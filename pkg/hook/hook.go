package hook

import (
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"google.golang.org/grpc"
	"kubevirt.io/client-go/log"

	"kubevirt.io/kubevirt/pkg/hooks"
	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"
)

const (
	onDefineDomainLoggingMessage  = "OnDefineDomain method has been called"
	preCloudInitIsoLoggingMessage = "PreCloudInitIso method has been called"
	onShutdownMessage             = "Hook's Shutdown callback method has been called"

	hookSocket = "dpdk.sock"

	onDefineDomainBin  = "onDefineDomain"
	preCloudInitIsoBin = "preCloudInitIso"
)

type v1Alpha2Server struct{}

func Hook(version string) {

	log.InitializeLogging("dpdk-hook")

	socketPath := filepath.Join(hooks.HookSocketsSharedDirectory, hookSocket)
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to initialized socket on path: %s", socket)
		log.Log.Error("Check whether given directory exists and socket name is not already taken by other file")
		os.Exit(1)
	}
	defer os.Remove(socketPath)

	server := grpc.NewServer([]grpc.ServerOption{}...)
	hooksInfo.RegisterInfoServer(server, infoServer{Version: version})
	hooksV1alpha2.RegisterCallbacksServer(server, v1Alpha2Server{})

	// Handle signals to properly shutdown process
	signalStopChan := make(chan os.Signal, 1)
	signal.Notify(signalStopChan, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	log.Log.Infof("shim is now exposing its services on socket %s", socketPath)
	errChan := make(chan error)
	go func() {
		errChan <- server.Serve(socket)
	}()

	select {
	case s := <-signalStopChan:
		log.Log.Infof("dpdk-hook received signal: %s", s.String())
	case err = <-errChan:
		log.Log.Reason(err).Error("Failed to run grpc server")

	}

	if err == nil {
		server.GracefulStop()
	}

}
