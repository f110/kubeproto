.PHONY: deps
deps:
	go mod tidy
	go mod vendor
	bazel run //:gazelle

kube.pb.go: kube.proto
	bazel build //:kubeproto_go_proto
	@cp bazel-bin/kubeproto_go_proto_/go.f110.dev/kubeproto/kube.pb.go ./
	@chmod 644 $@

.PHONY: gen
gen: kube.pb.go gen-proto gen-go

.PHONY: gen-proto
gen-proto: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto \
	k8s.io/apimachinery/pkg/api/resource/generated.proto \
	k8s.io/apimachinery/pkg/util/intstr/generated.proto \
	k8s.io/apimachinery/pkg/runtime/generated.proto \
	k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto \
	sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto \
	k8s.io/api/apps/v1/generated.proto \
	k8s.io/api/core/v1/generated.proto \
	k8s.io/api/batch/v1/generated.proto \
	k8s.io/api/admission/v1/generated.proto \
	k8s.io/api/authentication/v1/generated.proto \
	k8s.io/api/policy/v1/generated.proto \
	k8s.io/api/networking/v1/generated.proto \
	k8s.io/api/rbac/v1/generated.proto \
	k8s.io/api/admissionregistration/v1/generated.proto

.PHONY: gen-go
gen-go: go/apis/metav1/metav1_kubeproto.generated.object.go \
	go/apis/corev1/corev1_kubeproto.generated.object.go \
	go/apis/appsv1/appsv1_kubeproto.generated.object.go \
	go/apis/batchv1/batchv1_kubeproto.generated.object.go \
	go/apis/authenticationv1/authenticationv1_kubeproto.generated.object.go \
	go/apis/admissionv1/admissionv1_kubeproto.generated.object.go \
	go/apis/policyv1/policyv1_kubeproto.generated.object.go \
	go/apis/networkingv1/networkingv1_kubeproto.generated.object.go \
	go/apis/rbacv1/rbacv1_kubeproto.generated.object.go \
	go/k8sclient/go_client.generated.client.go \
	go/k8stestingclient/go_testingclient.generated.testingclient.go

.PHONY: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto
k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(shell pwd)/$@ --proto-package k8s.io.apimachinery.pkg.apis.meta.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/metav1" --all $(CURDIR)/vendor/$(@D)

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
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.core.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/corev1" --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/apps/v1/generated.proto
k8s.io/api/apps/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.apps.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/appsv1" --api-domain apps --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/batch/v1/generated.proto
k8s.io/api/batch/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.batch.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/batchv1" --api-domain apps --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/authentication/v1/generated.proto
k8s.io/api/authentication/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.authentication.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/authenticationv1" --api-domain authentication.k8s.io --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/admission/v1/generated.proto
k8s.io/api/admission/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.admission.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/admissionv1" --api-domain admission.k8s.io --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/policy/v1/generated.proto
k8s.io/api/policy/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.policy.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/policyv1" --api-domain policy --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/networking/v1/generated.proto
k8s.io/api/networking/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.networking.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/networkingv1" --api-domain networking.k8s.io --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/rbac/v1/generated.proto
k8s.io/api/rbac/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.rbac.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/rbacv1" --api-domain rbac.authorization.k8s.io --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/admissionregistration/v1/generated.proto
k8s.io/api/admissionregistration/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.admissionregistration.v1 --go-package $(@D) --kubeproto-package "go.f110.dev/kubeproto/go/apis/admissionregistrationv1" --api-domain admissionregistration.k8s.io --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto
k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.apiextensions_apiserver.pkg.apis.apiextensions.v1 --go-package $(@D) --api-domain apiextensions --api-version v1 --all $(CURDIR)/vendor/$(@D)

.PHONY: sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto
sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package sigs.k8s.io.gateway_api.apis.v1alpha2 --go-package $(@D) --api-domain gateway --api-sub-group networking.k8s.io --api-version v1alpha2 --all $(CURDIR)/vendor/$(@D)

.PHONY: go/apis/metav1/metav1_kubeproto.generated.object.go
go/apis/metav1/metav1_kubeproto.generated.object.go: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):metav1_kubeproto --action_env=KUBEPROTO_OPTS=all
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/corev1/corev1_kubeproto.generated.object.go
go/apis/corev1/corev1_kubeproto.generated.object.go: k8s.io/api/core/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):corev1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/appsv1/appsv1_kubeproto.generated.object.go
go/apis/appsv1/appsv1_kubeproto.generated.object.go: k8s.io/api/apps/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):appsv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/batchv1/batchv1_kubeproto.generated.object.go
go/apis/batchv1/batchv1_kubeproto.generated.object.go: k8s.io/api/batch/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):batchv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/authenticationv1/authenticationv1_kubeproto.generated.object.go
go/apis/authenticationv1/authenticationv1_kubeproto.generated.object.go: k8s.io/api/authentication/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):authenticationv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/admissionv1/admissionv1_kubeproto.generated.object.go
go/apis/admissionv1/admissionv1_kubeproto.generated.object.go: k8s.io/api/admission/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):admissionv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/policyv1/policyv1_kubeproto.generated.object.go
go/apis/policyv1/policyv1_kubeproto.generated.object.go: k8s.io/api/policy/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):policyv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/networkingv1/networkingv1_kubeproto.generated.object.go
go/apis/networkingv1/networkingv1_kubeproto.generated.object.go: k8s.io/api/networking/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):networkingv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/rbacv1/rbacv1_kubeproto.generated.object.go
go/apis/rbacv1/rbacv1_kubeproto.generated.object.go: k8s.io/api/rbac/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):rbacv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/admissionregistrationv1/admissionregistrationv1_kubeproto.generated.object.go
go/apis/admissionregistrationv1/admissionregistrationv1_kubeproto.generated.object.go: k8s.io/api/admissionregistration/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(<D):admissionregistrationv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/k8sclient/go_client.generated.client.go
go/k8sclient/go_client.generated.client.go: k8s.io/api/core/v1/generated.proto \
		k8s.io/api/admission/v1/generated.proto \
		k8s.io/api/admissionregistration/v1/generated.proto \
		k8s.io/api/apps/v1/generated.proto \
		k8s.io/api/authentication/v1/generated.proto \
		k8s.io/api/batch/v1/generated.proto \
		k8s.io/api/networking/v1/generated.proto \
		k8s.io/api/policy/v1/generated.proto \
		k8s.io/api/rbac/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(@D):go_client
	cp ./bazel-bin/$(@D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/k8stestingclient/go_testingclient.generated.testingclient.go
go/k8stestingclient/go_testingclient.generated.testingclient.go: k8s.io/api/core/v1/generated.proto \
		k8s.io/api/admission/v1/generated.proto \
		k8s.io/api/admissionregistration/v1/generated.proto \
		k8s.io/api/apps/v1/generated.proto \
		k8s.io/api/authentication/v1/generated.proto \
		k8s.io/api/batch/v1/generated.proto \
		k8s.io/api/networking/v1/generated.proto \
		k8s.io/api/policy/v1/generated.proto \
		k8s.io/api/rbac/v1/generated.proto
	@mkdir -p $(@D)
	bazel build //$(@D):go_testingclient
	cp ./bazel-bin/$(@D)/$(@F) $(@D)
	@chmod 0644 $@
