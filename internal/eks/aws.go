package eks

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	log "github.com/sirupsen/logrus"
)

type AwsConfig struct {
	J     janitorconfig.EnvConfig
	Token string
}

// aws specific EKS logic.
type EKSCluster struct {
	Name     string
	Endpoint string
	Token    string
	CAFile   string
}

func (a *AwsConfig) Init(cfg aws.Config, g Generator) EKSCluster {
	// return value
	cluster := EKSCluster{}
	eksClient := eks.NewFromConfig(cfg)

	// Get cluster endpoint
	result, err := g.DescribeCluster(eksClient, a.J.Cluster)
	if err != nil {
		log.Fatalf("Error calling DescribeCluster: %v", err)
	}
	cluster.Name = *result.Cluster.Name
	cluster.Endpoint = *result.Cluster.Endpoint

	// Get cluster token
	newToken, err := g.GenToken(result.Cluster.Name)
	if err != nil {
		panic("Shieeet")
	}

	// Verify connectivity to the cluster (list nodes or something)
	g.TestCluster(result.Cluster, newToken)

	// first generate a temp file to write the cluster CA to
	// handle file close in parent function
	file := g.WriteCA(*a, result)

	cluster.CAFile = file

	cluster.Token = newToken.Token
	a.Token = newToken.Token
	return cluster
}
