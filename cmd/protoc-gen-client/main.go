package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"go.f110.dev/kubeproto/internal/k8s"
)

func genClient() error {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var input pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(buf, &input)
	if err != nil {
		return err
	}
	files, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{File: input.ProtoFile})
	if err != nil {
		return err
	}

	var outFile, importPath string
	var FQNDNSetName bool
	opt := input.GetParameter()
	if strings.Contains(opt, ",") {
		s := strings.Split(opt, ",")
		outFile = s[0]
		importPath = s[1]
		if len(s) > 2 {
			for _, v := range s[2:] {
				switch v {
				case "fqdn-set":
					FQNDNSetName = true
				}
			}
		}
	} else {
		outFile = opt
	}

	var res pluginpb.CodeGeneratorResponse
	supportedFeatures := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	res.SupportedFeatures = &supportedFeatures
	out := new(bytes.Buffer)
	g := k8s.NewClientGenerator(input.FileToGenerate, files)
	packageName := path.Base(filepath.Dir(outFile))
	if err := g.Generate(out, packageName, importPath, FQNDNSetName); err != nil {
		return err
	}
	res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(outFile),
		Content: proto.String(out.String()),
	})

	output, err := proto.Marshal(&res)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(output); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := genClient(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
