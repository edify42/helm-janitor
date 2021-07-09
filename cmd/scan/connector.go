package scan

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	client "github.com/edify42/helm-janitor/internal/eks"
	internalhelm "github.com/edify42/helm-janitor/internal/helm"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

// TODO: rename this file...
// Supports the scan.go file by creating the interface that can be
// mocked out for testing.

type InputRun interface {
	Init()
	Makeawscfg() aws.Config
	Getekscluster(aws.Config) client.EKSCluster
	Config() janitorconfig.EnvConfig
	Getreleases(client.EKSCluster, *action.Configuration, internalhelm.HelmList) []*release.Release
	Deleterelease(*action.Configuration, *release.Release) error
}

type ScanClient struct {
	Selector          string
	Namespace         string
	AllNamespaces     bool
	IncludeNamespaces string
	ExcludeNamespaces string
	Env               janitorconfig.EnvConfig
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
func (sc *ScanClient) Getekscluster(c aws.Config) client.EKSCluster {
	a := client.AwsConfig{sc.Env}
	cluster := a.Init(c)
	return cluster
}

// Getreleases will return the array of helm releases
func (sc *ScanClient) Getreleases(c client.EKSCluster, a *action.Configuration, list internalhelm.HelmList) []*release.Release {
	releaseNamespace := sc.Namespace // TODO: use the parameters to figure out the namespaces to search.
	log.Debugf("Getting releases from namespace: %s", releaseNamespace)
	settings := cli.New()
	settings.KubeAPIServer = c.Endpoint
	settings.KubeToken = c.Token
	settings.KubeCaFile = c.CAFile
	if err := a.Init(settings.RESTClientGetter(), releaseNamespace, os.Getenv("HELM_DRIVER"), log.Infof); err != nil {
		panic(err)
	}

	iCli := list.NewList(a)
	iCli.Selector = sc.Env.JanitorLabel
	rel, err := list.RunCommand()
	if err != nil {
		log.Panic(err)
	}
	releaseList := NameList(rel)
	log.Debugf("Got a list of releases: %v", releaseList)

	return rel
}

// Deleterelease will try and delete a release
func (sc *ScanClient) Deleterelease(a *action.Configuration, rel *release.Release) error {
	settings := cli.New()
	if err := a.Init(settings.RESTClientGetter(), rel.Namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v)
	}); err != nil {
		log.Fatal(err)
	}
	run := action.NewUninstall(a)
	run.DryRun = false
	out, err := run.Run(rel.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("deleted: %s", out.Info)
	return nil
}
