package goparser

import (
	"os"
	"path/filepath"
	"strings"
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

	t.Run("Multiple files", func(t *testing.T) {
		files := map[string]string{
			"a.go": `package api
type ConditionType string

const ConditionTypeAccepted ConditionType = "Accepted"`,
			"b.go": `package api
const ConditionTypeResolvedRefs ConditionType = "ResolvedRefs"`,
			"c.go": `package api
const ConditionTypeProgrammed ConditionType = "Programmed"`,
		}
		tmpDir := t.TempDir()
		for name, code := range files {
			if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(code), 0644); err != nil {
				t.Fatal(err)
			}
		}

		for range 20 {
			g := New()
			if err := g.AddDir(tmpDir, true); err != nil {
				t.Fatal(err)
			}
			var names []string
			for _, v := range g.enumValueCandidates["ConditionType"] {
				names = append(names, v.Name)
			}
			expect := []string{"ConditionTypeAccepted", "ConditionTypeResolvedRefs", "ConditionTypeProgrammed"}
			if strings.Join(names, ",") != strings.Join(expect, ",") {
				t.Fatalf("expect %v but %v", expect, names)
			}
		}
	})
}

func TestMergeStrategy(t *testing.T) {
	code := `package api

type PodSpec struct {
	// +listType=map
	// +listMapKey=containerPort
	// +listMapKey=protocol
	Ports []Port ` + "`json:\"ports,omitempty\"`" + `
	// +listType=atomic
	Args []string ` + "`json:\"args,omitempty\"`" + `
	// +listType=set
	Finalizers []string ` + "`json:\"finalizers,omitempty\"`" + `
	Containers []Container ` + "`json:\"containers\" patchStrategy:\"merge\" patchMergeKey:\"name\"`" + `
	Names []string ` + "`json:\"names,omitempty\" patchStrategy:\"merge\"`" + `
	// +listType=atomic
	Volumes []Volume ` + "`json:\"volumes,omitempty\" patchStrategy:\"merge,retainKeys\" patchMergeKey:\"name\"`" + `
	// +mapType=atomic
	NodeSelector map[string]string ` + "`json:\"nodeSelector,omitempty\"`" + `
	Labels map[string]string ` + "`json:\"labels,omitempty\"`" + `
	Hosts []string ` + "`json:\"hosts,omitempty\"`" + `
}

// +structType=atomic
type Port struct {
	ContainerPort int32 ` + "`json:\"containerPort\"`" + `
	Protocol string ` + "`json:\"protocol,omitempty\"`" + `
}

type Container struct {
	Name string ` + "`json:\"name\"`" + `
}

type Volume struct {
	Name string ` + "`json:\"name\"`" + `
}
`
	tmpDir := t.TempDir()
	g := New()
	g.SetProtoPackage("api")
	if err := os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(code), 0644); err != nil {
		t.Fatal(err)
	}
	if err := g.AddDir(tmpDir, true); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "generated.proto")
	if err := g.WriteFile(out); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(b)

	expects := []string{
		`repeated Port ports = 1 [(dev.f110.kubeproto.field) = {go_name: "Ports", api_field_name: "ports", list_type: LIST_TYPE_MAP, list_map_keys: ["containerPort", "protocol"], inline: false}];`,
		`repeated string args = 2 [(dev.f110.kubeproto.field) = {go_name: "Args", api_field_name: "args", list_type: LIST_TYPE_ATOMIC, inline: false}];`,
		`repeated string finalizers = 3 [(dev.f110.kubeproto.field) = {go_name: "Finalizers", api_field_name: "finalizers", list_type: LIST_TYPE_SET, inline: false}];`,
		`repeated Container containers = 4 [(dev.f110.kubeproto.field) = {go_name: "Containers", api_field_name: "containers", list_type: LIST_TYPE_MAP, list_map_keys: ["name"], inline: false}];`,
		`repeated string names = 5 [(dev.f110.kubeproto.field) = {go_name: "Names", api_field_name: "names", list_type: LIST_TYPE_SET, inline: false}];`,
		`repeated Volume volumes = 6 [(dev.f110.kubeproto.field) = {go_name: "Volumes", api_field_name: "volumes", list_type: LIST_TYPE_ATOMIC, inline: false}];`,
		`map<string, string> node_selector = 7 [(dev.f110.kubeproto.field) = {go_name: "NodeSelector", api_field_name: "nodeSelector", map_type: MAP_TYPE_ATOMIC, inline: false}];`,
		`map<string, string> labels = 8 [(dev.f110.kubeproto.field) = {go_name: "Labels", api_field_name: "labels", inline: false}];`,
		`repeated string hosts = 9 [(dev.f110.kubeproto.field) = {go_name: "Hosts", api_field_name: "hosts", inline: false}];`,
		"option (dev.f110.kubeproto.message) = {struct_type: STRUCT_TYPE_ATOMIC};",
	}
	for _, v := range expects {
		if !strings.Contains(generated, v) {
			t.Errorf("generated proto does not contain %q\n%s", v, generated)
		}
	}
	if strings.Count(generated, "dev.f110.kubeproto.message") != 1 {
		t.Errorf("message option is expected only for Port\n%s", generated)
	}
}

