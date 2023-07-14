package v1alpha1

import (
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/retries"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/genericcondition"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ConditionDropletExists  = condition.Cond("DigitaloceanDropletExists")
	ConditionDropletReady   = condition.Cond("DigitaloceanDropletReady")
	ConditionDropletUpdated = condition.Cond("DigitaloceanDropletUpdated")

	ConditionKeyExists  = condition.Cond("KeyExists")
	ConditionKeyCreated = condition.Cond("DigitaloceanKeyCreated")

	RetryDeleteKey     = retries.NewRetry("DeleteKey", 0)
	RetryDeleteDroplet = retries.NewRetry("DeleteDroplet", 0)
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
	Machine string  `json:"machine"`
	Droplet v1.JSON `json:"droplet"`
}

// +k8s:deepcopy-gen=true

type DropletCreateRequest struct {
	godo.DropletCreateRequest
	Image   DropletCreateImage    `json:"image"`
	SSHKeys []DropletCreateSSHKey `json:"ssh_keys"`
}

// +k8s:deepcopy-gen=true

type DropletCreateImage struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

// +k8s:deepcopy-gen=true

type DropletCreateSSHKey struct {
	ID          int    `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

// +k8s:deepcopy-gen=true

type DropletStatus struct {
	Droplet    v1.JSON                             `json:"droplet"`
	Conditions []genericcondition.GenericCondition `json:"conditions"`
	Retries    []retries.GenericRetry              `json:"retries"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DropletList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Droplet `json:"items,omitempty"`
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
	Secret  string `json:"secret"`

	Key v1.JSON `json:"key"`
}

// +k8s:deepcopy-gen=true

type KeyStatus struct {
	Key        v1.JSON                             `json:"key"`
	Conditions []genericcondition.GenericCondition `json:"conditions"`
	Retries    []retries.GenericRetry              `json:"retries"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Key `json:"items,omitempty"`
}
