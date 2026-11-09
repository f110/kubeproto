package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldTag(t *testing.T) {
	cases := []struct {
		Name  string
		Field *Field
		Tag   string
	}{
		{
			Name:  "Required",
			Field: &Field{FieldName: "name"},
			Tag:   "`json:\"name\"`",
		},
		{
			Name:  "Optional",
			Field: &Field{FieldName: "name", Optional: true},
			Tag:   "`json:\"name,omitempty\"`",
		},
		{
			Name:  "Repeated",
			Field: &Field{FieldName: "tolerations", Repeated: true},
			Tag:   "`json:\"tolerations,omitempty\"`",
		},
		{
			Name:  "Inline",
			Field: &Field{Inline: true},
			Tag:   "`json:\",inline\"`",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Tag, tc.Field.Tag())
		})
	}
}
