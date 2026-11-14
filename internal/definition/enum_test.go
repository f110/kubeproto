package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto"
)

func TestNewEnum(t *testing.T) {
	enum := NewEnum(
		&descriptorpb.FileDescriptorProto{
			Package: new("testing.apis"),
			Options: &descriptorpb.FileOptions{
				GoPackage: new("go.f110.dev/kubeproto/internal/definition"),
			},
		},
		&descriptorpb.EnumDescriptorProto{
			Name: new("Phase"),
			Value: []*descriptorpb.EnumValueDescriptorProto{
				{Name: new("PHASE_CREATED"), Number: new(int32(0))},
				{Name: new("PHASE_CREATING"), Number: new(int32(1))},
				{Name: new("PHASE_PROVISION_PENDING"), Number: new(int32(2))},
			},
		},
	)

	assert.Equal(t, ".testing.apis.Phase", enum.Name)
	assert.Equal(t, "Phase", enum.ShortName)
	assert.Equal(t, []*EnumValue{
		{Name: "Created", Value: "Created"},
		{Name: "Creating", Value: "Creating"},
		{Name: "ProvisionPending", Value: "ProvisionPending"},
	}, enum.Values)
}

func TestNewEnumFromEnumDescriptor(t *testing.T) {
	emptyValue := &descriptorpb.EnumValueOptions{}
	proto.SetExtension(emptyValue, kubeproto.E_Value, &kubeproto.EnumValue{Value: ""})
	customValue := &descriptorpb.EnumValueOptions{}
	proto.SetExtension(customValue, kubeproto.E_Value, &kubeproto.EnumValue{Value: "custom"})

	fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:    new("testing.proto"),
		Package: new("testing.apis"),
		Syntax:  new("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: new("go.f110.dev/kubeproto/internal/definition"),
		},
		EnumType: []*descriptorpb.EnumDescriptorProto{
			{
				Name: new("Phase"),
				Value: []*descriptorpb.EnumValueDescriptorProto{
					{Name: new("PHASE_DEFAULT"), Number: new(int32(0)), Options: emptyValue},
					{Name: new("PHASE_CREATING"), Number: new(int32(1)), Options: customValue},
					{Name: new("PHASE_PROVISION_PENDING"), Number: new(int32(2))},
				},
			},
		},
	}, nil)
	require.NoError(t, err)

	enum := NewEnumFromEnumDescriptor(fd.Enums().Get(0), fd, false)
	assert.Equal(t, "testing.apis.Phase", enum.Name)
	assert.Equal(t, "Phase", enum.ShortName)
	assert.Equal(t, []*EnumValue{
		{Name: "Default", Value: ""},
		{Name: "Custom", Value: "custom"},
		{Name: "ProvisionPending", Value: "ProvisionPending"},
	}, enum.Values)
}