func TestUnresolvedPackage(t *testing.T) {
	cases := map[string]string{
		"Field":   "Ref unknown.Reference",
		"Pointer": "Ref *unknown.Reference",
		"Slice":   "Refs []unknown.Reference",
		"Map":     "Refs map[string]unknown.Reference",
	}
	for name, field := range cases {
		t.Run(name, func(t *testing.T) {
			code := `package api

import (
	"time"

	"example.com/unknown"
)

type Spec struct {
	Timeout time.Duration ` + "`json:\"timeout\"`" + `
	` + field + " `json:\"ref\"`" + `
}
`
			tmpDir := t.TempDir()
			if err := os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(code), 0644); err != nil {
				t.Fatal(err)
			}
			err := New().AddDir(tmpDir, true)
			if err == nil {
				t.Fatal("expect an error but got nil")
			}
			if !strings.Contains(err.Error(), "example.com/unknown") {
				t.Errorf("error does not contain the import path: %v", err)
			}
		})
	}

	t.Run("Resolved", func(t *testing.T) {
		code := `package api

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Spec struct {
	Timeout time.Duration ` + "`json:\"timeout\"`" + `
	Time metav1.Time ` + "`json:\"time\"`" + `
}
`
		tmpDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(code), 0644); err != nil {
			t.Fatal(err)
		}
		if err := New().AddDir(tmpDir, true); err != nil {
			t.Fatal(err)
		}
	})
}

func TestPredefinedType(t *testing.T) {
	code := `package api

import (
	corev1 "k8s.io/api/core/v1"
)

type PodResources struct {
	Limits corev1.ResourceList ` + "`json:\"limits,omitempty\"`" + `
	Requests *corev1.ResourceList ` + "`json:\"requests,omitempty\"`" + `
}
`
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(code), 0644); err != nil {
		t.Fatal(err)
	}
	g := New()
	g.SetProtoPackage("api")
	if err := g.AddDir(tmpDir, true); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "generated.proto")
	if err := g.WriteFile(out); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(b)

	expects := []string{
		`import "k8s.io/apimachinery/pkg/api/resource/generated.proto";`,
		`map<string, .k8s.io.apimachinery.pkg.api.resource.Quantity> limits = 1 [(dev.f110.kubeproto.field) = {go_name: "Limits", api_field_name: "limits", inline: false}];`,
		`map<string, .k8s.io.apimachinery.pkg.api.resource.Quantity> requests = 2 [(dev.f110.kubeproto.field) = {go_name: "Requests", api_field_name: "requests", inline: false}];`,
	}
	for _, v := range expects {
		if !strings.Contains(generated, v) {
			t.Errorf("generated proto does not contain %q\n%s", v, generated)
		}
	}
	if strings.Contains(generated, "k8s.io/api/core/v1/generated.proto") {
		t.Errorf("generated proto should not import k8s.io/api/core/v1\n%s", generated)
	}
}
