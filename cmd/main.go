package main

import (
	"fmt"
	"strings"

	"github.com/edify42/helm-janitor/cmd/scan"
	"github.com/edify42/helm-janitor/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var releaseNamespace string
	var allNamespaces bool
	var includeNamespaces string
	var excludeNamespaces string
	// logging setup
	log.SetFormatter(&log.JSONFormatter{})
	var cmdDelete = &cobra.Command{
		Use:   "delete [release]",
		Short: "delete a specific helm release",
		Long:  `simple wrapper around helm delete which can be used by this tool`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Delete function mock: " + strings.Join(args, " "))
		},
	}

	cmdDelete.PersistentFlags().StringVarP(&releaseNamespace, "namespace", "n", "", "namespace of the release")
	cmdDelete.Flags().BoolVarP(&allNamespaces, "all-namespaces", "a", false, "search all namespaces")
	cmdDelete.MarkPersistentFlagRequired("namespace")
	cmdDelete.Flags().StringVarP(&includeNamespaces, "include-namespaces", "i", "", "search for releases in namespaces matching this expression")
	cmdDelete.Flags().StringVarP(&excludeNamespaces, "exclude-namespaces", "e", "", "exclude releases in namespaces matching this expression")

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
			scanner.Init()
			scan.RunV2(scanner)
		},
	}
	cmdScan.Flags().StringVarP(&releaseNamespace, "namespace", "n", "", "namespace of the release")
	cmdScan.Flags().BoolVarP(&allNamespaces, "all-namespaces", "a", true, "search all namespaces")

	var rootCmd = &cobra.Command{Use: config.AppName}
	rootCmd.AddCommand(cmdDelete, cmdScan)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
