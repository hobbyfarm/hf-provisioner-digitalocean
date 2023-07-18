package controller

import (
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/controller/droplet"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/controller/key"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/controller/virtualmachine"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/labels"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/namespace"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/providerregistration"
	hfv1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	klabels "k8s.io/apimachinery/pkg/labels"
)

func routes(r *router.Router) {
	vmRouter := r.Type(&hfv1.VirtualMachine{}).Namespace(namespace.Resolve()).Selector(klabels.SelectorFromSet(map[string]string{
		labels.ProvisionerLabel: providerregistration.ProviderName(),
	}))

	vmRouter.FinalizeFunc(labels.Finalizer, virtualmachine.ProvisionerFinalizer)
	vmRouter.HandlerFunc(virtualmachine.SecretHandler)
	vmRouter.Middleware(virtualmachine.RequireSecret).HandlerFunc(virtualmachine.KeyHandler)
	vmRouter.Middleware(virtualmachine.RequireKey).
		HandlerFunc(virtualmachine.DropletHandler)

	keyRouter := r.Type(&v1alpha1.Key{}).Namespace(namespace.Resolve())

	keyRouter.HandlerFunc(key.EnsureStatus)
	keyRouter.Middleware(key.NotYetCreated).HandlerFunc(key.CreateKey)
	keyRouter.Middleware(key.Created).HandlerFunc(key.WriteVM)
	keyRouter.FinalizeFunc(key.EnsureDeletedFinalizer, key.EnsureDeleted)

	dropletRouter := r.Type(&v1alpha1.Droplet{}).Namespace(namespace.Resolve())

	dropletRouter.HandlerFunc(droplet.EnsureStatus)
	dropletRouter.Middleware(droplet.DropletNotCreated).HandlerFunc(droplet.CreateDroplet)
	dropletRouter.HandlerFunc(droplet.PeriodicUpdate)
	dropletRouter.HandlerFunc(droplet.UpdateVMStatus)
	dropletRouter.FinalizeFunc(droplet.EnsureDeletedFinalizer, droplet.EnsureDeleted)
}
