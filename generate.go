//go:generate go run github.com/acorn-io/baaah/cmd/deepcopy ./pkg/apis/provisioning.hobbyfarm.io/v1alpha1/
//go:generate go run k8s.io/kube-openapi/cmd/openapi-gen -i github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1 -p ./pkg/openapi/generated -h boilerplate/header.txt
package main

import (
	_ "github.com/acorn-io/baaah/pkg/deepcopy"
	_ "github.com/golang/mock/gomock"
	_ "k8s.io/kube-openapi/cmd/openapi-gen/args"
)
