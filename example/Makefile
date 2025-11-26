BAZEL ?= bazel

.PHONY: build
build: pkg/apis/blogv1alpha1/blog_proto_kubeproto.generated.object.go \
		pkg/apis/blogv1alpha2/blog_proto_kubeproto.generated.object.go \
		pkg/client/k8s.generated.client.go \
		pkg/client/testingclient/k8s.generated.testingclient.go \
		crd/blog.crd.yaml

.PHONY: build-thirdparty
build-thirdparty: proto/github.com/cert-manager/cert-manager/apis/metav1/generated.proto \
		proto/github.com/cert-manager/cert-manager/apis/acmev1/generated.proto \
		proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/generated.proto \
		proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/generated.proto \
		proto/github.com/cert-manager/cert-manager/apis/metav1/metav1_kubeproto.generated.object.go \
		proto/github.com/cert-manager/cert-manager/apis/acmev1/acmev1_kubeproto.generated.object.go \
		proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/certmanagerv1_kubeproto.generated.object.go \
		proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/monitoringv1_kubeproto.generated.object.go \
		pkg/thirdpartyclient/thirdpartyclient.generated.client.go \
		pkg/thirdpartyclient/testingthirdpartyclient/testingthirdpartyclient.generated.testingclient.go

.PHONY: pkg/apis/blogv1alpha1/blog_proto_kubeproto.generated.object.go
pkg/apis/blogv1alpha1/blog_proto_kubeproto.generated.object.go:
	bazel build //example/pkg/apis/blogv1alpha1:blog_proto_kubeproto
	@rm -f $@
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: pkg/apis/blogv1alpha2/blog_proto_kubeproto.generated.object.go
pkg/apis/blogv1alpha2/blog_proto_kubeproto.generated.object.go:
	bazel build //example/pkg/apis/blogv1alpha2:blog_proto_kubeproto
	@rm -f $@
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: pkg/client/k8s.generated.client.go
pkg/client/k8s.generated.client.go: pkg/apis/blogv1alpha1/blog.proto pkg/apis/blogv1alpha2/blog.proto
	bazel build //example/pkg/client:k8s
	@rm -f $@
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: pkg/client/testingclient/k8s.generated.testingclient.go
pkg/client/testingclient/k8s.generated.testingclient.go: pkg/apis/blogv1alpha1/blog.proto pkg/apis/blogv1alpha2/blog.proto
	bazel build //example/pkg/client/testingclient:k8s
	@rm -f $@
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: crd/blog.crd.yaml
crd/blog.crd.yaml: pkg/apis/blogv1alpha1/blog.proto pkg/apis/blogv1alpha2/blog.proto
	bazel build //example/crd:blog
	@rm -f $@
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: pkg/thirdpartyclient/thirdpartyclient.generated.client.go
pkg/thirdpartyclient/thirdpartyclient.generated.client.go:
	bazel build //example/pkg/thirdpartyclient:thirdpartyclient
	@rm -f $@
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: pkg/thirdpartyclient/testingthirdpartyclient/testingthirdpartyclient.generated.testingclient.go
pkg/thirdpartyclient/testingthirdpartyclient/testingthirdpartyclient.generated.testingclient.go:
	bazel build //example/pkg/thirdpartyclient/testingthirdpartyclient
	@rm -f $@
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/generated.proto
proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/generated.proto:
	$(BAZEL) build //example/$(@D):gen
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/cert-manager/cert-manager/apis/metav1/generated.proto
proto/github.com/cert-manager/cert-manager/apis/metav1/generated.proto:
	$(BAZEL) build //example/$(@D):gen
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/cert-manager/cert-manager/apis/acmev1/generated.proto
proto/github.com/cert-manager/cert-manager/apis/acmev1/generated.proto:
	$(BAZEL) build //example/$(@D):gen
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/generated.proto
proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/generated.proto:
	$(BAZEL) build //example/$(@D):gen
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/cert-manager/cert-manager/apis/acmev1/acmev1_kubeproto.generated.object.go
proto/github.com/cert-manager/cert-manager/apis/acmev1/acmev1_kubeproto.generated.object.go:
	$(BAZEL) build //example/$(@D):acmev1_kubeproto
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/certmanagerv1_kubeproto.generated.object.go
proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/certmanagerv1_kubeproto.generated.object.go:
	$(BAZEL) build //example/$(@D):certmanagerv1_kubeproto
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/cert-manager/cert-manager/apis/metav1/metav1_kubeproto.generated.object.go
proto/github.com/cert-manager/cert-manager/apis/metav1/metav1_kubeproto.generated.object.go:
	$(BAZEL) build //example/$(@D):metav1_kubeproto
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/monitoringv1_kubeproto.generated.object.go
proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/monitoringv1_kubeproto.generated.object.go:
	$(BAZEL) build //example/$(@D):monitoringv1_kubeproto
	mkdir -p $(@D)
	cp ../bazel-bin/example/$@ $(@D)
	@chmod 644 $@

.PHONY: clean
clean:
	rm -f pkg/apis/githubv1alpha1/github_proto_kubeproto.generated.object.go
	rm -f pkg/apis/githubv1alpha2/github_proto_kubeproto.generated.object.go
	rm -f pkg/apis/miniov1alpha1/minio_proto_kubeproto.generated.object.go
	rm -f crd/github.crd.yaml
	rm -f proto/github.com/cert-manager/cert-manager/apis/metav1/generated.proto
	rm -f proto/github.com/cert-manager/cert-manager/apis/acmev1/generated.proto
	rm -f proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/generated.proto
	rm -f proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/generated.proto
	rm -f proto/github.com/cert-manager/cert-manager/apis/metav1/metav1_kubeproto.generated.object.go
	rm -f proto/github.com/cert-manager/cert-manager/apis/acmev1/acmev1_kubeproto.generated.object.go
	rm -f proto/github.com/cert-manager/cert-manager/apis/certmanagerv1/certmanagerv1_kubeproto.generated.object.go
	rm -f proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/monitoringv1_kubeproto.generated.object.go
	rm -f pkg/thirdpartyclient/thirdpartyclient.generated.client.go
	rm -f pkg/thirdpartyclient/testingthirdpartyclient/testingthirdpartyclient.generated.testingclient.go
