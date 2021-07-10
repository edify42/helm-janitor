module github.com/edify42/helm-janitor

go 1.16

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go-v2 v1.7.0
	github.com/aws/aws-sdk-go-v2/config v1.4.0
	github.com/aws/aws-sdk-go-v2/service/eks v1.7.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.0
	helm.sh/helm/v3 v3.6.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/aws-iam-authenticator v0.5.3
)
