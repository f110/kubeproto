package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"go.f110.dev/kubeproto/internal/k8s"
)

func genDeepCopy() error {
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

	outFile := input.GetParameter()
	var res pluginpb.CodeGeneratorResponse
	supportedFeatures := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	res.SupportedFeatures = &supportedFeatures
	g, err := k8s.NewDeepCopyGenerator(input.FileToGenerate, files)
	if err != nil {
		return err
	}
	out := new(bytes.Buffer)
	if err := g.Generate(out); err != nil {
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
	if err := genDeepCopy(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
