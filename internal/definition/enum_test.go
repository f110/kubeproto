package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/descriptorpb"
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
	if assert.Len(t, enum.Values, 3) {
		assert.Equal(t, "Created", enum.Values[0])
		assert.Equal(t, "Creating", enum.Values[1])
		assert.Equal(t, "ProvisionPending", enum.Values[2])
	}
}
