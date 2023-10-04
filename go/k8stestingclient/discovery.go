package k8stestingclient

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"go.f110.dev/kubeproto/go/apis/metav1"
	"go.f110.dev/kubeproto/go/k8sclient"
)

type DiscoveryClient struct {
	k8sclient.DiscoveryClient

	Resources []*metav1.APIResourceList
}

func (d *DiscoveryClient) APIResourceLists(_ context.Context) (*metav1.APIGroupList, map[schema.GroupVersion]*metav1.APIResourceList, error) {
	var groups metav1.APIGroupList
	resources := make(map[schema.GroupVersion]*metav1.APIResourceList)
	for _, v := range d.Resources {
		gv := v.GroupVersionKind().GroupVersion()
		resources[gv] = v

		var g metav1.APIGroup
		g.Versions = []metav1.GroupVersionForDiscovery{{GroupVersion: gv.String(), Version: gv.Version}}
		groups.Groups = append(groups.Groups, g)
	}

	return &groups, resources, nil
}
