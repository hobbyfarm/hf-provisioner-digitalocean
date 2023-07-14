package v1alpha1

func (in *DropletCreateRequest) DeepCopyInto(out *DropletCreateRequest) {
	*out = *in
	out.DropletCreateRequest = in.DropletCreateRequest
	out.Image = in.Image
	if in.SSHKeys != nil {
		in, out := &in.SSHKeys, &out.SSHKeys
		*out = make([]DropletCreateSSHKey, len(*in))
		copy(*out, *in)
	}
}
