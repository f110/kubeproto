load("@bazel_skylib//lib:paths.bzl", "paths")
load("@rules_go//go:def.bzl", "GoLibrary", "go_context")
load("@rules_proto//proto:defs.bzl", "ProtoInfo")

def _crd_proto_manifest(ctx):
    args = ctx.actions.args()

    proto_files = []
    transitive_protos = []
    import_path = []
    for src in ctx.attr.srcs:
        proto = src[ProtoInfo]
        transitive_protos.append(proto.transitive_imports)
        import_path.append(proto.transitive_proto_path)
        for s in proto.direct_sources:
            args.add(s.path)
            proto_files.append(s)

    args.add("--plugin", ("protoc-gen-%s=%s" % (ctx.attr._compiler_name, ctx.executable._compiler.path)))
    for i in import_path:
        args.add_all(i, format_each = "--proto_path=%s")

    out = ctx.actions.declare_file("%s.crd.yaml" % ctx.label.name)
    args.add("--crd_out=%s:." % out.path)
    ctx.actions.run(
        executable = ctx.executable.protoc,
        tools = [ctx.executable._compiler],
        inputs = depset(
            direct = proto_files,
            transitive = transitive_protos,
        ),
        outputs = [out],
        arguments = [args],
    )

    return [DefaultInfo(files = depset([out]))]

crd_proto_manifest = rule(
    implementation = _crd_proto_manifest,
    attrs = {
        "srcs": attr.label_list(providers = [ProtoInfo]),
        "protoc": attr.label(
            executable = True,
            cfg = "host",
            default = "@protobuf//:protoc",
        ),
        "_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-crd",
        ),
        "_compiler_name": attr.string(default = "crd"),
    },
)

K8SClient = provider()

def _go_client(ctx):
    go = go_context(ctx)
    opts = ctx.attr.importpath
    if ctx.attr.fqdn:
        opts += ",fqdn-set"
    out = _execute_protoc(
        ctx,
        ctx.executable._compiler,
        ctx.attr._compiler_name,
        "generated.client.go",
        ctx.attr.srcs,
        opts,
    )
    library = go.new_library(go, srcs = [out])
    source = go.library_to_source(go, ctx.attr, library, False)
    k8s_client = K8SClient(fqdn = ctx.attr.fqdn)

    return [
        library,
        source,
        DefaultInfo(
            files = depset([out]),
        ),
        k8s_client,
    ]

go_client = rule(
    implementation = _go_client,
    attrs = {
        "srcs": attr.label_list(providers = [ProtoInfo]),
        "importpath": attr.string(mandatory = True),
        "fqdn": attr.bool(default = False, doc = "Generate with FQDN name"),
        "protoc": attr.label(
            executable = True,
            cfg = "host",
            default = "@protobuf//:protoc",
        ),
        "_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-client",
        ),
        "_compiler_name": attr.string(default = "client"),
        "_go_context_data": attr.label(
            default = "@rules_go//:go_context_data",
        ),
    },
    toolchains = ["@rules_go//go:toolchain"],
)

def _go_testing_client(ctx):
    go = go_context(ctx)
    opts = "%s,%s" % (ctx.attr.importpath, ctx.attr.client[GoLibrary].importpath)
    if ctx.attr.client[K8SClient].fqdn:
        opts += ",fqdn-set"
    out = _execute_protoc(
        ctx,
        ctx.executable._compiler,
        ctx.attr._compiler_name,
        "generated.testingclient.go",
        ctx.attr.srcs,
        opts,
    )
    library = go.new_library(go, srcs = [out])
    source = go.library_to_source(go, ctx.attr, library, False)

    return [
        library,
        source,
        DefaultInfo(
            files = depset([out]),
        ),
    ]

go_testing_client = rule(
    implementation = _go_testing_client,
    attrs = {
        "srcs": attr.label_list(providers = [ProtoInfo]),
        "importpath": attr.string(mandatory = True),
        "client": attr.label(mandatory = True, providers = [GoLibrary, K8SClient]),
        "protoc": attr.label(
            executable = True,
            cfg = "host",
            default = "@protobuf//:protoc",
        ),
        "_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-testing-client",
        ),
        "_compiler_name": attr.string(default = "testing-client"),
        "_go_context_data": attr.label(
            default = "@rules_go//:go_context_data",
        ),
    },
    toolchains = ["@rules_go//go:toolchain"],
)

def _execute_protoc(ctx, compiler, compiler_name, suffix, srcs, opts = ""):
    args = ctx.actions.args()
    args.add("--plugin", ("protoc-gen-%s=%s" % (compiler_name, compiler.path)))

    proto_files = []
    transitive_protos = []
    for src in srcs:
        proto = src[ProtoInfo]
        transitive_protos.append(proto.transitive_imports)
        args.add_all(proto.transitive_proto_path, format_each = "--proto_path=%s")

        for s in proto.direct_sources:
            args.add(s.path)
            proto_files.append(s)

    out = ctx.actions.declare_file("%s.%s" % (ctx.label.name, suffix))
    args.add("--%s_out=%s:." % (compiler_name, out.path))
    if opts:
        args.add("--%s_opt=%s" % (compiler_name, opts))

    ctx.actions.run(
        executable = ctx.executable.protoc,
        tools = [compiler],
        inputs = depset(
            direct = proto_files,
            transitive = transitive_protos,
        ),
        outputs = [out],
        arguments = [args],
        use_default_shell_env = True,
    )

    return out

def _kubeproto_go_api(ctx):
    go = go_context(ctx)

    objectOut = _execute_protoc(
        ctx,
        ctx.executable._object_compiler,
        ctx.attr._object_compiler_name,
        "generated.object.go",
        ctx.attr.srcs,
    )
    library = go.new_library(go, srcs = [objectOut])
    source = go.library_to_source(go, ctx.attr, library, False)

    return [
        library,
        source,
        DefaultInfo(
            files = depset([objectOut]),
        ),
    ]

kubeproto_go_api = rule(
    implementation = _kubeproto_go_api,
    attrs = {
        "srcs": attr.label_list(providers = [ProtoInfo]),
        "importpath": attr.string(mandatory = True),
        "protoc": attr.label(
            executable = True,
            cfg = "host",
            default = "@protobuf//:protoc",
        ),
        "_object_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-object",
        ),
        "_object_compiler_name": attr.string(default = "object"),
        "_go_context_data": attr.label(
            default = "@rules_go//:go_context_data",
        ),
    },
    toolchains = ["@rules_go//go:toolchain"],
)
