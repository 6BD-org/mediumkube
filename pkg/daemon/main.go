package main

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/daemon/mesh"
	"mediumkube/pkg/daemon/tasks"

	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	configDir := tmpFlagSet.String("config", "/etc/mediumkube/config.yaml", "Configuration file")
	profiling := tmpFlagSet.Bool("p", false, "Enable Profiling")
	profilingPort := tmpFlagSet.Int("pport", 7777, "Port of profiling service")
	tmpFlagSet.Parse(os.Args[1:])
	configurations.InitConfig(*configDir)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGTERM)

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
		tasks.ProcessIptables(bridge)
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

	for {
		c := false
		select {
		case <-time.After(3 * time.Second):
			processBridge(config)
			dnsMasq(config)
			processMesh(config)
		case c = <-close:
			mesh.StopMesh()
			dnsMasqProc.Kill()
			tasks.CleanUpIptables()
			break
		}
		if c {
			break
		}
	}

	klog.Info("Daemon exited")
}
