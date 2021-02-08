package main

import (
	"flag"
	"mediumkube/configurations"
	"mediumkube/daemon/tasks"
	"os"
	"os/signal"
	"sync"
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
	configDir := tmpFlagSet.String("config", "./config.yaml", "Configuration file")
	tmpFlagSet.Parse(os.Args)
	configurations.InitConfig(*configDir)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()

		select {
		case sig := <-c:
			klog.Info("Sig recvd: ", sig)
			stopDaemon()
			tasks.CleanUpIptables()
		}
	}()

	go func() {
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
	}()

	time.Sleep(1 * time.Second)

	go func() {
		wg.Add(1)
		defer wg.Done()
		config := configurations.Config()
		proc := tasks.StartDnsmasq(config.Bridge, *config)

		for on {
			time.Sleep(1 * time.Second)
		}
		proc.Kill()
	}()

	wg.Wait()
	klog.Info("Daemon exited")
}
