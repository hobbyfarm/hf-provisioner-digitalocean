package ownerreference

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetOwnerReference(obj client.Object) v1.OwnerReference {
	return v1.OwnerReference{
		Name:       obj.GetName(),
		Kind:       obj.GetObjectKind().GroupVersionKind().Kind,
		APIVersion: obj.GetObjectKind().GroupVersionKind().Version,
		UID:        obj.GetUID(),
	}
}
