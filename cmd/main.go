package main

import (
	"os"

	"github.com/lendi-au/helm-janitor/cmd/delete"
	"github.com/lendi-au/helm-janitor/cmd/scan"
	"github.com/lendi-au/helm-janitor/internal/config"
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

	var cmdDeleteSet = &cobra.Command{
		Use:   "deleteset [k8s label selector]",
		Short: "delete a set of releases that match the label selector",
		Long: `deleteset searches all releases using syntax similar to a
helm list -l <key>=<value>,<key>=<value>... style command, except
instead of listing the releases it will delete the ones it finds`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hasSelector, err := scan.ValidateScanArg(args)
			if err != nil {
				log.Fatalf("Bad selector input: %v", err)
			}
			if !hasSelector {
				log.Fatalf("Please supply a selector to this function") // unlikely to ever run as minNargs = 1
			}
			scanner := scan.NewScanClient()
			scanner.Selector = args[0]
			scanner.Dryrun = dryRun
			scanner.AllNamespaces = allNamespaces
			scanner.Namespace = releaseNamespace
			scanner.IncludeNamespaces = includeNamespaces
			scanner.ExcludeNamespaces = excludeNamespaces
			scanner.Init()
			delete.RunDeleteSet(scanner)
		},
	}

	cmdDeleteSet.Flags().StringVarP(&releaseNamespace, "namespace", "n", "", "namespace of the releases to delete")
	cmdDeleteSet.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "activate dry run mode (don't actually do anything)")
	cmdDeleteSet.Flags().StringVarP(&includeNamespaces, "include-namespaces", "i", "", "delete releases in namespaces matching this expression")
	cmdDeleteSet.Flags().StringVarP(&excludeNamespaces, "exclude-namespaces", "e", "", "exclude from delete releases in namespaces matching this expression")
	cmdDeleteSet.Flags().BoolVarP(&allNamespaces, "all-namespaces", "a", true, "search all namespaces")

	var cmdScan = &cobra.Command{
		Use:   "scan [k8s label selector]",
		Short: "scan the k8s cluster and delete releases past the ttl",
		Long: `scan will search for any extra labels on helm release that match the input.
		By default we always search for the helm-janitor: true label`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			hasSelector, err := scan.ValidateScanArg(args)
			if err != nil {
				log.Fatalf("Bad selector input: %v", err)
			}
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

	var rootCmd = &cobra.Command{
		Use:   config.AppName,
		Short: `helm-janitor cleans helm releases in your k8s cluster`,
		Long:  `Read the docs on the usage for the appropriate use-case to fulfil your needs.`,
	}
	rootCmd.AddCommand(cmdDelete, cmdScan, cmdDeleteSet)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
