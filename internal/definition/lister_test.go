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
