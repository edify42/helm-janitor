package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lendi-au/helm-janitor/cmd/delete"
	"github.com/lendi-au/helm-janitor/cmd/scan"
	"github.com/lendi-au/helm-janitor/internal/config"
	events "github.com/lendi-au/helm-janitor/pkg/lambda"
	log "github.com/sirupsen/logrus"
)

// runs the generic handler to execute helm delete...
// when the ttl expires.

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	logLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("Dodgy log level set: %s", logLevel)
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(level)
	}
}

type Test struct {
	Timmy string `json:"timmy"`
}

// EventBody is the git webhook event send by Stack Janitor
type EventBody struct {
	Name string `json:"name"`
	Time Test   `json:"time"`
}

// HandleRequest is the main lambda handler
func HandleRequest(ctx context.Context, event interface{}) error {
	scanner := scan.NewScanClient()
	// fmt.Println(reflect.TypeOf(event))
	test, _ := json.Marshal(event)
	// log.Debugf(string(test))
	switch event := event.(type) {
	case nil:
		log.Fatal("event is nil")
	case string:
		log.Fatalf("event was a string: %s", event)
	case EventBody:
		log.Infof("what kind of event: %v", event.Name)
		scanner.Selector = event.Name
	case events.GithubWebhookEvent:
		log.Debugf("my action is a %v with pr %v and repo %v", event.Action, event.PullRequest, event.Repository)
		a := fmt.Sprintf("BRANCH=%s,REPOSITORY=%s", event.PullRequest.State, event.PullRequest.Head.Repository.Name)
		scanner.Selector = a
	case events.BitbucketWebhookEvent:
		log.Debugf("my pr %v on the repo %v", event.PullRequest, event.Repository)
		a := fmt.Sprintf("BRANCH=%s,REPOSITORY=%s", event.PullRequest.Source.Branch.Name, event.Repository.Name)
		scanner.Selector = a
	default:
		a := new(events.BitbucketWebhookEvent)
		_ = json.Unmarshal(test, a)
		log.Infof("tried: %s on branch %s", a.Repository.Name, a.PullRequest.Source.Branch.Name)
		b := fmt.Sprintf("BRANCH=%s,REPOSITORY=%s", a.PullRequest.Source.Branch.Name, a.Repository.Name)
		scanner.Selector = b
	}
	scanner.Dryrun = config.GetenvWithDefaultBool("DRY_RUN", false)
	scanner.AllNamespaces = config.GetenvWithDefaultBool("ALL_NAMESPACES", true)
	scanner.Namespace = config.GetenvWithDefault("NAMESPACE", "")
	scanner.IncludeNamespaces = config.GetenvWithDefault("INCLUDE_NAMESPACES", "")
	scanner.ExcludeNamespaces = config.GetenvWithDefault("EXCLUDE_NAMESPACES", "")
	scanner.Context = ctx
	scanner.Init()
	delete.RunDeleteSet(scanner)
	return nil
}

func main() {
	log.Infof("starting")
	if os.Getenv("DEBUG") == "true" {
		ctx := context.Background()
		HandleRequest(ctx, EventBody{
			Name: "BRANCH=janitor-testing,REPOSITORY=router,helm-janitor=true",
			Time: Test{
				Timmy: "now",
			},
		})
	} else {
		lambda.Start(HandleRequest)
	}
	log.Infof("finished")
}
