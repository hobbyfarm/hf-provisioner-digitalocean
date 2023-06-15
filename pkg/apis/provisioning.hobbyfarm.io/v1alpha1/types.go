package v1alpha1

import (
	"encoding/json"
	"github.com/digitalocean/godo"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/genericcondition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ConditionDropletExists = condition.Cond("DropletExists")
	ConditionDropletReady  = condition.Cond("DropletReady")

	ConditionKeyExists       = condition.Cond("KeyExists")
	ConditionKeySecretExists = condition.Cond("KeySecretExists")
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Droplet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,inline"`

	Spec   DropletSpec   `json:"spec"`
	Status DropletStatus `json:"status"`
}

// +k8s:deepcopy-gen=true

type DropletSpec struct {
	// Name of the HobbyFarm machine that spawned this droplet
	Machine string `json:"machine"`

	godo.DropletCreateRequest
}

// MarshalJSON This method exists to prevent json from calling godo.DropletCreateImage's MarshalJSON()
// why? Because that method returns a scalar int or string (from the original struct{int, string})
// which doesn't play nice with a kubernetes server that expects to store a json object, not a json scalar.
// so we do some json merge trickery here so that the output json from this method
// uses the "overridden" 'image' tag from our embedded struct instead of
// the 'image' tag from godo.DropletCreateImage embedded in godo.DropletCreateRequest
// see https://choly.ca/post/go-json-marshalling/
//
// The above also applies to godo.DropletCreateSSHKey
func (d *DropletSpec) MarshalJSON() ([]byte, error) {
	type alias DropletSpec

	type SSHKey struct {
		ID          int    `json:"id"`
		Fingerprint string `json:"fingerprint"`
	}

	var keys = make([]SSHKey, len(d.SSHKeys))
	for i, k := range d.SSHKeys {
		keys[i] = SSHKey{
			ID:          k.ID,
			Fingerprint: k.Fingerprint,
		}
	}

	var out = &struct {
		Image   any `json:"image"`
		SSHKeys any `json:"ssh_keys"`
		*alias
	}{
		struct {
			ID   int    `json:"id"`
			Slug string `json:"slug"`
		}{ID: d.Image.ID, Slug: d.Image.Slug},
		keys,
		(*alias)(d),
	}

	return json.Marshal(out)
}

type DropletCreateImage struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

// +k8s:deepcopy-gen=true

type DropletStatus struct {
	godo.Droplet
	Conditions []genericcondition.GenericCondition
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DropletList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Droplet
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Key struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,inline"`

	Spec   KeySpec   `json:"spec"`
	Status KeyStatus `json:"status"`
}

// +k8s:deepcopy-gen=true

type KeySpec struct {
	// HF machine with which this key is associated
	Machine string `json:"machine"`

	godo.KeyCreateRequest
}

// +k8s:deepcopy-gen=true

type KeyStatus struct {
	// Name of the secret in which the key details are stored
	Secret string `json:"secret"`

	godo.Key

	Conditions []genericcondition.GenericCondition
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Key
}
