package virtualmachine

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/errors"
	labels2 "github.com/ebauman/hf-provisioner-digitalocean/pkg/labels"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// RequireKey not only ensures that the key has been created, but that
// it is also created in DigitalOcean
func RequireKey(next router.Handler) router.Handler {
	return router.HandlerFunc(func(req router.Request, resp router.Response) error {
		keyList := &v1alpha1.KeyList{}
		err := req.List(keyList, &kclient.ListOptions{
			Namespace: req.Object.GetNamespace(),
			LabelSelector: labels.SelectorFromSet(map[string]string{
				labels2.VirtualMachineLabel: req.Object.GetName(),
			}),
		})

		if err != nil {
			return err
		}

		if len(keyList.Items) == 0 {
			return nil
		}

		if v1alpha1.ConditionKeyCreated.GetStatus(keyList.Items[0]) != "true" {
			return nil // only passthrough keys that are created in DO
		}

		return next.Handle(req, resp)
	})
}

func KeyHandler(req router.Request, resp router.Response) error {
	secret, err := GetSecret(req)
	if err != nil {
		return err
	}

	// try to get existing key
	var key *v1alpha1.Key
	key, err = GetKey(req)
	if errors.IsNotFound(err) {
		// need to create
		key = &v1alpha1.Key{}
	} else if err != nil {
		return err // something else bad happened
	}

	name := fmt.Sprintf("%s-droplet-key", req.Object.GetName())
	key.Name = name
	key.Namespace = req.Object.GetNamespace()

	if len(key.Labels) == 0 {
		key.Labels = map[string]string{}
	}
	key.Labels[labels2.VirtualMachineLabel] = req.Object.GetName()
	key.Spec.Machine = req.Object.GetName()
	key.Spec.Secret = secret.Name

	var kcr = godo.KeyCreateRequest{
		Name:      name,
		PublicKey: string(secret.Data["public_key"]),
	}
	out, err := json.Marshal(kcr)
	if err != nil {
		return err
	}
	key.Spec.Key.Raw = out

	resp.Objects(key)

	return nil
}

func GetKey(req router.Request) (*v1alpha1.Key, error) {
	keyList := &v1alpha1.KeyList{}
	lo := &kclient.ListOptions{
		Namespace:     req.Object.GetNamespace(),
		LabelSelector: VMLabelSelector(req.Object.GetName()),
	}
	err := req.Client.List(req.Ctx, keyList, lo)
	if err != nil {
		return nil, err
	}

	if len(keyList.Items) > 0 {
		return &keyList.Items[0], nil
	}

	return nil, errors.NewNotFoundError("could not find any keys for virtualmachine %s", req.Object.GetName())
}
