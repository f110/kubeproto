package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestNewEnum(t *testing.T) {
	enum := NewEnum(
		&descriptorpb.FileDescriptorProto{
			Package: proto.String("testing.apis"),
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("go.f110.dev/kubeproto/internal/definition"),
			},
		},
		&descriptorpb.EnumDescriptorProto{
			Name: proto.String("Phase"),
			Value: []*descriptorpb.EnumValueDescriptorProto{
				{Name: proto.String("PHASE_CREATED"), Number: proto.Int32(0)},
				{Name: proto.String("PHASE_CREATING"), Number: proto.Int32(1)},
				{Name: proto.String("PHASE_CREATE_PENDING"), Number: proto.Int32(2)},
			},
		},
	)

	assert.Equal(t, ".testing.apis.Phase", enum.Name)
	assert.Equal(t, "Phase", enum.ShortName)
	if assert.Len(t, enum.Values, 3) {
		assert.Equal(t, "Created", enum.Values[0])
		assert.Equal(t, "Creating", enum.Values[1])
		assert.Equal(t, "CreatePending", enum.Values[2])
	}
}
