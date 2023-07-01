package blogv1alpha1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "blog.f110.dev"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1alpha1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Blog{},
		&BlogList{},
		&Post{},
		&PostList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type Blog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              BlogSpec   `json:"spec"`
	Status            BlogStatus `json:"status"`
}

func (in *Blog) DeepCopyInto(out *Blog) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *Blog) DeepCopy() *Blog {
	if in == nil {
		return nil
	}
	out := new(Blog)
	in.DeepCopyInto(out)
	return out
}

func (in *Blog) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type BlogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Blog `json:"items"`
}

func (in *BlogList) DeepCopyInto(out *BlogList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Blog, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *BlogList) DeepCopy() *BlogList {
	if in == nil {
		return nil
	}
	out := new(BlogList)
	in.DeepCopyInto(out)
	return out
}

func (in *BlogList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type Post struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PostSpec   `json:"spec"`
	Status            PostStatus `json:"status"`
}

func (in *Post) DeepCopyInto(out *Post) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *Post) DeepCopy() *Post {
	if in == nil {
		return nil
	}
	out := new(Post)
	in.DeepCopyInto(out)
	return out
}

func (in *Post) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type PostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Post `json:"items"`
}

func (in *PostList) DeepCopyInto(out *PostList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Post, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *PostList) DeepCopy() *PostList {
	if in == nil {
		return nil
	}
	out := new(PostList)
	in.DeepCopyInto(out)
	return out
}

func (in *PostList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type BlogSpec struct {
	Title string `json:"title"`
}

func (in *BlogSpec) DeepCopyInto(out *BlogSpec) {
	*out = *in
}

func (in *BlogSpec) DeepCopy() *BlogSpec {
	if in == nil {
		return nil
	}
	out := new(BlogSpec)
	in.DeepCopyInto(out)
	return out
}

type BlogStatus struct {
	Ready bool `json:"ready"`
}

func (in *BlogStatus) DeepCopyInto(out *BlogStatus) {
	*out = *in
}

func (in *BlogStatus) DeepCopy() *BlogStatus {
	if in == nil {
		return nil
	}
	out := new(BlogStatus)
	in.DeepCopyInto(out)
	return out
}

type PostSpec struct {
	Subject string `json:"subject"`
}

func (in *PostSpec) DeepCopyInto(out *PostSpec) {
	*out = *in
}

func (in *PostSpec) DeepCopy() *PostSpec {
	if in == nil {
		return nil
	}
	out := new(PostSpec)
	in.DeepCopyInto(out)
	return out
}

type PostStatus struct {
	Ready bool `json:"ready"`
}

func (in *PostStatus) DeepCopyInto(out *PostStatus) {
	*out = *in
}

func (in *PostStatus) DeepCopy() *PostStatus {
	if in == nil {
		return nil
	}
	out := new(PostStatus)
	in.DeepCopyInto(out)
	return out
}
