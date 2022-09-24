.PHONY: deps
deps:
	go mod tidy
	go mod vendor
	bazel run //:gazelle

kube.pb.go: kube.proto
	bazel build //:kubeproto_go_proto
	@cp bazel-bin/kubeproto_go_proto_/go.f110.dev/kubeproto/kube.pb.go ./
	@chmod 644 $@

.PHONY: gen-proto
gen-proto: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto \
	k8s.io/apimachinery/pkg/api/resource/generated.proto \
	k8s.io/apimachinery/pkg/util/intstr/generated.proto \
	k8s.io/apimachinery/pkg/runtime/generated.proto \
	k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto \
	sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto \
	k8s.io/api/apps/v1/generated.proto \
	k8s.io/api/core/v1/generated.proto

.PHONY: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto
k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(shell pwd)/$@ --proto-package k8s.io.apimachinery.pkg.apis.meta.v1 --go-package $(@D) --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/apimachinery/pkg/api/resource/generated.proto
k8s.io/apimachinery/pkg/api/resource/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.apimachinery.pkg.api.resource --go-package $(@D) --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/apimachinery/pkg/util/intstr/generated.proto
k8s.io/apimachinery/pkg/util/intstr/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.apimachinery.pkg.util.intstr --go-package $(@D) --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/apimachinery/pkg/runtime/generated.proto
k8s.io/apimachinery/pkg/runtime/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(shell pwd)/$@ --proto-package k8s.io.apimachinery.pkg.runtime --go-package $(@D) $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/core/v1/generated.proto
k8s.io/api/core/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.core.v1 --go-package $(@D) --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/apps/v1/generated.proto
k8s.io/api/apps/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.apps.v1 --go-package $(@D) --api-domain apps --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto
k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.apiextensions_apiserver.pkg.apis.apiextensions.v1 --go-package $(@D) --api-domain apiextensions --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto
sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package sigs.k8s.io.gateway_api.apis.v1alpha2 --go-package $(@D) --api-domain gateway --api-sub-group networking.k8s.io --api-version v1alpha2 --all $(CURDIR)/vendor/$(@D)
