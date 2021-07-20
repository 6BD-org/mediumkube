package main

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/pkg/common"
	"mediumkube/pkg/common/event"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/daemon/mesh"
	"mediumkube/pkg/daemon/mgrpc"
	"mediumkube/pkg/daemon/tasks"
	"mediumkube/pkg/services"
	"net"

	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

// DMux used for synchronizing go routines
type DMux struct {
	sync.Mutex
}

var (
	on    bool = true
	dmux  DMux = DMux{}
	close      = make(chan bool)

	dnsMasqProc *os.Process = nil
)

func stopDaemon() {
	close <- true
}

func main() {

	tmpFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	profiling := tmpFlagSet.Bool("p", false, "Enable Profiling")
	profilingPort := tmpFlagSet.Int("pport", 7777, "Port of profiling service")
	tmpFlagSet.Parse(os.Args[1:])

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGTERM)

	eventBus := event.GetEventBus()
	wg := sync.WaitGroup{}

	sigHandler := func() {
		wg.Add(1)
		defer wg.Done()

		select {
		case sig := <-c:
			klog.Info("Sig recvd: ", sig)
			stopDaemon()
		}
	}

	processBridge := func(config *common.OverallConfig) {
		bridge := configurations.Config().Bridge
		tasks.ProcessExistence(bridge)
		tasks.ProcessAddr(bridge)
		tasks.ProcessIptables(config)
	}

	dnsMasq := func(config *common.OverallConfig) {
		if dnsMasqProc == nil {
			dnsMasqProc = tasks.StartDnsmasq(config.Bridge, *config)
		}
	}

	profiler := func() {

		klog.Infof("Profiling service starting on localhost:%v/debug/pprof", *profilingPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *profilingPort), nil))
	}

	processMesh := func(config *common.OverallConfig) {
		if config.Overlay.Enabled {
			mesh.StartMesh()
		}
	}

	config := configurations.Config()
	go sigHandler()
	if *profiling {
		go profiler()
	}

	svr := grpc.NewServer()
	klog.Info("Registering Mediumkube server")
	mgrpc.RegisterDomainSerciceServer(svr, mgrpc.NewServer(config))
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Overlay.GRPCPort))
	if err != nil {
		klog.Error("Failed to start grpc server", err)
		stopDaemon()
	}

	go svr.Serve(lis)
	if err != nil {
		klog.Error("Failed to start grpc server", err)
		stopDaemon()
	}

	for {
		c := false
		select {
		case <-time.After(3 * time.Second):
			processBridge(config)
			dnsMasq(config)
			processMesh(config)
		case <-eventBus.DomainUpdate:
			mesh.SyncDomain(config)
		case c = <-close:
			mesh.StopMesh()
			dnsMasqProc.Kill()
			tasks.CleanUpIptables()
			services.GetNodeManager(config.Backend).Disconnect()
			break
		}
		if c {
			break
		}
	}

	klog.Info("Daemon exited")
}
