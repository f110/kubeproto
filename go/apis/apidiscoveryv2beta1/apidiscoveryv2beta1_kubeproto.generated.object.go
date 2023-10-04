package apidiscoveryv2beta1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "apidiscovery.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v2beta1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v2beta1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&APIGroupDiscovery{},
		&APIGroupDiscoveryList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type DiscoveryFreshness string

const (
	DiscoveryFreshnessCurrent DiscoveryFreshness = "Current"
	DiscoveryFreshnessStale   DiscoveryFreshness = "Stale"
)

type ResourceScope string

const (
	ResourceScopeCluster    ResourceScope = "Cluster"
	ResourceScopeNamespaced ResourceScope = "Namespaced"
)

type APIGroupDiscovery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// versions are the versions supported in this group. They are sorted in descending order of preference,
	// with the preferred version being the first entry.
	Versions []APIVersionDiscovery `json:"versions"`
}

func (in *APIGroupDiscovery) DeepCopyInto(out *APIGroupDiscovery) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Versions != nil {
		l := make([]APIVersionDiscovery, len(in.Versions))
		for i := range in.Versions {
			in.Versions[i].DeepCopyInto(&l[i])
		}
		out.Versions = l
	}
}

func (in *APIGroupDiscovery) DeepCopy() *APIGroupDiscovery {
	if in == nil {
		return nil
	}
	out := new(APIGroupDiscovery)
	in.DeepCopyInto(out)
	return out
}

func (in *APIGroupDiscovery) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type APIGroupDiscoveryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []APIGroupDiscovery `json:"items"`
}

func (in *APIGroupDiscoveryList) DeepCopyInto(out *APIGroupDiscoveryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]APIGroupDiscovery, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *APIGroupDiscoveryList) DeepCopy() *APIGroupDiscoveryList {
	if in == nil {
		return nil
	}
	out := new(APIGroupDiscoveryList)
	in.DeepCopyInto(out)
	return out
}

