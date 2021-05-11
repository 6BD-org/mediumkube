package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	flannel "mediumkube/common/flannel"
	"mediumkube/configurations"
	"mediumkube/daemon/tasks"
	"mediumkube/utils"
	"strings"

	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	clientv2 "go.etcd.io/etcd/client/v2"
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
			tasks.CleanUpRoute()

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
				flannel := configurations.Config().Overlay.Flannel
				tasks.ProcessExistence(bridge)
				tasks.ProcessAddr(bridge)
				tasks.ProcessIptables(bridge)
				tasks.ProcessRoute(bridge, flannel)
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

	startFlannel := func() {
		wg.Add(1)
		defer wg.Done()

		proc := tasks.StartFlannel()
		for on {
			time.Sleep(1 * time.Second)
		}
		proc.Kill()
	}

	etcd := func() {
		wg.Add(1)
		defer wg.Done()
		proc := tasks.StartEtcd()
		for on {
			time.Sleep(1 * time.Second)
		}
		proc.Kill()

	}

	initNetwork := func() {
		klog.Info("Initializing configurations for flannel")
		overlayConfig := configurations.Config().Overlay
		cli, err := clientv2.New(
			clientv2.Config{
				Endpoints: []string{
					utils.EtcdEp(overlayConfig.Master, overlayConfig.EtcdPort),
				},
			},
		)
		if err != nil {
			klog.Error("Failed to init network configurations to etcd")
		}

		k := strings.Join([]string{overlayConfig.Flannel.EtcdPrefix, "config"}, "/")
		v := flannel.NewConfig(configurations.Config()).ToStr()
		kpi := clientv2.NewKeysAPI(cli)

		_, err = kpi.Set(context.TODO(), k, v, &clientv2.SetOptions{})
		klog.Info(k, v)
		if err != nil {
			klog.Error(err)
		}

	}

	profiler := func() {

		klog.Infof("Profiling service starting on localhost:%v/debug/pprof", *profilingPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *profilingPort), nil))
	}

	go sigHandler()
	go etcd()

	klog.Info("Starting ETCD")
	time.Sleep(2 * time.Second)
	initNetwork()

	go startFlannel()

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
