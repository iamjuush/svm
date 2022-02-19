/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	// Used for flags.
	rootCmd = &cobra.Command{
		Use:   "svm",
		Short: "Yet another version manager, written in Go, to manage Apache Spark installations.",
		Long:  `Yet another version manager, written in Go, to manage Apache Spark installations.`,
		// Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

}
