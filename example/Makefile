.PHONY: build
build: githubv1alpha1 githubv1alpha2 miniov1alpha1 \
		pkg/client/github.generated.client.go \
		crd/github.crd.yaml

.PHONY: githubv1alpha1
githubv1alpha1: pkg/apis/githubv1alpha1/github.proto
	bazel build //example/pkg/apis/githubv1alpha1:github_proto_kubeproto
	@rm -f pkg/apis/$@/github_proto_kubeproto.generated.deepcopy.go pkg/apis/$@/github_proto_kubeproto.generated.register.go
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_kubeproto.generated.deepcopy.go pkg/apis/$@
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_kubeproto.generated.register.go pkg/apis/$@

.PHONY: githubv1alpha2
githubv1alpha2: pkg/apis/githubv1alpha2/github.proto
	bazel build //example/pkg/apis/githubv1alpha2:github_proto_kubeproto
	@rm -f pkg/apis/$@/github_proto_kubeproto.generated.deepcopy.go pkg/apis/$@/github_proto_kubeproto.generated.register.go
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_kubeproto.generated.deepcopy.go pkg/apis/$@
	cp ../bazel-bin/example/pkg/apis/$@/github_proto_kubeproto.generated.register.go pkg/apis/$@

.PHONY: miniov1alpha1
miniov1alpha1: pkg/apis/miniov1alpha1/minio.proto
	bazel build //example/pkg/apis/miniov1alpha1:minio_proto_kubeproto
	@rm -f pkg/apis/$@/minio_proto_kubeproto.generated.deepcopy.go pkg/apis/$@/minio_proto_kubeproto.generated.register.go
	cp ../bazel-bin/example/pkg/apis/$@/minio_proto_kubeproto.generated.deepcopy.go pkg/apis/$@
	cp ../bazel-bin/example/pkg/apis/$@/minio_proto_kubeproto.generated.register.go pkg/apis/$@

.PHONY: pkg/client/github.generated.client.go
pkg/client/github.generated.client.go: pkg/apis/githubv1alpha1/github.proto pkg/apis/githubv1alpha2/github.proto pkg/apis/miniov1alpha1/minio.proto
	bazel build //example/pkg/client:github
	@rm -f $@
	cp ../bazel-bin/example/pkg/client/github.generated.client.go $(@D)
	@chmod 644 $@

.PHONY: crd/github.crd.yaml
crd/github.crd.yaml: pkg/apis/githubv1alpha1/github.proto pkg/apis/githubv1alpha2/github.proto pkg/apis/miniov1alpha1/minio.proto
	bazel build //example/crd:github
	@rm -f $@
	cp ../bazel-bin/example/crd/github.crd.yaml $(@D)
	@chmod 644 $@

.PHONY: clean
clean:
	rm -f pkg/apis/githubv1alpha1/github_proto_kubeproto.generated.deepcopy.go
	rm -f pkg/apis/githubv1alpha1/github_proto_kubeproto.generated.register.go
	rm -f pkg/apis/githubv1alpha2/github_proto_kubeproto.generated.deepcopy.go
	rm -f pkg/apis/githubv1alpha2/github_proto_kubeproto.generated.register.go
	rm -f pkg/apis/miniov1alpha1/minio_proto_kubeproto.generated.deepcopy.go
	rm -f pkg/apis/miniov1alpha1/minio_proto_kubeproto.generated.register.go
	rm -f crd/github.crd.yaml