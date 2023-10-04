package k8sclient

import (
	"context"
	"encoding/json"
	"mime"
	"net/url"
	"path"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"go.f110.dev/kubeproto/go/apis/apidiscoveryv2beta1"
	"go.f110.dev/kubeproto/go/apis/metav1"
)

const (
	discoveryAcceptHeader = "application/json;g=apidiscovery.k8s.io;v=v2beta1;as=APIGroupDiscoveryList,application/json"
)

type DiscoveryClient struct {
	client         *rest.RESTClient
	maxConcurrency int
}

func NewDiscoveryClient(cfg *rest.Config) (*DiscoveryClient, error) {
	codec := runtime.NoopEncoder{Decoder: scheme.Codecs.UniversalDecoder()}
	cfg.NegotiatedSerializer = serializer.NegotiatedSerializerWrapper(runtime.SerializerInfo{Serializer: codec})
	c, err := rest.UnversionedRESTClientFor(cfg)
	if err != nil {
		return nil, err
	}
	return &DiscoveryClient{client: c, maxConcurrency: 10}, nil
}

func (d *DiscoveryClient) APIResourceLists(ctx context.Context) (*metav1.APIGroupList, map[schema.GroupVersion]*metav1.APIResourceList, error) {
	var contentType string
	body, err := d.client.Get().
		AbsPath("/apis").
		SetHeader("Accept", discoveryAcceptHeader).
		Do(ctx).
		ContentType(&contentType).
		Raw()
	if err != nil {
		return nil, nil, err
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, nil, err
	}
	var groupList metav1.APIGroupList
	resources := make(map[schema.GroupVersion]*metav1.APIResourceList)
	if mediaType == "application/json" && params["g"] == "apidiscovery.k8s.io" && params["v"] == "v2beta1" && params["as"] == "APIGroupDiscoveryList" {
		var discoveryList apidiscoveryv2beta1.APIGroupDiscoveryList
		if err := json.Unmarshal(body, &discoveryList); err != nil {
			return nil, nil, err
		}

		for _, v := range discoveryList.Items {
			group, resourceList := d.toAPIGroup(v)

			groupList.Groups = append(groupList.Groups, *group)
			for gv, list := range resourceList {
				resources[gv] = list
			}
		}
	} else {
		if err := json.Unmarshal(body, &groupList); err != nil {
			return nil, nil, err
		}

		sem := make(chan struct{}, d.maxConcurrency)
		var errs []error
		var wg sync.WaitGroup
		for _, group := range groupList.Groups {
			for _, version := range group.Versions {
				gv := schema.GroupVersion{Group: group.Name, Version: version.Version}
				wg.Add(1)
				go func(gv schema.GroupVersion) {
					sem <- struct{}{}
					defer func() {
						<-sem
						wg.Done()
					}()

					resourceList, err := d.fetchAPIResourceList(ctx, gv)
					if err != nil {
						errs = append(errs, err)
						return
					}
					resources[gv] = resourceList
				}(gv)
			}
		}
		wg.Wait()
		if len(errs) > 0 {
			return nil, nil, errs[0]
		}
	}
	if len(resources) > 0 {
		return &groupList, resources, nil
	}

	return &groupList, resources, nil
}

func (d *DiscoveryClient) SetMaxConcurrency(c int) {
	d.maxConcurrency = c
}

func (d *DiscoveryClient) fetchAPIResourceList(ctx context.Context, groupVersion schema.GroupVersion) (*metav1.APIResourceList, error) {
	u := url.URL{Path: path.Join("apis", groupVersion.String())}

	var resourceList metav1.APIResourceList
	err := d.client.Get().
		AbsPath(u.String()).
		Do(ctx).
		Into(&resourceList)
	if err != nil {
		return nil, err
	}
	return &resourceList, nil
}

func (*DiscoveryClient) toAPIGroup(in apidiscoveryv2beta1.APIGroupDiscovery) (*metav1.APIGroup, map[schema.GroupVersion]*metav1.APIResourceList) {
	var emptyKind = metav1.GroupVersionKind{}

	apiGroup := &metav1.APIGroup{}
	apiResources := make(map[schema.GroupVersion]*metav1.APIResourceList)
	for _, v := range in.Versions {
		gv := schema.GroupVersion{Group: in.Name, Version: v.Version}
		version := metav1.GroupVersionForDiscovery{GroupVersion: gv.String(), Version: v.Version}
		apiGroup.Versions = append(apiGroup.Versions, version)

		var resources *metav1.APIResourceList
		for _, r := range v.Resources {
			resource := metav1.APIResource{
				Name:         r.Resource,
				SingularName: r.SingularResource,
				Verbs:        r.Verbs,
				ShortNames:   r.ShortNames,
				Categories:   r.Categories,
				Namespaced:   r.Scope == apidiscoveryv2beta1.ResourceScopeNamespaced,
			}
			if r.ResponseKind != nil && *r.ResponseKind != emptyKind {
				resource.Group = r.ResponseKind.Group
				resource.Version = r.ResponseKind.Version
				resource.Kind = r.ResponseKind.Kind
			}
			resources.APIResources = append(resources.APIResources, resource)
		}
		apiResources[gv] = resources
	}

	return apiGroup, apiResources
}
