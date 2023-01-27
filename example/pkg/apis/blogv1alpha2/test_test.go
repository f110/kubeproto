package blogv1alpha2

import (
	"testing"
)

func Test(t *testing.T) {
	g := &Blog{
		Spec: BlogSpec{
			Tags: []string{"foo", "bar"},
			Categories: []Category{
				{
					Name: "foo",
				},
			},
		},
	}
	newG := g.DeepCopy()

	g.Spec.Tags[1] = "baz"
	g.Spec.Categories = append(g.Spec.Categories, Category{Name: "bar"})
	g.Spec.Categories[0].Name = "baz"

	if len(newG.Spec.Tags) != 2 {
		t.Fatal("DeepCopyInto is wrong")
	}
	if g.Spec.Tags[0] != "foo" {
		t.Error("DeepCopyInto is wrong")
	}
	if g.Spec.Tags[1] != "baz" {
		t.Error("DeepCopyInto is wrong")
	}
	if newG.Spec.Tags[0] != "foo" {
		t.Error("DeepCopyInto is wrong")
	}
	if newG.Spec.Tags[1] != "bar" {
		t.Error("DeepCopyInto is wrong")
	}

	if len(newG.Spec.Categories) != 1 {
		t.Fatal("DeepCopyInto is wrong")
	}
	if newG.Spec.Categories[0].Name != "foo" {
		t.Error("DeepCopyInto is wrong")
	}
}
