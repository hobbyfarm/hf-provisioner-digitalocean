package main

import (
	"context"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/controller"
	"log"
)

func main() {
	ctr, err := controller.New()
	if err != nil {
		log.Fatalf("unable to build controller: %s", err.Error())
	}

	ctx := context.Background()

	if err := ctr.Start(ctx); err != nil {
		log.Fatalf("unable to start controller: %s", err.Error())
	}

	<-ctx.Done()
}
