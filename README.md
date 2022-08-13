# awesome-eyes
awesomeeyes is deploying grafana and prometheus helm releases on eks cluster



## Setup

The following is required:

- Terraform v1.2.2
- GoLang v1.18.3

These enviroment variables are required:

- GRAFANA_PASSWORD
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY

## Deploy

```
terraform init
TF_VAR_grafana_admin_password=$GRAFANA_PASSWORD terraform apply --auto-approve
```

Remote backend is configured, if you want to use it, you have to create the following on your aws account on region eu-west-3:

- S3 bucket: "awesome-eyes-terrafrom-state"
- DynamoDB table: "awesome-eyes-locks"

Otherwise, you can comment this part and have terraform state locally

## Test

```
cd test && go test -v -timeout 3000s infra_test.go
```
Grafana password is initialized within test
New kubernetes clientset is being generated within test to authenticate cluster connection

## Destroy

```
TF_VAR_grafana_admin_password=$GRAFANA_PASSWORD terraform destroy --auto-approve
```

## Technical Description

Terraform script will do the following:

- Create eks cluster v1.22 with two nodes deployed on two subnets
- Deploy one replica of grafana helm release v6.24.1
- Deploy one replica of prometheus helm release v15.5.3

## Github Workflow

On every push for any branch will trigger ```deploy-terraform``` it will test, deploy then destroy (Will be locked only for main)
On main and develop pull requests will trigger ```validate-terraform``` it will format, validate and plan

Environment variables are hooked up as secret like the following:
```
env:
  PASSWD: ${{ secrets.GRAFANA_PASSWORD }}
  AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
  AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```