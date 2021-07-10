package main

import (
	"os"

	"github.com/edify42/helm-janitor/cmd/delete"
	"github.com/edify42/helm-janitor/cmd/scan"
	"github.com/edify42/helm-janitor/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

func main() {
	var releaseNamespace string
	var allNamespaces bool
	var dryRun bool
	var includeNamespaces string
	var excludeNamespaces string
	// logging setup
	var cmdDelete = &cobra.Command{
		Use:   "delete [release]",
		Short: "delete a specific helm release",
		Long:  `simple wrapper around helm delete which can be used by this tool`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			delete.ValidateArgs(args)
			purger := delete.NewClient()
			log.Infof("Deleting the release: %s in namespace %s", args[0], releaseNamespace)
			purger.Release = args[0]
			purger.Namespace = releaseNamespace
			purger.Dryrun = dryRun
			purger.Init()
			delete.RunV2(purger)
		},
	}

	cmdDelete.PersistentFlags().StringVarP(&releaseNamespace, "namespace", "n", "", "namespace of the release")
	cmdDelete.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "activate dry run mode (don't actually do anything)")
	// cmdDelete.Flags().BoolVarP(&allNamespaces, "all-namespaces", "a", false, "search all namespaces")
	cmdDelete.MarkPersistentFlagRequired("namespace")

	var cmdScan = &cobra.Command{
		Use:   "scan [k8s label selector]",
		Short: "scan the k8s cluster and delete releases past the ttl",
		Long: `scan will search for any extra labels on helm release that match the input.
		By default we always search for the helm-janitor: true label`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			hasSelector := scan.ValidateScanArg(args)
			scanner := scan.NewScanClient()
			if hasSelector {
				log.Infof("Scanning for releases which match: %s", args[0])
				scanner.Selector = args[0]
			}
			scanner.Dryrun = dryRun
			scanner.AllNamespaces = allNamespaces
			scanner.Namespace = releaseNamespace
			scanner.Init()
			scan.RunV2(scanner)
		},
	}
	cmdScan.Flags().StringVarP(&releaseNamespace, "namespace", "n", "", "namespace of the release")
	cmdScan.Flags().BoolVarP(&allNamespaces, "all-namespaces", "a", true, "search all namespaces")
	cmdScan.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "activate dry run mode (don't actually do anything)")
	cmdScan.Flags().StringVarP(&includeNamespaces, "include-namespaces", "i", "", "search for releases in namespaces matching this expression")
	cmdScan.Flags().StringVarP(&excludeNamespaces, "exclude-namespaces", "e", "", "exclude releases in namespaces matching this expression")

	var rootCmd = &cobra.Command{Use: config.AppName}
	rootCmd.AddCommand(cmdDelete, cmdScan)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
