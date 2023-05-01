package stringsutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUpperCamelCase(t *testing.T) {
	assert.Equal(t, "FooBar", ToUpperCamelCase("foo_bar"))
	assert.Equal(t, "CertManagerIo", ToUpperCamelCase("cert-manager.io"))
	assert.Equal(t, "FooBar", ToUpperCamelCase("FooBar"))
	assert.Equal(t, "Kubernetes", ToUpperCamelCase("kubernetes"))
}

func TestToLowerCamelCase(t *testing.T) {
	assert.Equal(t, "fooBar", ToLowerCamelCase("foo_bar"))
}

func TestToUpperSnakeCase(t *testing.T) {
	assert.Equal(t, "FOO", ToUpperSnakeCase("Foo"))
	assert.Equal(t, "FOO_BAR", ToUpperSnakeCase("FooBar"))
	assert.Equal(t, "HTTP_01", ToUpperSnakeCase("HTTP-01"))
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
	cases := []struct {
		In    string
		Split []string
	}{
		{In: "UserAdmin", Split: []string{"User", "Admin"}},
		{In: "APIGroup", Split: []string{"API", "Group"}},
		{In: "UserUIDGroup", Split: []string{"User", "UID", "Group"}},
		{In: "AdminUserUIDGroup", Split: []string{"Admin", "User", "UID", "Group"}},
		{In: "StorageVersionHash", Split: []string{"Storage", "Version", "Hash"}},
		{In: "ServerAddressByClientCIDRs", Split: []string{"Server", "Address", "By", "Client", "CIDRs"}},
		{In: "HTTP-01", Split: []string{"HTTP", "01"}},
		{In: "HTTP_HTTPS", Split: []string{"HTTP", "HTTPS"}},
	}

	for _, tc := range cases {
		t.Run(tc.In, func(t *testing.T) {
			assert.Equal(t, tc.Split, splitString(tc.In))
		})
	}
}

func TestInsert(t *testing.T) {
	assert.Equal(t, []string{"UserAdmin", "Admin"}, insert([]string{"UserAdmin"}, 1, "Admin"))
	assert.Equal(t, []string{"User", "Full", "FullName"}, insert([]string{"User", "FullName"}, 1, "Full"))
}

func TestIsSnakeCase(t *testing.T) {
	assert.True(t, IsSnakeCase("foo_bar"))
	assert.False(t, IsSnakeCase("FooBar"))
	assert.False(t, IsSnakeCase("cert-manager.io"))
}

func TestIsCamelCase(t *testing.T) {
	assert.True(t, IsCamelCase("FooBar"))
	assert.False(t, IsCamelCase("foo_bar"))
	assert.False(t, IsCamelCase("cert-manager.io"))
	assert.False(t, IsCamelCase("kubernetes"))
}
