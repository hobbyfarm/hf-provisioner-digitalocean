package virtualmachine

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/config"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/errors"
	labels2 "github.com/ebauman/hf-provisioner-digitalocean/pkg/labels"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/ssh"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func RequireSecret(next router.Handler) router.Handler {
	return router.HandlerFunc(func(req router.Request, resp router.Response) error {
		vm := req.Object.(*v1.VirtualMachine)

		secretList := &corev1.SecretList{}
		err := req.List(secretList, &kclient.ListOptions{
			Namespace: vm.Namespace,
			LabelSelector: labels.SelectorFromSet(map[string]string{
				labels2.VirtualMachineLabel: vm.Name,
			}),
		})
		if err != nil {
			return err
		}

		if len(secretList.Items) == 0 {
			return nil
		}

		return next.Handle(req, resp)
	})
}

func SecretHandler(req router.Request, resp router.Response) error {
	vm := req.Object.(*v1.VirtualMachine)

	// try to get secret
	var secret *corev1.Secret
	var public, private string
	secret, err := GetSecret(req)
	if errors.IsNotFound(err) {
		secret = &corev1.Secret{}

		public, private, err = ssh.GenKeyPair()
		if err != nil {
			return err
		}

		secret.Data = map[string][]byte{}

		secret.Data["public_key"] = []byte(public)
		secret.Data["private_key"] = []byte(private)
	} else if err != nil {
		return err
	}

	secret.Name = fmt.Sprintf("%s-droplet-keys", vm.Name)

	if len(secret.Labels) == 0 {
		secret.Labels = map[string]string{}
	}

	secret.Labels[labels2.VirtualMachineLabel] = vm.Name
	secret.Namespace = vm.Namespace
	secret.Data["password"] = []byte(config.ResolveConfigItem(vm, req, "password"))

	resp.Objects(secret)

	return nil
}

func GetSecret(req router.Request) (*corev1.Secret, error) {
	// list secrets associated with the vm
	secretList := &corev1.SecretList{}
	err := req.List(secretList, &kclient.ListOptions{
		Namespace: req.Object.GetNamespace(),
		LabelSelector: labels.SelectorFromSet(map[string]string{
			labels2.VirtualMachineLabel: req.Object.GetName(),
		}),
	})
	if err != nil {
		return nil, err
	}

	if len(secretList.Items) > 0 {
		return &secretList.Items[0], nil
	}

	return nil, errors.NewNotFoundError("could not find secret for VirtualMachine %s", req.Object.GetName())
}
