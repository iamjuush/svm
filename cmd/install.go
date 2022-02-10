/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"svm/io"
	"svm/parsers"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a specified version of Apache Spark",
	Long: `Installs a specified version of Apache Spark.

Example:
svm install 2.2.2-with-hadoop-2.7
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateArgs(args); err != nil {
			return err
		}
		VERSION := args[0]
		url := parsers.GetURLFromVersion(VERSION)
		fmt.Printf("Fetching from: %s \n", url)
		resource := io.Resource{Filename: "spark.tgz", Url: url}
		err := io.DownloadFile(resource)
		if err != nil {
			return err
		}
		return nil
	},
}

func validateArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("no version specified. Use `svm list` to view available versions \n")
	}
	if len(args) > 1 {
		return errors.New("multiple versions not supported. Please only specify one version \n")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}
