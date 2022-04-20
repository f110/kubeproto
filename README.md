kubeproto
---

kubeproto is a plugin of Protocol Buffers to defined k8s API.

**WORK IN PROGRESS**

# Synopsis

github.proto

```protobuf
syntax = "proto3";
package example.apis;
option go_package = "go.f110.dev/kubeproto/example/pkg/apis";
option (dev.f110.kubeproto.k8s) = {group: "grafana.f110.dev", version: "v1alpha1"};

import "kube.proto";

message Grafana {
  GrafanaSpec   spec   = 1;
  GrafanaStatus status = 2;

  option (dev.f110.kubeproto.kind) = {};
}
```

BUILD

```
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "github_proto",
    srcs = ["github.proto"],
    deps = [
        "//:kubeproto",
    ],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "github_proto_deepcopy",
    compiler = "//bazel/go:deepcopy",
    importpath = "go.f110.dev/kubeproto/example/pkg/apis",
    proto = ":github_proto",
    deps = [ # deps is required
        "//example/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:meta",
        "//example/vendor/k8s.io/apimachinery/pkg/runtime",
    ],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "github_proto_register",
    compiler = "//bazel/go:register",
    importpath = "go.f110.dev/kubeproto/example/pkg/apis",
    proto = ":github_proto",
    visibility = ["//visibility:public"],
    embed = [":github_proto_deepcopy"],
    deps = [
        "//example/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:meta",
        "//example/vendor/k8s.io/apimachinery/pkg/runtime",
        "//example/vendor/k8s.io/apimachinery/pkg/runtime/schema",
    ],
)
```

# Author

Fumihiro Ito