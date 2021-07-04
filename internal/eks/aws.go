package eks

import (
	"context"
	"encoding/base64"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AwsConfig struct {
	janitorconfig.EnvConfig
}

// aws specific EKS logic.
type EKSCluster struct {
	Name     string
	Endpoint string
	Token    string
	CAFile   string
}

func (a *AwsConfig) Init(cfg aws.Config) EKSCluster {
	// return value
	cluster := EKSCluster{}
	eksClient := eks.NewFromConfig(cfg)

	result, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{Name: &a.Cluster})
	if err != nil {
		log.Fatalf("Error calling DescribeCluster: %v", err)
	}
	cluster.Name = *result.Cluster.Name
	cluster.Endpoint = *result.Cluster.Endpoint

	clientset, err := New(result.Cluster)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting EKS nodes: %v", err)
	}
	log.Debugf("There are %d nodes associated with cluster %s", len(nodes.Items), a.Cluster)

	// first generate a temp file to write the cluster CA to
	// handle file close in parent function
	file, err := ioutil.TempFile(a.TmpFileLocation, a.TmpFilePrefix)
	if err != nil {
		log.Fatal(err)
	}
	decoded, _ := base64.StdEncoding.DecodeString(*result.Cluster.CertificateAuthority.Data)
	file.Write([]byte(decoded))
	if err = file.Close(); err != nil {
		log.Fatal(err)
	}

	cluster.CAFile = file.Name()

	// configure the helm api to re-use the EKS auth above.
	newToken, err := GenToken(result.Cluster.Name)
	if err != nil {
		panic("Shieeet")
	}

	cluster.Token = newToken.Token
	return cluster
}
