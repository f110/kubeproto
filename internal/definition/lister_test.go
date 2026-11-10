package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestResolveGoTypeOfTimestamp(t *testing.T) {
	files := new(protoregistry.Files)
	require.NoError(t, files.RegisterFile(timestamppb.File_google_protobuf_timestamp_proto))
	fileDesc, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:       new("testing/apis/blog.proto"),
		Package:    new("testing.apis"),
		Dependency: []string{"google/protobuf/timestamp.proto"},
		Syntax:     new("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: new("go.f110.dev/kubeproto/testing/apis"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: new("BlogStatus"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     new("published_at"),
						Number:   new(int32(1)),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: new(".google.protobuf.Timestamp"),
					},
				},
			},
		},
	}, files)
	require.NoError(t, err)
	require.NoError(t, files.RegisterFile(fileDesc))

	lister := NewLister([]string{"testing/apis/blog.proto"}, files, NewPackageNamespaceManager())
	msg := lister.GetMessages().Find("testing.apis.BlogStatus")
	require.NotNil(t, msg)
	require.Len(t, msg.Fields, 1)

	importPath, packageAlias, typ := lister.ResolveGoType("go.f110.dev/kubeproto/testing/apis", msg.Fields[0])
	assert.Equal(t, "go.f110.dev/kubeproto/go/apis/metav1", importPath)
	assert.Equal(t, "metav1", packageAlias)
	assert.Equal(t, "metav1.Time", typ)
}

func TestResolveGoTypeOfMapWithMessageValue(t *testing.T) {
	files := new(protoregistry.Files)
	resourceFile, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:    new("testing/resource/generated.proto"),
		Package: new("testing.resource"),
		Syntax:  new("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: new("k8s.io/apimachinery/pkg/api/resource"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: new("Quantity")},
		},
	}, files)
	require.NoError(t, err)
	require.NoError(t, files.RegisterFile(resourceFile))
	fileDesc, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:       new("testing/apis/pod.proto"),
		Package:    new("testing.apis"),
		Dependency: []string{"testing/resource/generated.proto"},
		Syntax:     new("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: new("go.f110.dev/kubeproto/testing/apis"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: new("PodResources"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     new("limits"),
						Number:   new(int32(1)),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName: new(".testing.apis.PodResources.LimitsEntry"),
					},
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name:    new("LimitsEntry"),
						Options: &descriptorpb.MessageOptions{MapEntry: new(true)},
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:   new("key"),
								Number: new(int32(1)),
								Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							},
							{
								Name:     new("value"),
								Number:   new(int32(2)),
								Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
								TypeName: new(".testing.resource.Quantity"),
							},
						},
					},
				},
			},
		},
	}, files)
	require.NoError(t, err)
	require.NoError(t, files.RegisterFile(fileDesc))

	nsm := NewPackageNamespaceManager()
	lister := NewLister([]string{"testing/apis/pod.proto"}, files, nsm)
	msg := lister.GetMessages().Find("testing.apis.PodResources")
	require.NotNil(t, msg)
	require.Len(t, msg.Fields, 1)

	importPath, packageAlias, typ := lister.ResolveGoType("go.f110.dev/kubeproto/testing/apis", msg.Fields[0])
	alias := nsm.Alias("k8s.io/apimachinery/pkg/api/resource")
	assert.Equal(t, "k8s.io/apimachinery/pkg/api/resource", importPath)
	assert.Equal(t, alias, packageAlias)
	assert.Equal(t, "map[string]"+alias+".Quantity", typ)
}
