package goparser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnum(t *testing.T) {
	t.Run("With type", func(t *testing.T) {
		code := `package api
type PolicyType string

const (
	// PolicyTypeIngress is a NetworkPolicy that affects ingress traffic on selected pods
	PolicyTypeIngress PolicyType = "Ingress"
	// PolicyTypeEgress is a NetworkPolicy that affects egress traffic on selected pods
	PolicyTypeEgress PolicyType = "Egress"
)`
		tmpDir := t.TempDir()
		g := New()
		err := os.WriteFile(filepath.Join(tmpDir, "enum.go"), []byte(code), 0644)
		if err != nil {
			t.Fatal(err)
		}
		err = g.AddDir(tmpDir, true)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := g.enumValueCandidates["PolicyType"]; !ok {
			t.Fatal("failed to parse enum")
		}
		if len(g.enumValueCandidates["PolicyType"]) != 2 {
			t.Errorf("found enum and expect two elements but %d elements.", len(g.enumValueCandidates["PolicyType"]))
		}
	})

	t.Run("Without type", func(t *testing.T) {
		code := `package api
// PathType represents the type of path referred to by a HTTPIngressPath.
// +enum
type PathType string

const (
	PathTypeExact = PathType("Exact")
	PathTypePrefix = PathType("Prefix")
	PathTypeImplementationSpecific = PathType("ImplementationSpecific")
)`

		tmpDir := t.TempDir()
		g := New()
		err := os.WriteFile(filepath.Join(tmpDir, "enum.go"), []byte(code), 0644)
		if err != nil {
			t.Fatal(err)
		}
		err = g.AddDir(tmpDir, true)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := g.enumValueCandidates["PathType"]; !ok {
			t.Fatal("failed to parse enum")
		}
		if len(g.enumValueCandidates["PathType"]) != 3 {
			t.Errorf("found enum and expect three elements but %d elements.", len(g.enumValueCandidates["PathType"]))
		}
	})
}
