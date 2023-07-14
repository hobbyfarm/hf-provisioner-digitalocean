package config

import (
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/namespace"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/sirupsen/logrus"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func ResolveConfigItemName(vmName string, req router.Request, item string) string {
	machine := &v1.VirtualMachine{}
	err := req.Client.Get(req.Ctx, kclient.ObjectKey{
		Namespace: namespace.Resolve(),
		Name:      vmName,
	}, machine)
	if err != nil {
		logrus.Errorf("error while looking up machine with name %s: %s", vmName, err.Error())
		return "" // nothing more we can do
	}

	return ResolveConfigItem(machine, req, item)
}

func ResolveConfigItem(obj *v1.VirtualMachine, req router.Request, item string) string {
	// go from most to least specific
	env := &v1.Environment{}
	err := req.Client.Get(req.Ctx, kclient.ObjectKey{
		Namespace: obj.Namespace,
		Name:      obj.Status.EnvironmentId,
	}, env)

	if err != nil {
		logrus.Warnf("error while looking up environment for config key %s: %s", item, err.Error())
	}
	if err == nil {
		// first, check specifics for the template
		if val, ok := env.Spec.TemplateMapping[obj.Spec.VirtualMachineTemplateId][item]; ok {
			return val
		}

		// if its not there, check the environment specs
		if val, ok := env.Spec.EnvironmentSpecifics[item]; ok {
			return val
		}
	}

	return ""
}
