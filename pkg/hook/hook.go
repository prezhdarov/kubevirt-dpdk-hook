package hook

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"kubevirt.io/client-go/log"

	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
	hooksV1alpha1 "kubevirt.io/kubevirt/pkg/hooks/v1alpha1"
	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"
	hooksV1alpha3 "kubevirt.io/kubevirt/pkg/hooks/v1alpha3"
)

const (
	onDefineDomainLoggingMessage  = "OnDefineDomain method has been called"
	preCloudInitIsoLoggingMessage = "PreCloudInitIso method has been called"
	onShutdownMessage             = "Hook's Shutdown callback method has been called"

	onDefineDomainBin  = "onDefineDomain"
	preCloudInitIsoBin = "preCloudInitIso"
)

type v1Alpha1Server struct{}
type v1Alpha2Server struct{}
type v1Alpha3Server struct {
	done chan struct{}
}

func Hook(version string) {

	log.InitializeLogging("shim-sidecar")

	socketPath, err := getSocketPath()
	if err != nil {
		log.Log.Reason(err).Errorf("Enviroment error")
		os.Exit(1)
	}

	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to initialized socket on path: %s", socket)
		os.Exit(1)
	}
	defer os.Remove(socketPath)

	server := grpc.NewServer([]grpc.ServerOption{}...)
	hooksInfo.RegisterInfoServer(server, infoServer{Version: version})
	hooksV1alpha1.RegisterCallbacksServer(server, v1Alpha1Server{})
	hooksV1alpha2.RegisterCallbacksServer(server, v1Alpha2Server{})

	shutdownChan := make(chan struct{})
	hooksV1alpha3.RegisterCallbacksServer(server, v1Alpha3Server{done: shutdownChan})

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
		log.Log.Infof("sidecar-shim received signal: %s", s.String())
	case err = <-errChan:
		log.Log.Reason(err).Error("Failed to run grpc server")
	case <-shutdownChan:
		log.Log.Info("Exiting")
	}

	if err == nil {
		server.GracefulStop()
	}

}
