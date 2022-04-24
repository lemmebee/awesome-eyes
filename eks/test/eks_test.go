package test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEKS(t *testing.T) {
	t.Parallel()

	terraformOpts := &terraform.Options{
		TerraformDir: "../",
	}

	defer terraform.Destroy(t, terraformOpts)

	terraform.InitAndApply(t, terraformOpts)

	kubeconfigFilename := terraform.Output(t, terraformOpts, "kubeconfig_filename")

	assert.Contains(t, kubeconfigFilename, "./kubeconfig_awesome-eyes-eks-")
}
