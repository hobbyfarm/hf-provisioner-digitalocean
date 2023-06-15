package v1alpha1

import "github.com/digitalocean/godo"

func (in *DropletSpec) DeepCopyInto(out *DropletSpec) {
	*out = *in
	out = in.DeepCopy()
}

func (in *DropletSpec) DeepCopy() *DropletSpec {
	out := &DropletSpec{}

	out.Machine = in.Machine

	out.DropletCreateRequest = godo.DropletCreateRequest{
		Name:   in.Name,
		Region: in.Region,
		Size:   in.Size,
		Image:  *deepCopyDropletCreateImage(&in.Image),
		SSHKeys: func() []godo.DropletCreateSSHKey {
			var outKeys = make([]godo.DropletCreateSSHKey, len(in.SSHKeys))
			for i, k := range in.SSHKeys {
				outKeys[i] = *deepCopyDropletCreateSSHKey(&k)
			}
			return outKeys
		}(),
		Backups:           in.Backups,
		IPv6:              in.IPv6,
		PrivateNetworking: in.PrivateNetworking,
		Monitoring:        in.Monitoring,
		UserData:          in.UserData,
		Volumes: func() []godo.DropletCreateVolume {
			var outVols = make([]godo.DropletCreateVolume, len(in.Volumes))
			for i, k := range in.Volumes {
				outVols[i] = *deepCopyDropletCreateVolume(&k)
			}
			return outVols
		}(),
		Tags:             in.Tags,
		VPCUUID:          in.VPCUUID,
		WithDropletAgent: in.WithDropletAgent,
	}

	return out
}

func (in *DropletStatus) DeepCopyInto(out *DropletStatus) {
	*out = *in
	out = in.DeepCopy()
}

func (in *DropletStatus) DeepCopy() *DropletStatus {
	out := &DropletStatus{}

	*out = *in

	out.Droplet = godo.Droplet{
		ID:               in.ID,
		Name:             in.Name,
		Memory:           in.Memory,
		Vcpus:            in.Vcpus,
		Disk:             in.Disk,
		Region:           deepCopyRegion(in.Region),
		Image:            deepCopyImage(in.Image),
		Size:             deepCopySize(in.Size),
		SizeSlug:         in.SizeSlug,
		BackupIDs:        in.BackupIDs,
		NextBackupWindow: in.NextBackupWindow,
		SnapshotIDs:      in.SnapshotIDs,
		Features:         in.Features,
		Locked:           in.Locked,
		Status:           in.Status,
		Networks:         deepCopyNetworks(in.Networks),
		Created:          in.Created,
		Kernel:           deepCopyKernel(in.Kernel),
		Tags:             in.Tags,
		VolumeIDs:        in.VolumeIDs,
		VPCUUID:          in.VPCUUID,
	}

	return out
}

func deepCopyRegion(in *godo.Region) *godo.Region {
	if in == nil {
		return nil
	}

	return &godo.Region{
		Slug:      in.Slug,
		Name:      in.Name,
		Sizes:     in.Sizes,
		Available: in.Available,
		Features:  in.Features,
	}
}

func deepCopyImage(in *godo.Image) *godo.Image {
	if in == nil {
		return nil
	}

	return &godo.Image{
		ID:            in.ID,
		Name:          in.Name,
		Type:          in.Type,
		Distribution:  in.Distribution,
		Slug:          in.Slug,
		Public:        in.Public,
		Regions:       in.Regions,
		MinDiskSize:   in.MinDiskSize,
		SizeGigaBytes: in.SizeGigaBytes,
		Created:       in.Created,
		Description:   in.Description,
		Tags:          in.Tags,
		Status:        in.Status,
		ErrorMessage:  in.ErrorMessage,
	}
}

func deepCopySize(in *godo.Size) *godo.Size {
	if in == nil {
		return nil
	}

	return &godo.Size{
		Slug:         in.Slug,
		Memory:       in.Memory,
		Vcpus:        in.Vcpus,
		Disk:         in.Disk,
		PriceMonthly: in.PriceMonthly,
		PriceHourly:  in.PriceHourly,
		Regions:      in.Regions,
		Available:    in.Available,
		Transfer:     in.Transfer,
		Description:  in.Description,
	}
}

func deepCopyNetworks(in *godo.Networks) *godo.Networks {
	if in == nil {
		return nil
	}

	return &godo.Networks{
		V4: in.V4,
		V6: in.V6,
	}
}

func deepCopyKernel(in *godo.Kernel) *godo.Kernel {
	if in == nil {
		return nil
	}

	return &godo.Kernel{
		ID:      in.ID,
		Name:    in.Name,
		Version: in.Version,
	}
}

func deepCopyDropletCreateImage(in *godo.DropletCreateImage) *godo.DropletCreateImage {
	if in == nil {
		return nil
	}

	return &godo.DropletCreateImage{
		ID:   in.ID,
		Slug: in.Slug,
	}
}

func deepCopyDropletCreateSSHKey(in *godo.DropletCreateSSHKey) *godo.DropletCreateSSHKey {
	if in == nil {
		return nil
	}

	return &godo.DropletCreateSSHKey{
		ID:          in.ID,
		Fingerprint: in.Fingerprint,
	}
}

func deepCopyDropletCreateVolume(in *godo.DropletCreateVolume) *godo.DropletCreateVolume {
	if in == nil {
		return nil
	}

	return &godo.DropletCreateVolume{
		ID:   in.ID,
		Name: in.Name,
	}
}

func (in *KeySpec) DeepCopy() *KeySpec {
	out := &KeySpec{}

	out.Machine = in.Machine
	out.KeyCreateRequest = *deepCopyKeyCreateRequest(&in.KeyCreateRequest)

	return out
}

func (in *KeySpec) DeepCopyInto(out *KeySpec) {
	*out = *in

	out = in.DeepCopy()
}

func (in *KeyStatus) DeepCopy() *KeyStatus {
	out := &KeyStatus{}

	out.Conditions = in.Conditions
	out.Secret = in.Secret
	out.Key = *deepCopyKey(&in.Key)

	return out
}

func (in *KeyStatus) DeepCopyInto(out *KeyStatus) {
	*out = *in
	out.Secret = in.Secret
	out.Conditions = in.Conditions
	out.Key = *deepCopyKey(&in.Key)
}

func deepCopyKey(in *godo.Key) *godo.Key {
	if in == nil {
		return nil
	}

	return &godo.Key{
		ID:          in.ID,
		Name:        in.Name,
		Fingerprint: in.Fingerprint,
		PublicKey:   in.PublicKey,
	}
}

func deepCopyKeyCreateRequest(in *godo.KeyCreateRequest) *godo.KeyCreateRequest {
	if in == nil {
		return nil
	}

	return &godo.KeyCreateRequest{
		Name:      in.Name,
		PublicKey: in.PublicKey,
	}
}
