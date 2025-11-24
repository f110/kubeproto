BAZEL ?= bazel
GO ?= $(BAZEL) run @rules_go//go --

.PHONY: deps
deps:
	$(GO) mod tidy
	$(GO) mod vendor
	$(BAZEL) run //:gazelle

kube.pb.go: kube.proto
	$(BAZEL) build //:kubeproto_go_proto
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
	k8s.io/api/admissionregistration/v1/generated.proto \
	k8s.io/api/certificates/v1/generated.proto \
	k8s.io/api/authorization/v1/generated.proto \
	k8s.io/api/discovery/v1/generated.proto \
	k8s.io/api/autoscaling/v1/generated.proto \
	k8s.io/api/autoscaling/v2/generated.proto \
	k8s.io/api/coordination/v1/generated.proto \
	k8s.io/api/events/v1/generated.proto \
	k8s.io/api/scheduling/v1/generated.proto \
	k8s.io/api/storage/v1/generated.proto \
	k8s.io/api/apidiscovery/v2beta1/generated.proto

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
	go/apis/certificatesv1/certificatesv1_kubeproto.generated.object.go \
	go/apis/authorizationv1/authorizationv1_kubeproto.generated.object.go \
	go/apis/discoveryv1/discoveryv1_kubeproto.generated.object.go \
	go/apis/autoscalingv1/autoscalingv1_kubeproto.generated.object.go \
	go/apis/autoscalingv2/autoscalingv2_kubeproto.generated.object.go \
	go/apis/coordinationv1/coordinationv1_kubeproto.generated.object.go \
	go/apis/eventsv1/eventsv1_kubeproto.generated.object.go \
	go/apis/schedulingv1/schedulingv1_kubeproto.generated.object.go \
	go/apis/storagev1/storagev1_kubeproto.generated.object.go \
	go/apis/apidiscoveryv2beta1/apidiscoveryv2beta1_kubeproto.generated.object.go \
	go/k8sclient/go_client.generated.client.go \
	go/k8stestingclient/go_testingclient.generated.testingclient.go

