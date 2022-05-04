package test

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespace string = "monitoring"
var kubeOptions = k8s.NewKubectlOptions("", "", namespace)

// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
var tlsConfig tls.Config = tls.Config{}

// Returns pod name
func awaitPods(t *testing.T, kubeOptions, filter string) string {
	var podName string

	pods := k8s.ListPods(t, kubeOptions, v1.ListOptions{FieldSelector: "status.phase=Running"})

	for _, pod := range pods {

		// Await all deployed pods to be available and ready
		k8s.WaitUntilPodAvailable(t, kubeOptions, pod.Name, 60, 1*time.Second)

		if strings.Contains(pod.Name, filter) {
			podName := pod.Name
			return podName
		}
	}

	return podName
}

// Returns service name
func awaitServices(t *testing.T, kubeOptions kubeOptions, filter string) string {
	var serviceName string

	services := k8s.ListServices(t, kubeOptions, v1.ListOptions{FieldSelector: "metadata.namespace=monitoring"})

	for _, service := range services {

		// Await all deployed services to be available and ready
		k8s.WaitUntilServiceAvailable(t, kubeOptions, service.Name, 60, 1*time.Second)

		if strings.Contains(service.Name, filter) {
			serviceName := service.Name
			return serviceName
		}
	}
	return serviceName
}

func verifyPrometheusWelcomePage(statusCode int, body string) bool {
	if statusCode != 200 {
		return false
	}

	return strings.Contains(body, `<title>Prometheus Time Series Collection and Processing Server</title>`)
}

func verifyGrafanaWelcomePage(statusCode int, body string) bool {
	if statusCode != 200 {
		return false
	}

	return strings.Contains(body, `<title>Grafana</title>`)
}

func TestingPrometheusDeployment(t *testing.T) {

	var filter string = "server"

	prometheusPod := awaitPods(t, kubeOptions, filter)

	awaitServices(t, kubeOptions, filter)

	tunnel := k8s.NewTunnel(kubeOptions, k8s.ResourceTypePod, prometheusPod, 0, 9090)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

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

func TestingGrafanaDeployment(t *testing.T) {

	var filter string = "grafana"

	grafanaPod := awaitPods(t, kubeOptions, filter)

	awaitServices(t, kubeOptions, filter)

	tunnel := k8s.NewTunnel(kubeOptions, k8s.ResourceTypePod, grafanaPod, 0, 3000)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// Try to access the grafana service on the local port, retrying until we get a good response for up to 5 minutes
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://%s", tunnel.Endpoint()),
		&tlsConfig,
		60,
		5*time.Second,
		verifyGrafanaWelcomePage,
	)

}

func TestInfrastructure(t *testing.T) {

	terraformOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"grafana_admin_password": strings.ToLower(random.UniqueId()),
		},
	})

	defer terraform.Destroy(t, terraformOpts)

	terraform.InitAndApply(t, terraformOpts)

	kubeconfigFilename := terraform.Output(t, terraformOpts, "kubeconfig_filename")

	assert.Contains(t, kubeconfigFilename, "config")

	TestingPrometheusDeployment(t)

	TestingGrafanaDeployment(t)

}
