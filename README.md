kubeproto
---

kubeproto is a plugin for Protocol Buffers to define k8s API.

The code generated by these plugins is incompatible with the code generated by code-generator.

**Status: Stabilizing API**

# Synopsis

## With Bazel

github.proto

```protobuf
syntax = "proto3";
package example.apis.githubv1alpha1;
option go_package = "go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1";
option (dev.f110.kubeproto.k8s) = {
  domain: "f110.dev",
  sub_group: "grafana",
  version: "v1alpha1",
};

import "kube.proto";

message Grafana {
  GrafanaSpec   spec   = 1;
  GrafanaStatus status = 2;
        
  option (dev.f110.kubeproto.kind) = {};
}
```

BUILD at pkg/apis/GROUP_AND_VERSION for generating deepcopy and register

```
load("@rules_proto//proto:defs.bzl", "proto_library")
load("//bazel:def.bzl", "kubeproto_go_api")

proto_library(
    name = "github_proto",
    srcs = ["github.proto"],
    deps = [
        "//:kubeproto",
    ],
    visibility = ["//visibility:public"],
)

kubeproto_go_api(
    name = "github_proto_kubeproto",
    srcs = [":github_proto"],
    importpath = "go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1",
)
```

BUILD file for generating for test

```
load("//bazel:def.bzl", "go_testing_client")

go_testing_client(
    name = "github",
    srcs = [
        "//example/pkg/apis/githubv1alpha1:github_proto",
        "//example/pkg/apis/githubv1alpha2:github_proto",
        "//example/pkg/apis/miniov1alpha1:minio_proto",
    ],
    importpath = "go.f110.dev/kubeproto/example/pkg/client/testingclient",
    client = "//example/pkg/client:github",
    visibility = ["//visibility:public"],
)
```

BUILD file for generating CustomResourceDefinition

```
load("//bazel/crd:def.bzl", "crd_proto_manifest")

crd_proto_manifest(
    name = "github",
    srcs = [
        "//example/pkg/apis/githubv1alpha1:github_proto",
        "//example/pkg/apis/githubv1alpha2:github_proto",
    ],
    visibility = ["//visibility:public"],
)
```

# How to use generated client

```go
cfg, err := rest.InClusterConfig()
if err != nil {
    return nil, err
}
apiClient, err := client.NewSet(cfg)
if err != nil {
    return nil, err
}

factory := client.NewInformerFactory(apiClient, client.NewInformerCache(), metav1.NamespaceAll, 30*time.Second)
githubInformers := client.NewGithubV1alpha1Informer(factory.Cache(), apiClient.GiothubV1alpha1, metav1.NamespaceAll, 30*time.Second)
githubGrafanaInformer := githubInformers.GrafanaInformer()
grafanaLister := githubInformers.GrafanaLister()
```

# Why use the extension number for internal?

These plugins are intended to use my projects.

If you want to use these plugins in your project, I would consider registering with Global Extension Registry.
Please feel free to file an issue! 

# Author

Fumihiro Ito
