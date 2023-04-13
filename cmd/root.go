package main

import (
	"github.com/mikezuff/apidemo/pkg/api"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apidemo",
	Short: "demo API",
	Long:  "run the demo API (eqivalent to `apidemo run`)",
	Run:   apiDemoRun,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the demo API",
	Long:  "run the demo API",
	Run:   apiDemoRun,
}

func apiDemoRun(cmd *cobra.Command, args []string) {
	api.Run()
}

func init() {
	rootCmd.AddCommand(runCmd)
}
