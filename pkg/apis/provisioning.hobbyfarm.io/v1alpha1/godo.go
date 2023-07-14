package v1alpha1

import "github.com/digitalocean/godo"

func (d *DropletCreateRequest) ToGodo() *godo.DropletCreateRequest {
	return &godo.DropletCreateRequest{
		Name:   d.Name,
		Region: d.Region,
		Size:   d.Size,
		Image:  *d.Image.ToGodo(),
		SSHKeys: func() []godo.DropletCreateSSHKey {
			var out = make([]godo.DropletCreateSSHKey, len(d.SSHKeys))
			for i, k := range d.SSHKeys {
				out[i] = *k.ToGodo()
			}
			return out
		}(),
		Backups:           d.Backups,
		IPv6:              d.IPv6,
		PrivateNetworking: d.PrivateNetworking,
		Monitoring:        d.Monitoring,
		UserData:          d.UserData,
		Volumes:           d.Volumes,
		Tags:              d.Tags,
		VPCUUID:           d.VPCUUID,
		WithDropletAgent:  d.WithDropletAgent,
	}
}

func (d *DropletCreateImage) ToGodo() *godo.DropletCreateImage {
	return &godo.DropletCreateImage{
		ID:   d.ID,
		Slug: d.Slug,
	}
}

func (d *DropletCreateSSHKey) ToGodo() *godo.DropletCreateSSHKey {
	return &godo.DropletCreateSSHKey{
		ID:          d.ID,
		Fingerprint: d.Fingerprint,
	}
}
