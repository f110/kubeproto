package githubv1alpha1

import (
	"testing"

	"go.f110.dev/kubeproto/example/pkg/apis"
)

func Test(t *testing.T) {
	g := &apis.Grafana{
		Spec: apis.GrafanaSpec{
			FeatureGates: []string{"foo", "bar"},
			Volumes: []apis.Volume{
				{
					Name: "foo",
				},
			},
		},
	}
	newG := g.DeepCopy()

	g.Spec.FeatureGates[1] = "baz"
	g.Spec.Volumes = append(g.Spec.Volumes, apis.Volume{Name: "bar"})
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
