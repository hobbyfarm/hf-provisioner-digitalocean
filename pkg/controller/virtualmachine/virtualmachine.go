package virtualmachine

import (
	labels2 "github.com/ebauman/hf-provisioner-digitalocean/pkg/labels"
	"k8s.io/apimachinery/pkg/labels"
)

func VMLabelSelector(vmName string) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		labels2.VirtualMachineLabel: vmName,
	})
}
