package test

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

var namespace string = "awesomeeyes"
var kubeOptions *k8s.KubectlOptions = k8s.NewKubectlOptions("", "", namespace)

// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
var tlsConfig tls.Config = tls.Config{}

// https://stackoverflow.com/questions/60547409/unable-to-obtain-kubeconfig-of-an-aws-eks-cluster-in-go-code/60573982#60573982
func newClientset(cluster *eks.Cluster) (*kubernetes.Clientset, error) {
	const (
		RecommendedHomeDir  = ".kube"
		RecommendedFileName = "config"
	)

	var RecommendedHomeFile = path.Join(homedir.HomeDir(), RecommendedHomeDir, RecommendedFileName)

	// log.Printf("%+v", cluster)

	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}

	opts := &token.GetTokenOptions{
		ClusterID: aws.StringValue(cluster.Name),
	}

	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, err
	}

	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(
		&rest.Config{
			Host:        aws.StringValue(cluster.Endpoint),
			BearerToken: tok.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters[*cluster.Name] = &clientcmdapi.Cluster{
		Server:                   *cluster.Endpoint,
		CertificateAuthorityData: ca,
	}

	contexts := make(map[string]*clientcmdapi.Context)
	contexts[*cluster.Arn] = &clientcmdapi.Context{
		Cluster:   *cluster.Name,
		Namespace: namespace,
		AuthInfo:  namespace,
	}

	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos[namespace] = &clientcmdapi.AuthInfo{
		Token: tok.Token,
	}

	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: *cluster.Arn,
		AuthInfos:      authinfos,
	}

	clientcmd.WriteToFile(clientConfig, RecommendedHomeFile)

	return clientset, nil
}

// Returns running pods
func getRunningPods(clientset *kubernetes.Clientset, result *eks.DescribeClusterOutput) *corev1.PodList {

	// https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go#L64
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{FieldSelector: "status.phase=Running"})
	if err != nil {
		log.Fatalf("Error getting EKS pods: %v", err)
	}
	log.Printf("There are %d pods associated with cluster: %s", len(pods.Items), *result.Cluster.Name)

	return pods
}

// Wait services
func waitServices(t *testing.T, clientset *kubernetes.Clientset, result *eks.DescribeClusterOutput) {

	svcs, err := clientset.CoreV1().Services(namespace).List(context.TODO(), v1.ListOptions{FieldSelector: "metadata.namespace=" + namespace})
	if err != nil {
		log.Fatalf("Error getting EKS services: %v", err)
	}
	log.Printf("There are %d services associated with cluster: %s", len(svcs.Items), *result.Cluster.Name)

	for i := range svcs.Items {

		// Wait all deployed services to be available and ready
		k8s.WaitUntilServiceAvailable(t, kubeOptions, svcs.Items[i].Name, 60, 1*time.Second)

	}

}

func verifyPrometheus(statusCode int, body string) bool {
	if statusCode != 200 {
		return false
	}

	return strings.Contains(body, "Prometheus Time Series Collection and Processing Server")
}

func ValidatePrometheusDeployment(t *testing.T, clientset *kubernetes.Clientset, result *eks.DescribeClusterOutput) {

	pods := getRunningPods(clientset, result)

	for i := range pods.Items {

		// Wait all deployed pods to be available and ready
		k8s.WaitUntilPodAvailable(t, kubeOptions, pods.Items[i].Name, 60, 1*time.Second)

		var filter string = "server"

		if strings.Contains(pods.Items[i].Name, filter) {

			log.Printf("Prometheus server: %v", pods.Items[i].Name)

			tunnel := k8s.NewTunnel(kubeOptions, k8s.ResourceTypePod, pods.Items[i].Name, 0, 9090)
			defer tunnel.Close()
			tunnel.ForwardPort(t)

			// Try to access prometheus service on the local port, retrying until we get a good response for up to 5 minutes
			http_helper.HttpGetWithRetryWithCustomValidation(
				t,
				fmt.Sprintf("http://%s", tunnel.Endpoint()),
				&tlsConfig,
				60,
				5*time.Second,
				verifyPrometheus,
			)
		}
	}

}

func verifyGrafana(statusCode int, body string) bool {
	if statusCode != 200 {
		return false
	}

	return strings.Contains(body, "Grafana")
}

func ValidateGrafanaDeployment(t *testing.T, clientset *kubernetes.Clientset, result *eks.DescribeClusterOutput) {

	pods := getRunningPods(clientset, result)

	for i := range pods.Items {

		// Wait all deployed pods to be available and ready
		k8s.WaitUntilPodAvailable(t, kubeOptions, pods.Items[i].Name, 60, 1*time.Second)

		var filter string = "grafana"

		if strings.Contains(pods.Items[i].Name, filter) {

			log.Printf("Grafana server: %v", pods.Items[i].Name)

			tunnel := k8s.NewTunnel(kubeOptions, k8s.ResourceTypePod, pods.Items[i].Name, 0, 3000)
			defer tunnel.Close()
			tunnel.ForwardPort(t)

			// Try to access grafana service on the local port, retrying until we get a good response for up to 5 minutes
			http_helper.HttpGetWithRetryWithCustomValidation(
				t,
				fmt.Sprintf("http://%s", tunnel.Endpoint()),
				&tlsConfig,
				60,
				5*time.Second,
				verifyGrafana,
			)
		}
	}

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

	clusterName := terraform.Output(t, terraformOpts, "cluster_name")
	region := terraform.Output(t, terraformOpts, "region")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	eksSvc := eks.New(sess)

	input := &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}

	result, err := eksSvc.DescribeCluster(input)
	assert.NoError(t, err)

	clientset, err := newClientset(result.Cluster)

	assert.NoError(t, err)

	waitServices(t, clientset, result)

	ValidatePrometheusDeployment(t, clientset, result)

	ValidateGrafanaDeployment(t, clientset, result)

}
