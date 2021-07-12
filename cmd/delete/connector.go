package delete

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	client "github.com/edify42/helm-janitor/internal/eks"
	internalhelm "github.com/edify42/helm-janitor/internal/helm"
	log "github.com/sirupsen/logrus"
)

// Client is the data object which contains the item to delete
type Client struct {
	Dryrun    bool
	Release   string
	Namespace string
	Env       janitorconfig.EnvConfig
}

// NewClient will return the Client struct
func NewClient() *Client {
	return &Client{}
}

// InputRun is our interface which defines the main delete methods
type InputRun interface {
	Config() *Client
	Init()
	Makeawscfg() aws.Config
	Getekscluster(aws.Config, client.Generator) client.EKSCluster
	Deleterelease(client.EKSCluster, *action.Configuration, *release.Release, internalhelm.HelmDelete) error
	Makeekscfg() client.Generator // Experimental. Using this to mock...
}

// Config - return it!
func (c *Client) Config() *Client {
	return c
}

// Init it!
func (d *Client) Init() {
	test := janitorconfig.EnvConfig{}
	test.Init() // get the default values...again.
	d.Env = test
	log.Infof("Delete client initialised with values %v", d)
}

// Loose - experimental...
func (d *Client) Makeekscfg() client.Generator {
	return &client.GeneratorType{}
}

// Makeawscfg - creates the cfg object
func (d *Client) Makeawscfg() aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(d.Env.Region),
	)
	if err != nil {
		// handle error :(
		log.Panic("aws config management issue...")
	}
	return cfg
}

// Getekscluster - Return the cluster, endpoints and auth token!
func (d *Client) Getekscluster(c aws.Config, g client.Generator) client.EKSCluster {
	a := client.AwsConfig{J: d.Env}
	cluster := a.Init(c, g)
	return cluster
}

// Deleterelease will try and delete a release -> Need to reconfigure...
func (c *Client) Deleterelease(eks client.EKSCluster, a *action.Configuration, rel *release.Release, del internalhelm.HelmDelete) error {
	settings := cli.New()
	settings.KubeAPIServer = eks.Endpoint
	settings.KubeToken = eks.Token
	settings.KubeCaFile = eks.CAFile
	if err := a.Init(settings.RESTClientGetter(), rel.Namespace, os.Getenv("HELM_DRIVER"), log.Infof); err != nil {
		log.Fatal(err)
	}
	run := del.ActionNewUninstall(a)
	run.DryRun = c.Dryrun
	out, err := del.RunCommand(rel.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("deleted: %s", out.Info)

	if c.Dryrun {
		log.Infof("dry-run mode enabled - release: %s not actually deleted", rel.Name)
	} else {
		log.Infof("deleted: %s", out.Info)
	}
	return nil
}
