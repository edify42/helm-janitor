package scan

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	janitorconfig "github.com/lendi-au/helm-janitor/internal/config"
	client "github.com/lendi-au/helm-janitor/internal/eks"
	internalhelm "github.com/lendi-au/helm-janitor/internal/helm"
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
	Getekscluster(aws.Config, client.Generator) client.EKSCluster
	Config() janitorconfig.EnvConfig
	Getreleases(client.EKSCluster, *action.Configuration, internalhelm.HelmList) []*release.Release
	Deleterelease(client.EKSCluster, *action.Configuration, *release.Release, internalhelm.HelmDelete) error
	Makeekscfg() client.Generator // Experimental. Using this to mock...
}

type ScanClient struct {
	Selector          string
	Namespace         string
	AllNamespaces     bool
	IncludeNamespaces string
	ExcludeNamespaces string
	Env               janitorconfig.EnvConfig
	Dryrun            bool
	Context           context.Context
}

func NewScanClient() *ScanClient {
	return &ScanClient{}
}

// Loose - experimental...
func (sc *ScanClient) Makeekscfg() client.Generator {
	if sc.Context != nil {
		return &client.GeneratorType{
			Context: sc.Context,
		}
	}
	return &client.GeneratorType{
		Context: context.TODO(),
	}
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
		sc.Context,
		config.WithRegion(sc.Env.Region),
	)
	if err != nil {
		// handle error :(
		log.Panic("aws config management issue...")
	}

	if os.Getenv("ROLE_ARN") != "" {
		// reference: https://stackoverflow.com/questions/67605787/how-do-i-use-aws-sts-assume-role-with-mfa-using-aws-sdk-go-v2
		creds := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg), os.Getenv("ROLE_ARN"), func(o *stscreds.AssumeRoleOptions) {})
		cfg.Credentials = aws.NewCredentialsCache(creds)
	}

	return cfg
}

// Getekscluster - Return the cluster, endpoints and auth token!
func (sc *ScanClient) Getekscluster(c aws.Config, g client.Generator) client.EKSCluster {
	a := client.AwsConfig{J: sc.Env}
	cluster := a.Init(c, g)
	return cluster
}

// Getreleases will return the array of helm releases
func (sc *ScanClient) Getreleases(c client.EKSCluster, a *action.Configuration, list internalhelm.HelmList) []*release.Release {
	releaseNamespace := sc.Namespace
	if releaseNamespace != "" {
		log.Debugf("Getting releases from namespace: %s", releaseNamespace)
	} else {
		log.Debugf("No namespace defined - AllNamespaces is %v", sc.AllNamespaces)
	}

	settings := cli.New()
	settings.KubeAPIServer = c.Endpoint
	settings.KubeToken = c.Token
	settings.KubeCaFile = c.CAFile
	if err := a.Init(settings.RESTClientGetter(), releaseNamespace, "secrets", log.Infof); err != nil {
		panic(err)
	}

	iCli := list.ActionNewList(a) // super confusing? same obj in memory so...
	iCli.Selector = sc.Env.JanitorLabel
	iCli.AllNamespaces = sc.AllNamespaces

	if sc.Selector != "" {
		iCli.Selector = fmt.Sprintf("%s,%s", sc.Env.JanitorLabel, sc.Selector)
	}
	rel, err := list.RunCommand()
	if err != nil {
		log.Panic(err)
	}
	releaseList := NameList(rel)
	log.Debugf("Got a list of releases: %v", releaseList)

	return rel
}

// Deleterelease will try and delete a release
func (sc *ScanClient) Deleterelease(eks client.EKSCluster, a *action.Configuration, rel *release.Release, del internalhelm.HelmDelete) error {
	os.Setenv("HELM_NAMESPACE", rel.Namespace) // holy hack batman. why don't they expose this at the API level?
	settings := cli.New()
	settings.KubeAPIServer = eks.Endpoint
	settings.KubeToken = eks.Token
	settings.KubeCaFile = eks.CAFile
	if err := a.Init(settings.RESTClientGetter(), rel.Namespace, "secrets", log.Infof); err != nil {
		log.Fatal(err)
	}
	run := del.ActionNewUninstall(a)
	if sc.Dryrun {
		log.Debug("dry run mode enabled")
		run.DryRun = sc.Dryrun
	}
	out, err := del.RunCommand(rel.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("deleted: %s", out.Info)
	return nil
}
