package apis

import (
	_ "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	_ "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	_ "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	_ "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
)
