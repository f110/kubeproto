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
	assert.Equal(t, "uid", ToLowerSnakeCase("UID"))
	assert.Equal(t, "api_group", ToLowerSnakeCase("APIGroup"))
	assert.Equal(t, "storage_version_hash", ToLowerSnakeCase("StorageVersionHash"))
	assert.Equal(t, "server_address_by_client_cidrs", ToLowerSnakeCase("ServerAddressByClientCIDRs"))
}

func TestSplitString(t *testing.T) {
	assert.Equal(t, []string{"User", "Admin"}, splitString("UserAdmin"))
	assert.Equal(t, []string{"API", "Group"}, splitString("APIGroup"))
	assert.Equal(t, []string{"User", "UID", "Group"}, splitString("UserUIDGroup"))
	assert.Equal(t, []string{"Admin", "User", "UID", "Group"}, splitString("AdminUserUIDGroup"))
	assert.Equal(t, []string{"Storage", "Version", "Hash"}, splitString("StorageVersionHash"))
	assert.Equal(t, []string{"Server", "Address", "By", "Client", "CIDRs"}, splitString("ServerAddressByClientCIDRs"))
}

func TestInsert(t *testing.T) {
	assert.Equal(t, []string{"UserAdmin", "Admin"}, insert([]string{"UserAdmin"}, 1, "Admin"))
	assert.Equal(t, []string{"User", "Full", "FullName"}, insert([]string{"User", "FullName"}, 1, "Full"))
}
