package main

import (
	"fmt"
	"strings"

	"github.com/edify42/helm-janitor/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	var echoTimes int
	var releaseNamespace string
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
	cmdDelete.MarkPersistentFlagRequired("namespace")

	var cmdScan = &cobra.Command{
		Use:   "scan [k8s label selector]",
		Short: "scan the k8s cluster and delete releases past the ttl",
		Long: `scan will search for any extra labels on helm release that match the input.
		By default we always search for the helm-janitor: true label`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("scanning for releases which match: " + strings.Join(args, " "))
		},
	}

	cmdScan.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: config.AppName}
	rootCmd.AddCommand(cmdDelete, cmdScan)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