func (in *APIGroupDiscoveryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type APIVersionDiscovery struct {
	// version is the name of the version within a group version.
	Version string `json:"version"`
	// resources is a list of APIResourceDiscovery objects for the corresponding group version.
	Resources []APIResourceDiscovery `json:"resources"`
	// freshness marks whether a group version's discovery document is up to date.
	// "Current" indicates the discovery document was recently
	// refreshed. "Stale" indicates the discovery document could not
	// be retrieved and the returned discovery document may be
	// significantly out of date. Clients that require the latest
	// version of the discovery information be retrieved before
	// performing an operation should not use the aggregated document
	Freshness DiscoveryFreshness `json:"freshness,omitempty"`
}

func (in *APIVersionDiscovery) DeepCopyInto(out *APIVersionDiscovery) {
	*out = *in
	if in.Resources != nil {
		l := make([]APIResourceDiscovery, len(in.Resources))
		for i := range in.Resources {
			in.Resources[i].DeepCopyInto(&l[i])
		}
		out.Resources = l
	}
}

func (in *APIVersionDiscovery) DeepCopy() *APIVersionDiscovery {
	if in == nil {
		return nil
	}
	out := new(APIVersionDiscovery)
	in.DeepCopyInto(out)
	return out
}

type APIResourceDiscovery struct {
	// resource is the plural name of the resource.  This is used in the URL path and is the unique identifier
	// for this resource across all versions in the API group.
	// Resources with non-empty groups are located at /apis/<APIGroupDiscovery.objectMeta.name>/<APIVersionDiscovery.version>/<APIResourceDiscovery.Resource>
	// Resources with empty groups are located at /api/v1/<APIResourceDiscovery.Resource>
	Resource string `json:"resource"`
	// responseKind describes the group, version, and kind of the serialization schema for the object type this endpoint typically returns.
	// APIs may return other objects types at their discretion, such as error conditions, requests for alternate representations, or other operation specific behavior.
	// This value will be null if an APIService reports subresources but supports no operations on the parent resource
	ResponseKind *metav1.GroupVersionKind `json:"responseKind,omitempty"`
	// scope indicates the scope of a resource, either Cluster or Namespaced
	Scope ResourceScope `json:"scope"`
	// singularResource is the singular name of the resource.  This allows clients to handle plural and singular opaquely.
	// For many clients the singular form of the resource will be more understandable to users reading messages and should be used when integrating the name of the resource into a sentence.
	// The command line tool kubectl, for example, allows use of the singular resource name in place of plurals.
	// The singular form of a resource should always be an optional element - when in doubt use the canonical resource name.
	SingularResource string `json:"singularResource"`
	// verbs is a list of supported API operation types (this includes
	// but is not limited to get, list, watch, create, update, patch,
	// delete, deletecollection, and proxy).
	Verbs []string `json:"verbs"`
	// shortNames is a list of suggested short names of the resource.
	ShortNames []string `json:"shortNames"`
	// categories is a list of the grouped resources this resource belongs to (e.g. 'all').
	// Clients may use this to simplify acting on multiple resource types at once.
	Categories []string `json:"categories"`
	// subresources is a list of subresources provided by this resource. Subresources are located at /apis/<APIGroupDiscovery.objectMeta.name>/<APIVersionDiscovery.version>/<APIResourceDiscovery.Resource>/name-of-instance/<APIResourceDiscovery.subresources[i].subresource>
	Subresources []APISubresourceDiscovery `json:"subresources"`
}

func (in *APIResourceDiscovery) DeepCopyInto(out *APIResourceDiscovery) {
	*out = *in
	if in.ResponseKind != nil {
		in, out := &in.ResponseKind, &out.ResponseKind
		*out = new(metav1.GroupVersionKind)
		(*in).DeepCopyInto(*out)
	}
	if in.Verbs != nil {
		t := make([]string, len(in.Verbs))
		copy(t, in.Verbs)
		out.Verbs = t
	}
	if in.ShortNames != nil {
		t := make([]string, len(in.ShortNames))
		copy(t, in.ShortNames)
		out.ShortNames = t
	}
	if in.Categories != nil {
		t := make([]string, len(in.Categories))
		copy(t, in.Categories)
		out.Categories = t
	}
	if in.Subresources != nil {
		l := make([]APISubresourceDiscovery, len(in.Subresources))
		for i := range in.Subresources {
			in.Subresources[i].DeepCopyInto(&l[i])
		}
		out.Subresources = l
	}
}

func (in *APIResourceDiscovery) DeepCopy() *APIResourceDiscovery {
	if in == nil {
		return nil
	}
	out := new(APIResourceDiscovery)
	in.DeepCopyInto(out)
	return out
}

type APISubresourceDiscovery struct {
	// subresource is the name of the subresource.  This is used in the URL path and is the unique identifier
	// for this resource across all versions.
	Subresource string `json:"subresource"`
	// responseKind describes the group, version, and kind of the serialization schema for the object type this endpoint typically returns.
	// Some subresources do not return normal resources, these will have null return types.
	ResponseKind *metav1.GroupVersionKind `json:"responseKind,omitempty"`
	// acceptedTypes describes the kinds that this endpoint accepts.
	// Subresources may accept the standard content types or define
	// custom negotiation schemes. The list may not be exhaustive for
	// all operations.
	AcceptedTypes []metav1.GroupVersionKind `json:"acceptedTypes"`
	// verbs is a list of supported API operation types (this includes
	// but is not limited to get, list, watch, create, update, patch,
	// delete, deletecollection, and proxy). Subresources may define
	// custom verbs outside the standard Kubernetes verb set. Clients
	// should expect the behavior of standard verbs to align with
	// Kubernetes interaction conventions.
	Verbs []string `json:"verbs"`
}

func (in *APISubresourceDiscovery) DeepCopyInto(out *APISubresourceDiscovery) {
	*out = *in
	if in.ResponseKind != nil {
		in, out := &in.ResponseKind, &out.ResponseKind
		*out = new(metav1.GroupVersionKind)
		(*in).DeepCopyInto(*out)
	}
	if in.AcceptedTypes != nil {
		l := make([]metav1.GroupVersionKind, len(in.AcceptedTypes))
		for i := range in.AcceptedTypes {
			in.AcceptedTypes[i].DeepCopyInto(&l[i])
		}
		out.AcceptedTypes = l
	}
	if in.Verbs != nil {
		t := make([]string, len(in.Verbs))
		copy(t, in.Verbs)
		out.Verbs = t
	}
}

func (in *APISubresourceDiscovery) DeepCopy() *APISubresourceDiscovery {
	if in == nil {
		return nil
	}
	out := new(APISubresourceDiscovery)
	in.DeepCopyInto(out)
	return out
}
