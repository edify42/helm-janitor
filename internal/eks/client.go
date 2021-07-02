package eks

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

// New client to generate
func New(cluster *types.Cluster) (*kubernetes.Clientset, error) {
	ca, err := base64.StdEncoding.DecodeString(aws.ToString(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, err
	}
	tok, err := GenToken(cluster.Name)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(
		&rest.Config{
			Host:        aws.ToString(cluster.Endpoint),
			BearerToken: tok.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// GenToken will get the EKS cluster oauth2 token.
// Consider refresh flow instead and make token private
// and accessible via function call.
func GenToken(cluster *string) (token.Token, error) {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return token.Token{}, err
	}
	opts := &token.GetTokenOptions{
		ClusterID: aws.ToString(cluster),
	}
	return gen.GetWithOptions(opts)
}
