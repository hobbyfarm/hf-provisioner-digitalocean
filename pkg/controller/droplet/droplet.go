package droplet

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/digitalocean"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const EnsureDeletedFinalizer = "provisioner.hobbyfarm.io/digitalocean-droplet-deleted"

func PeriodicUpdate(req router.Request, resp router.Response) error {
	k8sDroplet := req.Object.(*v1alpha1.Droplet)

	// get the old godo droplet from k8s
	var godoDroplet = godo.Droplet{}
	err := json.Unmarshal(k8sDroplet.Status.Droplet.Raw, &godoDroplet)
	if err != nil {
		logrus.Errorf("Unable to unmarshal DigitalOcean droplet from k8s droplet object: %s", err.Error())
		return nil
	}

	// get an updated droplet object from DO
	dClient, err := digitalocean.GetGodoClient(k8sDroplet.Spec.Machine, req)
	newGodoDroplet, _, err := dClient.Droplets.Get(req.Ctx, godoDroplet.ID)
	if err != nil {
		v1alpha1.ConditionDropletUpdated.False(k8sDroplet)
		v1alpha1.ConditionDropletUpdated.SetError(k8sDroplet, "digitalocean error", err)
		return req.Client.Status().Update(req.Ctx, k8sDroplet)
	}

	jsonDroplet, err := json.Marshal(newGodoDroplet)
	if err != nil {
		logrus.Errorf("Unable to marshal DigitalOcean droplet to []byte for storage in k8s: %s", err.Error())
		return nil
	}

	k8sDroplet.Status.Droplet.Raw = jsonDroplet

	switch newGodoDroplet.Status {
	case "new":
		resp.RetryAfter(10 * time.Second)
	case "active":
		resp.RetryAfter(30 * time.Second)
	}

	return req.Client.Status().Update(req.Ctx, k8sDroplet)
}

func UpdateVMStatus(req router.Request, resp router.Response) error {
	k8sDroplet := req.Object.(*v1alpha1.Droplet)

	droplet := godo.Droplet{}
	err := json.Unmarshal(k8sDroplet.Status.Droplet.Raw, &droplet)
	if err != nil {
		logrus.Errorf("Unable to unmarshal DigitalOcean droplet from k8s droplet object: %s", err.Error())
		return nil
	}

	// get the corresponding VM for the droplet
	vm := &v1.VirtualMachine{}
	err = req.Client.Get(req.Ctx, client.ObjectKey{Namespace: k8sDroplet.Namespace, Name: k8sDroplet.Spec.Machine}, vm)
	if err != nil {
		// could not get VM
		logrus.Errorf("Could not get VirtualMachine with name %s: %s", k8sDroplet.Spec.Machine, err.Error())
		return nil
	}

	switch droplet.Status {
	case "new":
		vm.Status.Status = v1.VmStatusProvisioned
	case "active":
		vm.Status.Status = v1.VmStatusRunning
	case "off":
	case "archive":
		vm.Status.Status = v1.VmStatusTerminating
	}

	for _, n := range droplet.Networks.V4 {
		switch n.Type {
		case "public":
			vm.Status.PublicIP = n.IPAddress
		case "private":
			vm.Status.PrivateIP = n.IPAddress
		}
	}

	return req.Client.Status().Update(req.Ctx, vm)
}

func EnsureDeleted(req router.Request, resp router.Response) error {
	droplet := req.Object.(*v1alpha1.Droplet)

	var godoDroplet = godo.Droplet{}
	if err := json.Unmarshal(droplet.Status.Droplet.Raw, &godoDroplet); err != nil || godoDroplet.ID == 0 {
		logrus.Errorf("Could not obtain DO droplet from object JSON. This is unrecoverable, may result in "+
			"an orphan droplet in DigitalOcean. Droplet (k8s) name was %s", droplet.Name)
		return nil
	}

	// droplet exists, delete it
	dClient, err := digitalocean.GetGodoClient(droplet.Spec.Machine, req)
	if err != nil {
		logrus.Errorf("building digitalocean client: %s", err.Error())
		return req.Client.Status().Update(req.Ctx, droplet)
	}

	_, err = dClient.Droplets.Delete(req.Ctx, godoDroplet.ID)
	if err != nil {
		logrus.Errorf("deleting droplet in digitalocean: %s", err.Error())
		return req.Client.Status().Update(req.Ctx, droplet)
	}

	return nil
}

func EnsureStatus(req router.Request, resp router.Response) error {
	droplet := req.Object.(*v1alpha1.Droplet)

	if len(droplet.Status.Conditions) == 0 {
		v1alpha1.ConditionDropletExists.SetStatus(droplet, "unknown")
		v1alpha1.ConditionDropletReady.SetStatus(droplet, "unknown")

		return req.Client.Status().Update(req.Ctx, droplet)
	}

	return nil
}

func DropletNotCreated(next router.Handler) router.Handler {
	return router.HandlerFunc(func(req router.Request, resp router.Response) error {
		droplet := req.Object.(*v1alpha1.Droplet)

		if v1alpha1.ConditionDropletExists.GetStatus(droplet) == "unknown" {
			return next.Handle(req, resp)
		}

		return nil
	})
}

func CreateDroplet(req router.Request, _ router.Response) error {
	droplet := req.Object.(*v1alpha1.Droplet)

	var dcr v1alpha1.DropletCreateRequest
	if err := json.Unmarshal(droplet.Spec.Droplet.Raw, &dcr); err != nil {
		return fmt.Errorf("error unmarshalling droplet: %s", err.Error())
	}
	dClient, err := digitalocean.GetGodoClient(droplet.Spec.Machine, req)
	if err != nil {
		v1alpha1.ConditionDropletExists.SetStatus(droplet, "false")
		v1alpha1.ConditionDropletExists.SetError(droplet, "error creating digitalocean client", err)
		return req.Client.Status().Update(req.Ctx, droplet)
	}

	newDroplet, _, err := dClient.Droplets.Create(req.Ctx, dcr.ToGodo())
	if err != nil {
		v1alpha1.ConditionDropletExists.SetStatus(droplet, "false")
		v1alpha1.ConditionDropletExists.SetError(droplet, "digitalocean error", err)
		return req.Client.Status().Update(req.Ctx, droplet)
	}

	dropletJson, err := json.Marshal(newDroplet)
	if err != nil {
		return fmt.Errorf("error marshalling droplet: %s", err.Error())
	}
	droplet.Status.Droplet.Raw = dropletJson
	v1alpha1.ConditionDropletExists.SetStatus(droplet, "true")
	v1alpha1.ConditionDropletExists.Reason(droplet, "droplet created")

	return req.Client.Status().Update(req.Ctx, droplet)
}
