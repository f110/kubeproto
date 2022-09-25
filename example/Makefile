.PHONY: build
build: githubv1alpha1 githubv1alpha2 miniov1alpha1 \
		pkg/client/github.generated.client.go \
		crd/github.crd.yaml

.PHONY: githubv1alpha1
githubv1alpha1: pkg/apis/githubv1alpha1/github.proto
	bazel build //example/pkg/apis/githubv1alpha1:github_proto_kubeproto
	@rm -f pkg/apis/$@/github_proto_kubeproto.generated.object.go
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_kubeproto.generated.object.go pkg/apis/$@

.PHONY: githubv1alpha2
githubv1alpha2: pkg/apis/githubv1alpha2/github.proto
	bazel build //example/pkg/apis/githubv1alpha2:github_proto_kubeproto
	@rm -f pkg/apis/$@/github_proto_kubeproto.generated.object.go
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_kubeproto.generated.object.go pkg/apis/$@

.PHONY: miniov1alpha1
miniov1alpha1: pkg/apis/miniov1alpha1/minio.proto
	bazel build //example/pkg/apis/miniov1alpha1:minio_proto_kubeproto
	@rm -f pkg/apis/$@/minio_proto_kubeproto.generated.object.go
	cp ../bazel-bin/example/pkg/apis/$@/minio_proto_kubeproto.generated.object.go pkg/apis/$@

.PHONY: pkg/client/github.generated.client.go
pkg/client/github.generated.client.go: pkg/apis/githubv1alpha1/github.proto pkg/apis/githubv1alpha2/github.proto pkg/apis/miniov1alpha1/minio.proto
	bazel build //example/pkg/client:github
	@rm -f $@
	cp ../bazel-bin/example/pkg/client/github.generated.client.go $(@D)
	@chmod 644 $@

.PHONY: pkg/client/testingclient/github.generated.testingclient.go
pkg/client/testingclient/github.generated.testingclient.go: pkg/apis/githubv1alpha1/github.proto pkg/apis/githubv1alpha2/github.proto pkg/apis/miniov1alpha1/minio.proto
	bazel build //example/pkg/client/testingclient:github
	@rm -f $@
	cp ../bazel-bin/example/pkg/client/testingclient/github.generated.testingclient.go $(@D)
	@chmod 644 $@

.PHONY: crd/github.crd.yaml
crd/github.crd.yaml: pkg/apis/githubv1alpha1/github.proto pkg/apis/githubv1alpha2/github.proto pkg/apis/miniov1alpha1/minio.proto
	bazel build //example/crd:github
	@rm -f $@
	cp ../bazel-bin/example/crd/github.crd.yaml $(@D)
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

.PHONY: clean
clean:
	rm -f pkg/apis/githubv1alpha1/github_proto_kubeproto.generated.object.go
	rm -f pkg/apis/githubv1alpha2/github_proto_kubeproto.generated.object.go
	rm -f pkg/apis/miniov1alpha1/minio_proto_kubeproto.generated.object.go
	rm -f crd/github.crd.yaml
	rm -f proto/certmanageracmev1.proto proto/certmanagermetav1.proto proto/certmanagerv1.proto
