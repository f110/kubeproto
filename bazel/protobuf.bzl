load("@rules_go//go:def.bzl", "GoInfo")
load("@bazel_skylib//lib:paths.bzl", "paths")

KubeProtoLibrary = provider()

def _gen_protobuf_impl(ctx):
    import_path = ""
    srcs = []
    dir = ""
    for src in ctx.attr.srcs:
        info = src[GoInfo]
        if not import_path:
            import_path = info.importpath
            dir = paths.dirname(info.srcs[0].path)
        elif import_path != info.importpath:
            fail("Can not generate once with multiple packages")
        srcs.extend(info.srcs)

    args = ctx.actions.args()
    args.add("--proto-package=%s" % ctx.attr.proto_package_name)
    args.add("--go-package=%s" % ctx.attr.importpath)
    if ctx.attr.api_domain:
        args.add("--api-domain=%s" % ctx.attr.api_domain)
    if ctx.attr.api_sub_group:
        args.add("--api-sub-group=%s" % ctx.attr.api_sub_group)
    if ctx.attr.api_version:
        args.add("--api-version=%s" % ctx.attr.api_version)
    if ctx.attr.deps:
        for d in ctx.attr.deps:
            dep = d[KubeProtoLibrary]
            args.add("--imports=%s:%s:%s" % (dep.origin_go_package, dep.proto_package, dep.file_import_path))
    if ctx.attr.kubeproto_importpath:
        args.add("--kubeproto-package=%s" % ctx.attr.kubeproto_importpath)
    if ctx.attr.all:
        args.add("--all")
    out = ctx.actions.declare_file("generated.proto")
    args.add("--out=%s" % out.path)
    args.add(dir)
    ctx.actions.run(
        executable = ctx.executable._cmd,
        inputs = depset(
            direct = srcs,
        ),
        outputs = [out],
        arguments = [args]
    )
    return [
        DefaultInfo(
            files = depset([out]),
        ),
        KubeProtoLibrary(
            proto_package = ctx.attr.proto_package_name,
            go_package = ctx.attr.importpath,
            file_import_path = ctx.attr.dir,
            origin_go_package = import_path,
        )
    ]

_gen_protobuf = rule(
    implementation = _gen_protobuf_impl,
    attrs = {
        "srcs": attr.label_list(providers = [GoInfo]),
        "proto_package_name": attr.string(),
        "importpath": attr.string(),
        "api_domain": attr.string(),
        "api_sub_group": attr.string(),
        "api_version": attr.string(),
        "dir": attr.string(),
        "all": attr.bool(default = False),
        "deps": attr.label_list(providers = [KubeProtoLibrary]),
        "kubeproto_importpath": attr.string(),
        "_cmd": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/gen-go-to-protobuf",
        )
    }
)

def gen_protobuf(name, **kwargs):
    if not "dir" in kwargs:
        dir = native.package_name()
        kwargs["dir"] = dir

    _gen_protobuf(name = name, **kwargs)
