package virtualmachine

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/errors"
	labels2 "github.com/ebauman/hf-provisioner-digitalocean/pkg/labels"
	"k8s.io/apimachinery/pkg/labels"
	"time"
)

func VMLabelSelector(vmName string) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		labels2.VirtualMachineLabel: vmName,
	})
}

func ProvisionerFinalizer(req router.Request, resp router.Response) error {
	// before deleting the VM, make sure the droplet and key are gone
	droplet, err := GetDroplet(req)
	// if the droplet is not found, move on. anything else, report!
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("error fetching droplet: %s", err.Error())
	}

	if droplet != nil {
		// droplet exists, delete it
		err = req.Client.Delete(req.Ctx, droplet)
		if err != nil {
			return fmt.Errorf("error deleting droplet: %s", err.Error())
		}
		resp.RetryAfter(5 * time.Second)
	}

	key, err := GetKey(req)
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("error fetching key: %s", err.Error())
	}

	if key != nil {
		err = req.Client.Delete(req.Ctx, key)
		if err != nil {
			return fmt.Errorf("error deleting key: %s", err.Error())
		}
		resp.RetryAfter(5 * time.Second)
	}

	return nil
}
