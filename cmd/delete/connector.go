package delete

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	janitorconfig "github.com/lendi-au/helm-janitor/internal/config"
	client "github.com/lendi-au/helm-janitor/internal/eks"
	internalhelm "github.com/lendi-au/helm-janitor/internal/helm"
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
func (c *Client) Init() {
	test := janitorconfig.EnvConfig{}
	test.Init() // get the default values...again.
	c.Env = test
	log.Infof("Delete client initialised with values %v", c)
}

// Makeekscfg returns an empty EKS config
func (c *Client) Makeekscfg() client.Generator {
	return &client.GeneratorType{}
}

// Makeawscfg - creates the cfg object
func (c *Client) Makeawscfg() aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(c.Env.Region),
	)
	if err != nil {
		// handle error :(
		log.Panic("aws config management issue...")
	}
	return cfg
}

// Getekscluster - Return the cluster, endpoints and auth token!
func (c *Client) Getekscluster(aws aws.Config, g client.Generator) client.EKSCluster {
	a := client.AwsConfig{J: c.Env}
	cluster := a.Init(aws, g)
	return cluster
}

// Deleterelease will try and delete a release -> Need to reconfigure...
func (c *Client) Deleterelease(eks client.EKSCluster, a *action.Configuration, rel *release.Release, del internalhelm.HelmDelete) error {
	os.Setenv("HELM_NAMESPACE", rel.Namespace) // holy hack batman. why don't they expose this at the API level?
	settings := cli.New()
	settings.KubeAPIServer = eks.Endpoint
	settings.KubeToken = eks.Token
	settings.KubeCaFile = eks.CAFile
	if err := a.Init(settings.RESTClientGetter(), rel.Namespace, "secrets", log.Infof); err != nil {
		log.Error(err)
		return err
	}
	run := del.ActionNewUninstall(a)
	run.DryRun = c.Dryrun
	out, err := del.RunCommand(rel.Name)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("deleted: %s", out.Info)

	if c.Dryrun {
		log.Infof("dry-run mode enabled - release: %s not actually deleted", rel.Name)
	} else {
		log.Infof("deleted: %s", out.Info)
	}
	return nil
}