.PHONY: go/k8sclient/go_client.generated.client.go
go/k8sclient/go_client.generated.client.go: k8s.io/api/core/v1/generated.proto \
		k8s.io/api/admission/v1/generated.proto \
		k8s.io/api/admissionregistration/v1/generated.proto \
		k8s.io/api/apps/v1/generated.proto \
		k8s.io/api/authentication/v1/generated.proto \
		k8s.io/api/batch/v1/generated.proto \
		k8s.io/api/networking/v1/generated.proto \
		k8s.io/api/policy/v1/generated.proto \
		k8s.io/api/rbac/v1/generated.proto \
		k8s.io/api/certificates/v1/generated.proto \
		k8s.io/api/authorization/v1/generated.proto \
		k8s.io/api/discovery/v1/generated.proto \
		k8s.io/api/autoscaling/v1/generated.proto \
		k8s.io/api/autoscaling/v2/generated.proto \
		k8s.io/api/coordination/v1/generated.proto \
		k8s.io/api/events/v1/generated.proto \
		k8s.io/api/scheduling/v1/generated.proto \
		k8s.io/api/storage/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(@D):go_client
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
		k8s.io/api/rbac/v1/generated.proto \
		k8s.io/api/certificates/v1/generated.proto \
		k8s.io/api/authorization/v1/generated.proto \
		k8s.io/api/discovery/v1/generated.proto \
		k8s.io/api/autoscaling/v1/generated.proto \
		k8s.io/api/autoscaling/v2/generated.proto \
		k8s.io/api/coordination/v1/generated.proto \
		k8s.io/api/events/v1/generated.proto \
		k8s.io/api/scheduling/v1/generated.proto \
		k8s.io/api/storage/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(@D):go_testingclient
	cp ./bazel-bin/$(@D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto
k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/apimachinery/pkg/api/resource/generated.proto
k8s.io/apimachinery/pkg/api/resource/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/apimachinery/pkg/util/intstr/generated.proto
k8s.io/apimachinery/pkg/util/intstr/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/apimachinery/pkg/runtime/generated.proto
k8s.io/apimachinery/pkg/runtime/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/core/v1/generated.proto
k8s.io/api/core/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/apps/v1/generated.proto
k8s.io/api/apps/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/batch/v1/generated.proto
k8s.io/api/batch/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/authentication/v1/generated.proto
k8s.io/api/authentication/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/admission/v1/generated.proto
k8s.io/api/admission/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/policy/v1/generated.proto
k8s.io/api/policy/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/networking/v1/generated.proto
k8s.io/api/networking/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/rbac/v1/generated.proto
k8s.io/api/rbac/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/admissionregistration/v1/generated.proto
k8s.io/api/admissionregistration/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/certificates/v1/generated.proto
k8s.io/api/certificates/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/authorization/v1/generated.proto
k8s.io/api/authorization/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/discovery/v1/generated.proto
k8s.io/api/discovery/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/autoscaling/v1/generated.proto
k8s.io/api/autoscaling/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/autoscaling/v2/generated.proto
k8s.io/api/autoscaling/v2/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/coordination/v1/generated.proto
k8s.io/api/coordination/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/events/v1/generated.proto
k8s.io/api/events/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/scheduling/v1/generated.proto
k8s.io/api/scheduling/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/storage/v1/generated.proto
k8s.io/api/storage/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/api/apidiscovery/v2beta1/generated.proto
k8s.io/api/apidiscovery/v2beta1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto
k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto
sigs.k8s.io/gateway-api/apis/v1alpha2/generated.proto:
	$(BAZEL) build //$(@D):gen
	mkdir -p $(@D)
	cp ./bazel-bin/$@ $(@D)

.PHONY: go/apis/metav1/metav1_kubeproto.generated.object.go
go/apis/metav1/metav1_kubeproto.generated.object.go: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):metav1_kubeproto --action_env=KUBEPROTO_OPTS=all
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/corev1/corev1_kubeproto.generated.object.go
go/apis/corev1/corev1_kubeproto.generated.object.go: k8s.io/api/core/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):corev1_kubeproto --action_env=KUBEPROTO_OPTS=all
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/appsv1/appsv1_kubeproto.generated.object.go
go/apis/appsv1/appsv1_kubeproto.generated.object.go: k8s.io/api/apps/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):appsv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/batchv1/batchv1_kubeproto.generated.object.go
go/apis/batchv1/batchv1_kubeproto.generated.object.go: k8s.io/api/batch/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):batchv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/authenticationv1/authenticationv1_kubeproto.generated.object.go
go/apis/authenticationv1/authenticationv1_kubeproto.generated.object.go: k8s.io/api/authentication/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):authenticationv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/admissionv1/admissionv1_kubeproto.generated.object.go
go/apis/admissionv1/admissionv1_kubeproto.generated.object.go: k8s.io/api/admission/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):admissionv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/policyv1/policyv1_kubeproto.generated.object.go
go/apis/policyv1/policyv1_kubeproto.generated.object.go: k8s.io/api/policy/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):policyv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/networkingv1/networkingv1_kubeproto.generated.object.go
go/apis/networkingv1/networkingv1_kubeproto.generated.object.go: k8s.io/api/networking/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):networkingv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/rbacv1/rbacv1_kubeproto.generated.object.go
go/apis/rbacv1/rbacv1_kubeproto.generated.object.go: k8s.io/api/rbac/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):rbacv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/admissionregistrationv1/admissionregistrationv1_kubeproto.generated.object.go
go/apis/admissionregistrationv1/admissionregistrationv1_kubeproto.generated.object.go: k8s.io/api/admissionregistration/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):admissionregistrationv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/certificatesv1/certificatesv1_kubeproto.generated.object.go
go/apis/certificatesv1/certificatesv1_kubeproto.generated.object.go: k8s.io/api/certificates/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):certificatesv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/authorizationv1/authorizationv1_kubeproto.generated.object.go
go/apis/authorizationv1/authorizationv1_kubeproto.generated.object.go: k8s.io/api/authorization/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):authorizationv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/discoveryv1/discoveryv1_kubeproto.generated.object.go
go/apis/discoveryv1/discoveryv1_kubeproto.generated.object.go: k8s.io/api/discovery/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):discoveryv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/autoscalingv1/autoscalingv1_kubeproto.generated.object.go
go/apis/autoscalingv1/autoscalingv1_kubeproto.generated.object.go: k8s.io/api/autoscaling/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):autoscalingv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/autoscalingv2/autoscalingv2_kubeproto.generated.object.go
go/apis/autoscalingv2/autoscalingv2_kubeproto.generated.object.go: k8s.io/api/autoscaling/v2/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):autoscalingv2_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/coordinationv1/coordinationv1_kubeproto.generated.object.go
go/apis/coordinationv1/coordinationv1_kubeproto.generated.object.go: k8s.io/api/coordination/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):coordinationv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/eventsv1/eventsv1_kubeproto.generated.object.go
go/apis/eventsv1/eventsv1_kubeproto.generated.object.go: k8s.io/api/events/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):eventsv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/schedulingv1/schedulingv1_kubeproto.generated.object.go
go/apis/schedulingv1/schedulingv1_kubeproto.generated.object.go: k8s.io/api/scheduling/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):schedulingv1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/storagev1/storagev1_kubeproto.generated.object.go
go/apis/storagev1/storagev1_kubeproto.generated.object.go: k8s.io/api/storage/v1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):storagev1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@

.PHONY: go/apis/apidiscoveryv2beta1/apidiscoveryv2beta1_kubeproto.generated.object.go
go/apis/apidiscoveryv2beta1/apidiscoveryv2beta1_kubeproto.generated.object.go: k8s.io/api/apidiscovery/v2beta1/generated.proto
	@mkdir -p $(@D)
	$(BAZEL) build //$(<D):apidiscoveryv2beta1_kubeproto
	cp ./bazel-bin/$(<D)/$(@F) $(@D)
	@chmod 0644 $@
