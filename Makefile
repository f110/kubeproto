.PHONY: deps
deps:
	go mod tidy
	go mod vendor
	bazel run //:gazelle

kube.pb.go: kube.proto
	bazel build //:kubeproto_go_proto
	@cp bazel-bin/kubeproto_go_proto_/go.f110.dev/kubeproto/kube.pb.go ./
	@chmod 644 $@

.PHONY: k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto
k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(shell pwd)/$@ --proto-package k8s.io.apimachinery.pkg.apis.meta.v1 --go-package $(@D) --all $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/apimachinery/pkg/runtime/generated.proto
k8s.io/apimachinery/pkg/runtime/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(shell pwd)/$@ --proto-package k8s.io.apimachinery.pkg.runtime --go-package $(@D) $(CURDIR)/vendor/$(@D)

.PHONY: k8s.io/api/core/v1/generated.proto
k8s.io/api/core/v1/generated.proto:
	mkdir -p $(@D)
	bazel run //cmd/gen-go-to-protobuf -- --out $(CURDIR)/$@ --proto-package k8s.io.api.core.v1 --go-package $(@D) --all $(CURDIR)/vendor/$(@D)