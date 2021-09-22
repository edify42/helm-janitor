module github.com/lendi-au/helm-janitor

go 1.16

require (
	github.com/aws/aws-lambda-go v1.26.0
	github.com/aws/aws-sdk-go-v2 v1.9.1
	github.com/aws/aws-sdk-go-v2/config v1.8.2
	github.com/aws/aws-sdk-go-v2/credentials v1.4.2
	github.com/aws/aws-sdk-go-v2/service/eks v1.10.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.7.1
	github.com/fnproject/fdk-go v0.0.8
	github.com/mitchellh/mapstructure v1.4.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/tencentyun/scf-go-lib v0.0.0-20200624065115-ba679e2ec9c9
	helm.sh/helm/v3 v3.7.0
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/aws-iam-authenticator v0.5.3
)
