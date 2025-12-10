package k8sclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"go.f110.dev/kubeproto/go/apis/corev1"
)

func (c *CoreV1) PortForward(ctx context.Context, pod *corev1.Pod, port int) (*portforward.PortForwarder, uint16, error) {
	req := c.backend.RESTClient().Post().Resource("pods").Namespace(pod.Namespace).Name(pod.Name).SubResource("portforward")
	transport, upgrader, err := spdy.RoundTripperFor(c.config)
	if err != nil {
		return nil, 0, err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())

	readyCh := make(chan struct{})
	pf, err := portforward.New(dialer, []string{fmt.Sprintf(":%d", port)}, ctx.Done(), readyCh, nil, nil)
	if err != nil {
		return nil, 0, err
	}
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)

		err := pf.ForwardPorts()
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case <-readyCh:
	case err := <-errCh:
		return nil, 0, err
	case <-time.After(5 * time.Second):
		return nil, 0, errors.New("timed out")
	}

	ports, err := pf.GetPorts()
	if err != nil {
		return nil, 0, err
	}

	return pf, ports[0].Local, nil
}
