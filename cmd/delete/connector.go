package delete

import (
	"fmt"
	"os"

	"github.com/prometheus/common/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

// Client is the data object which contains the item to delete
type Client struct {
	Release   string
	Namespace string
}

// NewClient will return the Client struct
func NewClient() *Client {
	return &Client{}
}

// InputRun is our interface which defines the main delete methods
type InputRun interface {
	Deleterelease()
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
	run.DryRun = false
	out, err := run.Run(rel.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("deleted: %s", out.Info)
	return nil
}
