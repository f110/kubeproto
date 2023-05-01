package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"go.f110.dev/kubeproto/internal/goparser"
)

func genGoToProtobuf(args []string) error {
	var out, protoPackage, goPackage, apiDomain, apiSubGroup, apiVersion, importPrefix, kubeprotoPackage string
	var imports []string
	var allStructs bool
	fs := pflag.NewFlagSet("gen-to-to-protobuf", pflag.PanicOnError)
	fs.StringVar(&out, "out", "", "Output file")
	fs.StringVar(&protoPackage, "proto-package", "", "Protobuf package name")
	fs.StringVar(&goPackage, "go-package", "", "GO package name")
	fs.BoolVar(&allStructs, "all", false, "Generate protobuf for all structs except marked as without generation")
	fs.StringVar(&apiDomain, "api-domain", "", "API domain name (e,g, f110.dev)")
	fs.StringVar(&apiSubGroup, "api-sub-group", "", "API sub group (e,g, minio)")
	fs.StringVar(&apiVersion, "api-version", "", "API version")
	fs.StringSliceVar(&imports, "imports", nil, "Imported packages. The value consists three element. "+
		"The first element is the name of go package. "+
		"The second element is the package name in protobuf. "+
		"The third element is the import path.")
	fs.StringVar(&importPrefix, "import-prefix", "", "Import path prefix")
	fs.StringVar(&kubeprotoPackage, "kubeproto-package", "", "")
	fs.MarkHidden("kubeproto-package")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if out == "" {
		return errors.New("--out is mandatory")
	}

	for _, v := range fs.Args() {
		g := goparser.New()
		g.SetProtoPackage(protoPackage)
		g.SetGoPackage(goPackage)
		g.SetAPIDomain(apiDomain, apiSubGroup)
		g.SetAPIVersion(apiVersion)
		g.AddImport(imports...)
		g.SetImportPrefix(importPrefix)
		g.SetKubeprotoPackage(kubeprotoPackage)
		if err := g.AddDir(v, allStructs); err != nil {
			return err
		}
		if err := g.WriteFile(out); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := genGoToProtobuf(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
