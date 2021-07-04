package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	"github.com/edify42/helm-janitor/internal/eks"
	client "github.com/edify42/helm-janitor/internal/eks"
	log "github.com/sirupsen/logrus"
)

// TODO: rename this file...
// Supports the scan.go file by creating the interface that can be
// mocked out for testing.

type InputRun interface {
	Init()
	Makeawscfg() aws.Config
	Getekscluster(aws.Config) client.EKSCluster
	Config() janitorconfig.EnvConfig
}

type ScanClient struct {
	Selector string
	Env      janitorconfig.EnvConfig
}

func NewScanClient() *ScanClient {
	return &ScanClient{}
}

// Init - initialise!
func (sc *ScanClient) Init() {
	test := janitorconfig.EnvConfig{}
	test.Init() // get the default values...again.
	sc.Env = test
	log.Infof("ScanClient initialised with values %v", sc)
}

// Config - return it!
func (sc *ScanClient) Config() janitorconfig.EnvConfig {
	return janitorconfig.EnvConfig(sc.Env)
}

// Makeawscfg - creates the cfg object
func (sc *ScanClient) Makeawscfg() aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(sc.Env.Region),
	)
	if err != nil {
		// handle error :(
		log.Panic("aws config management issue...")
	}
	return cfg
}

// Getekscluster - Return the cluster, endpoints and auth token!
func (sc *ScanClient) Getekscluster(c aws.Config) eks.EKSCluster {
	a := eks.AwsConfig{sc.Env}
	cluster := a.Init(c)
	return cluster
}
