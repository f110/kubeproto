package stringsutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUpperCamelCase(t *testing.T) {
	assert.Equal(t, "FooBar", ToUpperCamelCase("foo_bar"))
}

func TestToLowerCamelCase(t *testing.T) {
	assert.Equal(t, "fooBar", ToLowerCamelCase("foo_bar"))
}

func TestToUpperSnakeCase(t *testing.T) {
	assert.Equal(t, "FOO", ToUpperSnakeCase("Foo"))
	assert.Equal(t, "FOO_BAR", ToUpperSnakeCase("FooBar"))
}

func TestToLowerSnakeCase(t *testing.T) {
	assert.Equal(t, "foo", ToLowerSnakeCase("Foo"))
	assert.Equal(t, "foo_bar", ToLowerSnakeCase("FooBar"))
}
