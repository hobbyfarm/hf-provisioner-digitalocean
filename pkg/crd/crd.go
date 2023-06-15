package crd

import (
	"github.com/ebauman/crder"
	provisioning_hobbyfarm_io "github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
)

func Setup() []crder.CRD {
	droplet := crder.NewCRD(v1alpha1.Droplet{}, provisioning_hobbyfarm_io.Group, func(c *crder.CRD) {
		c.WithShortNames("drop")
		c.IsNamespaced(true)
		c.AddVersion(v1alpha1.Version, v1alpha1.Droplet{}, func(cv *crder.Version) {
			cv.IsStored(true).IsServed(true)
			cv.WithStatus()
		})
	})

	key := crder.NewCRD(v1alpha1.Key{}, provisioning_hobbyfarm_io.Group, func(c *crder.CRD) {
		c.IsNamespaced(true)
		c.AddVersion(v1alpha1.Version, v1alpha1.Key{}, func(cv *crder.Version) {
			cv.IsStored(true).IsServed(true)
			cv.WithStatus()
		})
	})

	return []crder.CRD{
		*droplet,
		*key,
	}
}
