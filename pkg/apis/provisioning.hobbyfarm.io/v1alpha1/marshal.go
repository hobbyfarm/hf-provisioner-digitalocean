package v1alpha1

import (
	"k8s.io/apimachinery/pkg/util/json"
)

func (d *DropletCreateImage) UnmarshalJSON(data []byte) error {
	var num int
	var err = json.Unmarshal(data, &num)
	if err != nil {
		var str string
		err = json.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		d.Slug = str
		return nil
	}
	d.ID = num

	return nil
}

func (d *DropletCreateImage) MarshalJSON() ([]byte, error) {
	return d.ToGodo().MarshalJSON()
}

func (d *DropletCreateSSHKey) MarshalJSON() ([]byte, error) {
	return d.ToGodo().MarshalJSON()
}

func (d *DropletCreateSSHKey) UnmarshalJSON(data []byte) error {
	var num int
	var err = json.Unmarshal(data, &num)
	if err != nil {
		var str string
		err = json.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		d.Fingerprint = str
		return nil
	}

	d.ID = num

	return nil
}
