package main

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/pkg/configurations"
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
	on   bool = true
	dmux DMux = DMux{}
)

func stopDaemon() {
	dmux.Lock()
	on = false
	dmux.Unlock()
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
			tasks.CleanUpIptables()
		}
	}

	bridgeProcessor := func() {
		wg.Add(1)
		defer wg.Done()
		for on {
			dmux.Lock()
			if on {
				time.Sleep(1 * time.Second)
				bridge := configurations.Config().Bridge
				tasks.ProcessExistence(bridge)
				tasks.ProcessAddr(bridge)
				tasks.ProcessIptables(bridge)
			}
			dmux.Unlock()
		}
	}

	dnsMasq := func() {
		wg.Add(1)
		defer wg.Done()
		config := configurations.Config()
		proc := tasks.StartDnsmasq(config.Bridge, *config)

		for on {
			time.Sleep(1 * time.Second)
		}
		proc.Kill()
	}

	profiler := func() {

		klog.Infof("Profiling service starting on localhost:%v/debug/pprof", *profilingPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *profilingPort), nil))
	}

	go sigHandler()

	klog.Info("Starting ETCD")
	time.Sleep(2 * time.Second)

	go bridgeProcessor()
	time.Sleep(1 * time.Second)
	go dnsMasq()

	klog.Info("Starting Flannel")
	if *profiling {
		go profiler()
	}

	wg.Wait()
	klog.Info("Daemon exited")
}
