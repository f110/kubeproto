load("@bazel_skylib//lib:paths.bzl", "paths")
load("@io_bazel_rules_go//go:def.bzl", "go_context")
load("@rules_proto//proto:defs.bzl", "ProtoInfo")
load("@io_bazel_rules_go//proto:compiler.bzl", "GoProtoCompiler")

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