package internal

import (
	_ "k8s.io/api/admission/v1"
	_ "k8s.io/api/apps/v1"
	_ "k8s.io/api/batch/v1"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/api/networking/v1"
	_ "k8s.io/api/policy/v1"
	_ "k8s.io/api/rbac/v1"
	_ "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	_ "sigs.k8s.io/gateway-api/apis/v1alpha2"
)
