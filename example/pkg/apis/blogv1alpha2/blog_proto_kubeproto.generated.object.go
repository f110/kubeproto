package blogv1alpha2

import (
	metav1_1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "blog.f110.dev"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1alpha2"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha2"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Author{},
		&AuthorList{},
		&Blog{},
		&BlogList{},
		&Post{},
		&PostList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type PostPhase string

const (
	PostPhaseCREATED      PostPhase = "CREATED"
	PostPhasePROVISIONING PostPhase = "PROVISIONING"
	PostPhasePROVISIONED  PostPhase = "PROVISIONED"
)

type Author struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AuthorSpec   `json:"spec"`
	Status            AuthorStatus `json:"status"`
}

func (in *Author) DeepCopyInto(out *Author) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *Author) DeepCopy() *Author {
	if in == nil {
		return nil
	}
	out := new(Author)
	in.DeepCopyInto(out)
	return out
}

func (in *Author) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type AuthorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Author `json:"items"`
}

func (in *AuthorList) DeepCopyInto(out *AuthorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Author, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *AuthorList) DeepCopy() *AuthorList {
	if in == nil {
		return nil
	}
	out := new(AuthorList)
	in.DeepCopyInto(out)
	return out
}

func (in *AuthorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
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

type AuthorSpec struct {
}

func (in *AuthorSpec) DeepCopyInto(out *AuthorSpec) {
	*out = *in
}

func (in *AuthorSpec) DeepCopy() *AuthorSpec {
	if in == nil {
		return nil
	}
	out := new(AuthorSpec)
	in.DeepCopyInto(out)
	return out
}

type AuthorStatus struct {
}

func (in *AuthorStatus) DeepCopyInto(out *AuthorStatus) {
	*out = *in
}

func (in *AuthorStatus) DeepCopy() *AuthorStatus {
	if in == nil {
		return nil
	}
	out := new(AuthorStatus)
	in.DeepCopyInto(out)
	return out
}

type BlogSpec struct {
	// blog title
	Title          string               `json:"title"`
	AuthorSelector metav1.LabelSelector `json:"authorSelector"`
	// A list of all tags.
	// A tag is one of metadata of the post.
	Tags               []string                  `json:"tags"`
	Categories         []Category                `json:"categories"`
	ServiceAccountJSON *corev1.SecretKeySelector `json:"serviceAccountJSON,omitempty"`
	EditorSelector     LabelSelector             `json:"editorSelector"`
	IssuerRef          metav1_1.ObjectReference  `json:"issuerRef"`
}

func (in *BlogSpec) DeepCopyInto(out *BlogSpec) {
	*out = *in
	in.AuthorSelector.DeepCopyInto(&out.AuthorSelector)
	if in.Tags != nil {
		t := make([]string, len(in.Tags))
		copy(t, in.Tags)
		out.Tags = t
	}
	if in.Categories != nil {
		l := make([]Category, len(in.Categories))
		for i := range in.Categories {
			in.Categories[i].DeepCopyInto(&l[i])
		}
		out.Categories = l
	}
	if in.ServiceAccountJSON != nil {
		in, out := &in.ServiceAccountJSON, &out.ServiceAccountJSON
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	in.EditorSelector.DeepCopyInto(&out.EditorSelector)
	in.IssuerRef.DeepCopyInto(&out.IssuerRef)
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
	Ready              bool         `json:"ready"`
	ObservedGeneration int64        `json:"observedGeneration"`
	Url                string       `json:"url"`
	ObservedTime       *metav1.Time `json:"observedTime,omitempty"`
}

func (in *BlogStatus) DeepCopyInto(out *BlogStatus) {
	*out = *in
	if in.ObservedTime != nil {
		in, out := &in.ObservedTime, &out.ObservedTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
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
	Subject string   `json:"subject"`
	Authors []string `json:"authors"`
}

func (in *PostSpec) DeepCopyInto(out *PostSpec) {
	*out = *in
	if in.Authors != nil {
		t := make([]string, len(in.Authors))
		copy(t, in.Authors)
		out.Authors = t
	}
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
	Ready bool      `json:"ready"`
	Phase PostPhase `json:"phase"`
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

type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (in *Category) DeepCopyInto(out *Category) {
	*out = *in
}

func (in *Category) DeepCopy() *Category {
	if in == nil {
		return nil
	}
	out := new(Category)
	in.DeepCopyInto(out)
	return out
}

type LabelSelector struct {
	metav1.LabelSelector `json:",inline"`
	Namespace            string `json:"namespace,omitempty"`
}

func (in *LabelSelector) DeepCopyInto(out *LabelSelector) {
	*out = *in
	out.LabelSelector = in.LabelSelector
}

func (in *LabelSelector) DeepCopy() *LabelSelector {
	if in == nil {
		return nil
	}
	out := new(LabelSelector)
	in.DeepCopyInto(out)
	return out
}
