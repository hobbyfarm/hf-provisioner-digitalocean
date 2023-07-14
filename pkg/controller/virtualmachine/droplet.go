package virtualmachine

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/config"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/errors"
	labels2 "github.com/ebauman/hf-provisioner-digitalocean/pkg/labels"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/parse"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func DropletHandler(req router.Request, resp router.Response) error {
	obj := req.Object.(*v1.VirtualMachine)

	var drop *v1alpha1.Droplet
	drop, err := GetDroplet(req)

	if errors.IsNotFound(err) {
		name := fmt.Sprintf("%s-droplet", obj.Name)
		key, err := GetKey(req)
		if err != nil {
			return err
		}

		droplet := v1alpha1.Droplet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: obj.Namespace,
				Labels: map[string]string{
					labels2.VirtualMachineLabel: obj.Name,
				},
			},
			Spec: v1alpha1.DropletSpec{
				Machine: obj.Name,
			},
		}

		dcr := buildDropletCreateRequest(name, obj, req)

		dcr.Image = godo.DropletCreateImage{
			Slug: config.ResolveConfigItem(obj, req, "image"),
		}

		var doKey = godo.Key{}
		if err := json.Unmarshal(key.Status.Key.Raw, &doKey); err != nil {
			return err
		}
		dcr.SSHKeys = []godo.DropletCreateSSHKey{
			{
				ID:          doKey.ID,
				Fingerprint: doKey.Fingerprint,
			},
		}

		dcrJson, err := json.Marshal(dcr)
		if err != nil {
			return fmt.Errorf("error marshalling json: %s", err.Error())
		}
		droplet.Spec.Droplet.Raw = dcrJson
		drop = &droplet
	}

	resp.Objects(drop)

	return nil
}

func GetDroplet(req router.Request) (*v1alpha1.Droplet, error) {
	dropletList := &v1alpha1.DropletList{}
	err := req.List(dropletList, &kclient.ListOptions{
		Namespace:     req.Object.GetNamespace(),
		LabelSelector: VMLabelSelector(req.Object.GetName()),
	})

	if err != nil {
		return nil, err
	}

	if len(dropletList.Items) > 0 {
		return &dropletList.Items[0], nil
	}

	return nil, errors.NewNotFoundError("could not find any droplets for virtualmachine %s", req.Object.GetName())
}

func buildDropletCreateRequest(name string, vm *v1.VirtualMachine, req router.Request) *godo.DropletCreateRequest {
	return &godo.DropletCreateRequest{
		Name:              name,
		Region:            config.ResolveConfigItem(vm, req, "region"),
		Size:              config.ResolveConfigItem(vm, req, "size"),
		Backups:           parse.ParseBoolOrFalse(config.ResolveConfigItem(vm, req, "backups")),
		IPv6:              parse.ParseBoolOrFalse(config.ResolveConfigItem(vm, req, "ipv6")),
		PrivateNetworking: parse.ParseBoolOrFalse(config.ResolveConfigItem(vm, req, "private_networking")),
		Monitoring:        parse.ParseBoolOrFalse(config.ResolveConfigItem(vm, req, "monitoring")),
		UserData:          config.ResolveConfigItem(vm, req, "user_data"),
	}
}
