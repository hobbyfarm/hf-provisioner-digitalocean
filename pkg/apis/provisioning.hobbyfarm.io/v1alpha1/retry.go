package v1alpha1

import "github.com/ebauman/hf-provisioner-digitalocean/pkg/retries"

func (k *Key) GetRetries() []retries.GenericRetry {
	return k.Status.Retries
}

func (k *Key) SetRetries(retries []retries.GenericRetry) {
	k.Status.Retries = retries
}

func (d *Droplet) GetRetries() []retries.GenericRetry {
	return d.Status.Retries
}

func (d *Droplet) SetRetries(retries []retries.GenericRetry) {
	d.Status.Retries = retries
}
