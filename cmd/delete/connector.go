package delete

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	client "github.com/edify42/helm-janitor/internal/eks"
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
	Getekscluster(aws.Config) client.EKSCluster
	Deleterelease(*action.Configuration, *release.Release) error
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
func (d *Client) Getekscluster(c aws.Config) client.EKSCluster {
	a := client.AwsConfig{J: d.Env}
	cluster := a.Init(c)
	return cluster
}

// Deleterelease will try and delete a release -> Need to reconfigure...
func (c *Client) Deleterelease(a *action.Configuration, rel *release.Release) error {
	settings := cli.New()
	if err := a.Init(settings.RESTClientGetter(), rel.Namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v)
	}); err != nil {
		log.Fatal(err)
	}
	run := action.NewUninstall(a)
	run.DryRun = c.Dryrun
	out, err := run.Run(rel.Name)
	if err != nil {
		log.Fatal(err)
	}
	if c.Dryrun {
		log.Infof("dry-run mode enabled - release: %s not actually deleted", rel.Name)
	} else {
		log.Infof("deleted: %s", out.Info)
	}
	return nil
}
