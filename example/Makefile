.PHONY: build
build: pkg/apis/blogv1alpha1/blog_proto_kubeproto.generated.object.go \
		pkg/apis/blogv1alpha2/blog_proto_kubeproto.generated.object.go \
		pkg/client/k8s.generated.client.go \
		pkg/client/testingclient/k8s.generated.testingclient.go \
		crd/k8s.crd.yaml

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
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ \
		--proto-package github.com.cert_manager.cert_manager.apis.certmanagerv1 \
		--go-package github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1 \
		--api-sub-group cert-manager.io \
		--api-version v1 \
		--imports github.com/cert-manager/cert-manager/pkg/apis/meta/v1:github.com.cert_manager.cert_manager.apis.metav1:github.com/cert-manager/cert-manager/apis/metav1 \
		--imports github.com/cert-manager/cert-manager/pkg/apis/acme/v1:github.com.cert_manager.cert_manager.apis.acmev1:github.com/cert-manager/cert-manager/apis/acmev1 \
		--import-prefix example/proto \
		--all \
		$(CURDIR)/vendor/github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1

.PHONY: proto/github.com/cert-manager/cert-manager/apis/metav1/generated.proto
proto/github.com/cert-manager/cert-manager/apis/metav1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ \
		--proto-package github.com.cert_manager.cert_manager.apis.metav1 \
		--go-package github.com/cert-manager/cert-manager/pkg/apis/meta/v1 \
		--api-sub-group cert-manager.io \
		--api-version v1 \
		--all \
		$(CURDIR)/vendor/github.com/cert-manager/cert-manager/pkg/apis/meta/v1

.PHONY: proto/github.com/cert-manager/cert-manager/apis/acmev1/generated.proto
proto/github.com/cert-manager/cert-manager/apis/acmev1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ \
		--proto-package github.com.cert_manager.cert_manager.apis.acmev1 \
		--go-package github.com/cert-manager/cert-manager/pkg/apis/acme/v1 \
		--api-sub-group cert-manager.io \
		--api-version v1 \
		--imports github.com/cert-manager/cert-manager/pkg/apis/meta/v1:github.com.cert_manager.cert_manager.apis.metav1:github.com/cert-manager/cert-manager/apis/metav1 \
		--import-prefix example/proto \
		--all \
		$(CURDIR)/vendor/github.com/cert-manager/cert-manager/pkg/apis/acme/v1

.PHONY: proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/generated.proto
proto/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoringv1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ \
		--proto-package github.com.prometheus_operator.prometheus_operator.apis.monitoringv1 \
		--go-package github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1 \
		--api-sub-group coreos.com \
		--api-domain monitor \
		--api-version v1 \
		--all \
		$(CURDIR)/vendor/github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1

.PHONY: clean
clean:
	rm -f pkg/apis/githubv1alpha1/github_proto_kubeproto.generated.object.go
	rm -f pkg/apis/githubv1alpha2/github_proto_kubeproto.generated.object.go
	rm -f pkg/apis/miniov1alpha1/minio_proto_kubeproto.generated.object.go
	rm -f crd/github.crd.yaml
