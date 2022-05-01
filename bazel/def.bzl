load("@bazel_skylib//lib:paths.bzl", "paths")
load("@io_bazel_rules_go//go:def.bzl", "go_context")
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
            default = "@com_google_protobuf//:protoc",
        ),
        "_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-crd",
        ),
        "_compiler_name": attr.string(default = "crd"),
    },
)

def _go_client(ctx):
    args = ctx.actions.args()
    args.add("--plugin", ("protoc-gen-%s=%s" % (ctx.attr._compiler_name, ctx.executable._compiler.path)))

    proto_files = []
    transitive_protos = []
    for src in ctx.attr.srcs:
        proto = src[ProtoInfo]
        transitive_protos.append(proto.transitive_imports)
        args.add_all(proto.transitive_proto_path, format_each = "--proto_path=%s")

        for s in proto.direct_sources:
            args.add(s.path)
            proto_files.append(s)

    out = ctx.actions.declare_file("%s.generated.client.go" % ctx.label.name)
    args.add("--client_out=%s:." % out.path)
    args.add("--client_opt=%s" % ctx.attr.importpath)

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

go_client = rule(
    implementation = _go_client,
    attrs = {
        "srcs": attr.label_list(providers = [ProtoInfo]),
        "importpath": attr.string(mandatory = True),
        "protoc": attr.label(
            executable = True,
            cfg = "host",
            default = "@com_google_protobuf//:protoc",
        ),
        "_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-client",
        ),
        "_compiler_name": attr.string(default = "client"),
    },
)

def _execute_protoc(ctx, compiler, compiler_name, suffix, srcs):
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

    ctx.actions.run(
        executable = ctx.executable.protoc,
        tools = [compiler],
        inputs = depset(
            direct = proto_files,
            transitive = transitive_protos,
        ),
        outputs = [out],
        arguments = [args],
    )

    return out

def _kubeproto_go_api(ctx):
    deepcopyOut = _execute_protoc(
        ctx,
        ctx.executable._deepcopy_compiler,
        ctx.attr._deepcopy_compiler_name,
        "generated.deepcopy.go",
        ctx.attr.srcs,
    )
    registerOut = _execute_protoc(
        ctx,
        ctx.executable._register_compiler,
        ctx.attr._register_compiler_name,
        "generated.register.go",
        ctx.attr.srcs,
    )

    return [DefaultInfo(files = depset([deepcopyOut, registerOut]))]

kubeproto_go_api = rule(
    implementation = _kubeproto_go_api,
    attrs = {
        "srcs": attr.label_list(providers = [ProtoInfo]),
        "importpath": attr.string(mandatory = True),
        "protoc": attr.label(
            executable = True,
            cfg = "host",
            default = "@com_google_protobuf//:protoc",
        ),
        "_deepcopy_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-deepcopy",
        ),
        "_deepcopy_compiler_name": attr.string(default = "deepcopy"),
        "_register_compiler": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/protoc-gen-register",
        ),
        "_register_compiler_name": attr.string(default = "register"),
    }
)