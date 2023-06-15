package controller

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/controller/virtualmachine"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/labels"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/providerregistration"
	hfv1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	klabels "k8s.io/apimachinery/pkg/labels"
)

func routes(r *router.Router) {
	vmRouter := r.Type(&hfv1.VirtualMachine{}).Selector(klabels.SelectorFromSet(map[string]string{
		labels.ProvisionerLabel: providerregistration.ProviderName(),
	}))

	vmRouter = vmRouter.Middleware(func(h router.Handler) router.Handler {
		return router.HandlerFunc(func(req router.Request, resp router.Response) error {
			fmt.Printf("calling. key %s\n", req.Key)

			return h.Handle(req, resp)
		})
	})

	vmRouter.HandlerFunc(virtualmachine.SecretHandler)
	vmRouter.Middleware(virtualmachine.RequireSecret).HandlerFunc(virtualmachine.KeyHandler)
	vmRouter.Middleware(virtualmachine.RequireKey).HandlerFunc(virtualmachine.DropletHandler)
}
