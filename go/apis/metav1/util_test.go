package metav1

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldsV1JSON(t *testing.T) {
	in := `{"manager":"kubectl-rollout","operation":"Update","fieldsType":"FieldsV1","fieldsV1":{"f:spec":{"f:template":{}}}}`

	var e ManagedFieldsEntry
	require.NoError(t, json.Unmarshal([]byte(in), &e))
	require.NotNil(t, e.FieldsV1)
	assert.Equal(t, `{"f:spec":{"f:template":{}}}`, string(e.FieldsV1.Raw))

	out, err := json.Marshal(e)
	require.NoError(t, err)
	assert.JSONEq(t, in, string(out))
}

func TestFieldsV1JSONNull(t *testing.T) {
	var f FieldsV1
	require.NoError(t, json.Unmarshal([]byte("null"), &f))
	assert.Nil(t, f.Raw)

	out, err := json.Marshal(f)
	require.NoError(t, err)
	assert.Equal(t, "null", string(out))
}

func TestDurationJSON(t *testing.T) {
	type spec struct {
		MaxAge *Duration `json:"max_age,omitempty"`
	}

	var s spec
	require.NoError(t, json.Unmarshal([]byte(`{"max_age":"1h"}`), &s))
	require.NotNil(t, s.MaxAge)
	assert.Equal(t, time.Hour, s.MaxAge.TimeDuration())

	out, err := json.Marshal(s)
	require.NoError(t, err)
	assert.JSONEq(t, `{"max_age":"1h0m0s"}`, string(out))

	assert.Error(t, json.Unmarshal([]byte(`{"max_age":"1x"}`), &s))
	assert.Error(t, json.Unmarshal([]byte(`{"max_age":3600}`), &s))
}
