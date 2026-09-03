package hook

import (
	"context"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"google.golang.org/grpc"
	"kubevirt.io/client-go/log"

	"kubevirt.io/kubevirt/pkg/hooks"
	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
	hooksV1alpha3 "kubevirt.io/kubevirt/pkg/hooks/v1alpha3"
)

const (
	onDefineDomainLoggingMessage  = "OnDefineDomain method has been called"
	preCloudInitIsoLoggingMessage = "PreCloudInitIso method has been called"
	onShutdownMessage             = "Hook's Shutdown callback method has been called"

	hookSocket = "dpdk.sock"
)

type v1Alpha3Server struct {
	done chan struct{}
}

func Hook(version string) {

	log.InitializeLogging("ovs")

	socketPath := filepath.Join(hooks.HookSocketsSharedDirectory, hookSocket)
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to initialized socket on path: %s", socketPath)
		log.Log.Error("Check whether given directory exists and socket name is not already taken by other file")
		os.Exit(1)
	}
	defer func() {
		if err := os.Remove(socketPath); err != nil {
			log.Log.Reason(err).Errorf("Failed to remove socket %s", socketPath)
		}
	}()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	hooksInfo.RegisterInfoServer(server, infoServer{Version: version})

	shutdownChan := make(chan struct{})

	hooksV1alpha3.RegisterCallbacksServer(server, v1Alpha3Server{done: shutdownChan})

	Serve(server, socket, shutdownChan)
}

func Serve(server *grpc.Server, socket net.Listener, shutdownChan <-chan struct{}) {
	errChan := make(chan error)
	go func() {
		errChan <- server.Serve(socket)
	}()
	waitForShutdown(server, errChan, shutdownChan)
}

func waitForShutdown(server *grpc.Server, errChan <-chan error, shutdownChan <-chan struct{}) {
	signalStopChan := make(chan os.Signal, 1)
	signal.Notify(signalStopChan, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	var err error
	select {
	case s := <-signalStopChan:
		log.Log.Infof("received signal: %s", s.String())
	case err = <-errChan:
		log.Log.Reason(err).Error("grpc server failed")
	case <-shutdownChan:
		log.Log.Info("received shutdown, exiting")
	}
	if err == nil {
		server.GracefulStop()
	}
}

func (s v1Alpha3Server) Shutdown(_ context.Context, _ *hooksV1alpha3.ShutdownParams) (*hooksV1alpha3.ShutdownResult, error) {
	log.Log.Info(onShutdownMessage)
	s.done <- struct{}{}
	return &hooksV1alpha3.ShutdownResult{}, nil
}
