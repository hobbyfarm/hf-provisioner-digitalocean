package main

import (
	"context"
	"flag"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/cleanup"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/controller"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/log"
	"github.com/sirupsen/logrus"
)

var (
	LogLevel      = flag.Int("loglevel", int(logrus.InfoLevel), "log level")
	CleanupPeriod = flag.Int("cleanup-period", 300, "period, in seconds, between cleanup executions")
)

func init() {
	flag.Parse()
}

func main() {
	log.SetLogLevel(logrus.Level(*LogLevel))
	log.Debugf("building controller")
	ctr, err := controller.New()
	if err != nil {
		log.Fatalf("unable to build controller: %s", err.Error())
	}

	ctx := context.Background()

	log.Infof("starting controller")
	if err := ctr.Start(ctx); err != nil {
		log.Fatalf("unable to start controller: %s", err.Error())
	}

	stopCh := make(chan struct{})
	errCh := make(chan error)
	defer close(stopCh)
	log.Infof("starting cleanup goroutines")
	go cleanup.RunCleanup(ctr.Router.Backend(), *CleanupPeriod, stopCh, errCh)

	for {
		select {
		case err := <-errCh:
			log.Errorf(err.Error())
		case <-ctx.Done():
			return
		}
	}
}
