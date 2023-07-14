package finalizer

import "k8s.io/utils/strings/slices"

type GetFinalizers interface {
	GetFinalizers() []string
}

type SetFinalizers interface {
	SetFinalizers([]string)
}

type UpdateFinalizers interface {
	GetFinalizers
	SetFinalizers
}

func ClearFinalizer(finalizer string, obj UpdateFinalizers) {
	var finalizers = obj.GetFinalizers()

	if i := slices.Index(finalizers, finalizer); i >= 0 {
		obj.SetFinalizers(append(finalizers[i:], finalizers[i+1:]...))
	}

	obj.SetFinalizers(finalizers)
}
