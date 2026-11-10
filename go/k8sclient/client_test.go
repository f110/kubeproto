package k8sclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/client-go/rest"

	"go.f110.dev/kubeproto/go/apis/appsv1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	"go.f110.dev/kubeproto/go/internal/assertion"
)

func TestApply(t *testing.T) {
	var req *http.Request
	var body []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		body, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"test","namespace":"default"}}`))
	}))
	t.Cleanup(srv.Close)
	s, err := NewSet(&rest.Config{Host: srv.URL})
	assertion.MustNoError(t, err)

	d, err := s.AppsV1.ApplyDeployment(
		t.Context(),
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: metav1.NamespaceDefault}},
		metav1.ApplyOptions{FieldManager: "test", Force: true},
	)
	assertion.MustNoError(t, err)
	assertion.Equal(t, "test", d.Name)

	assertion.Equal(t, http.MethodPatch, req.Method)
	assertion.Equal(t, "/apis/apps/v1/namespaces/default/deployments/test", req.URL.Path)
	assertion.Equal(t, "test", req.URL.Query().Get("fieldManager"))
	assertion.Equal(t, "true", req.URL.Query().Get("force"))
	assertion.Equal(t, "application/apply-patch+yaml", req.Header.Get("Content-Type"))
	var applied map[string]any
	assertion.MustNoError(t, json.Unmarshal(body, &applied))
	assertion.Equal(t, "apps/v1", applied["apiVersion"])
	assertion.Equal(t, "Deployment", applied["kind"])
}
