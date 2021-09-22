package main

import (
	"context"
	"os"

	"github.com/lendi-au/helm-janitor/cmd/scan"
	"github.com/lendi-au/helm-janitor/internal/config"
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

// HandleRequest runs the scan package code to look for old releases and deletes them
func HandleRequest() error {
	ctx := context.Background()

	scanner := scan.NewScanClient()
	scanner.Dryrun = config.GetenvWithDefaultBool("DRY_RUN", false)
	scanner.AllNamespaces = config.GetenvWithDefaultBool("ALL_NAMESPACES", true)
	scanner.Namespace = config.GetenvWithDefault("NAMESPACE", "")
	scanner.IncludeNamespaces = config.GetenvWithDefault("INCLUDE_NAMESPACES", "")
	scanner.ExcludeNamespaces = config.GetenvWithDefault("EXCLUDE_NAMESPACES", "")
	scanner.Context = ctx
	scanner.Init()
	log.Info("starting...")
	scan.RunV2(scanner)
	return nil
}

func main() {
	log.Infof("starting")
	HandleRequest()
	// lambda.Start(HandleRequest)
	log.Infof("finished")
}
