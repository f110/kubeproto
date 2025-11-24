load("@rules_go//go:def.bzl", "GoInfo")
load("@bazel_skylib//lib:paths.bzl", "paths")

def _gen_protobuf(ctx):
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
        )
    ]

gen_protobuf = rule(
    implementation = _gen_protobuf,
    attrs = {
        "srcs": attr.label_list(providers = [GoInfo]),
        "proto_package_name": attr.string(),
        "importpath": attr.string(),
        "api_domain": attr.string(),
        "api_sub_group": attr.string(),
        "api_version": attr.string(),
        "kubeproto_importpath": attr.string(),
        "all": attr.bool(default = False),
        "_cmd": attr.label(
            executable = True,
            cfg = "host",
            default = "//cmd/gen-go-to-protobuf",
        )
    }
)
