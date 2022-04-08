package test

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deploymentTest struct {
	namespace   string
	kubeOptions *k8s.KubectlOptions
}

func awaitPrometheusPods(t *testing.T, kubeOptions *k8s.KubectlOptions) string {

	var promPod string

	pods := k8s.ListPods(t, kubeOptions, v1.ListOptions{LabelSelector: "app=prometheus"})

	for _, pod := range pods {

		k8s.WaitUntilPodAvailable(t, kubeOptions, pod.Name, 60, 1*time.Second)

		// Should get prometheus-server pod name
		if strings.Contains(pod.Name, "server") {
			promPod := pod.Name
			return promPod
		}
	}
	return promPod
}

func awaitPrometheusServices(t *testing.T, kubeOptions *k8s.KubectlOptions) {

	services := k8s.ListServices(t, kubeOptions, v1.ListOptions{LabelSelector: "app=prometheus"})

	for _, service := range services {

		k8s.WaitUntilServiceAvailable(t, kubeOptions, service.Name, 60, 1*time.Second)

	}
}

func verifyPrometheusWelcomePage(statusCode int, body string) bool {
	if statusCode != 200 {
		return false
	}

	return strings.Contains(body, `<title>Prometheus Time Series Collection and Processing Server</title>`)
}

func TestPrometheusDeployment(t *testing.T) {
	t.Parallel()

	namespace := "monitoring"

	kubeOptions := k8s.NewKubectlOptions("", "", namespace)

	// Wait for pods to be ready and store prometheus server pod name
	promPodName := awaitPrometheusPods(t, kubeOptions)

	// Wait for services to be ready
	awaitPrometheusServices(t, kubeOptions)

	fmt.Println("::::CEREBRAL_DEBUG:::promPodName::", promPodName)

	tunnel := k8s.NewTunnel(kubeOptions, k8s.ResourceTypePod, promPodName, 0, 9090)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Try to access the prometheus service on the local port, retrying until we get a good response for up to 5 minutes
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://%s", tunnel.Endpoint()),
		&tlsConfig,
		60,
		5*time.Second,
		verifyPrometheusWelcomePage,
	)

}
