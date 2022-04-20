package apis

import (
	"testing"
)

func Test(t *testing.T) {
	g := &Grafana{
		Spec: GrafanaSpec{
			FeatureGates: []string{"foo", "bar"},
			Volumes: []Volume{
				{
					Name: "foo",
				},
			},
		},
	}
	newG := g.DeepCopy()

	g.Spec.FeatureGates[1] = "baz"
	g.Spec.Volumes = append(g.Spec.Volumes, Volume{Name: "bar"})
	g.Spec.Volumes[0].Name = "baz"

	if len(newG.Spec.FeatureGates) != 2 {
		t.Fatal("DeepCopyInto is wrong")
	}
	if g.Spec.FeatureGates[0] != "foo" {
		t.Error("DeepCopyInto is wrong")
	}
	if g.Spec.FeatureGates[1] != "baz" {
		t.Error("DeepCopyInto is wrong")
	}
	if newG.Spec.FeatureGates[0] != "foo" {
		t.Error("DeepCopyInto is wrong")
	}
	if newG.Spec.FeatureGates[1] != "bar" {
		t.Error("DeepCopyInto is wrong")
	}

	if len(newG.Spec.Volumes) != 1 {
		t.Fatal("DeepCopyInto is wrong")
	}
	if newG.Spec.Volumes[0].Name != "foo" {
		t.Error("DeepCopyInto is wrong")
	}
}
