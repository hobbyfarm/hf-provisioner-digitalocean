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
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func DropletHandler(req router.Request, resp router.Response) error {
	obj := req.Object.(*v1.VirtualMachine)
	name := fmt.Sprintf("%s-droplet", obj.Name)

	key, err := GetKey(req)
	if err != nil {
		return err
	}

	var droplet *v1alpha1.Droplet
	droplet, err = GetDroplet(req)
	if errors.IsNotFound(err) {
		droplet = &v1alpha1.Droplet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: obj.Namespace,
				Name:      name,
				Labels: map[string]string{
					labels2.VirtualMachineLabel: obj.Name,
				},
			},
			Spec: v1alpha1.DropletSpec{
				Machine: obj.Name,
				DropletCreateRequest: godo.DropletCreateRequest{
					Name:   name,
					Region: config.ResolveConfigItem(obj, req, "region"),
					Size:   config.ResolveConfigItem(obj, req, "size"),
					Image: godo.DropletCreateImage{
						Slug: config.ResolveConfigItem(obj, req, "image"),
					},
					SSHKeys: []godo.DropletCreateSSHKey{
						{
							ID:          key.Status.ID,
							Fingerprint: key.Status.Fingerprint,
						},
					},
					Backups:           parse.ParseBoolOrFalse(config.ResolveConfigItem(obj, req, "backups")),
					IPv6:              parse.ParseBoolOrFalse(config.ResolveConfigItem(obj, req, "ipv6")),
					PrivateNetworking: parse.ParseBoolOrFalse(config.ResolveConfigItem(obj, req, "private_networking")),
					Monitoring:        parse.ParseBoolOrFalse(config.ResolveConfigItem(obj, req, "monitoring")),
					UserData:          config.ResolveConfigItem(obj, req, "user_data"),
					Volumes:           nil,
					Tags:              nil,
					VPCUUID:           "",
					WithDropletAgent:  nil,
				},
			},
		}
	} else if err != nil {
		return err
	}

	resp.Objects(droplet)

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
