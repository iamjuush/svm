/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	svmio "svm/io"
	"svm/parsers"
	"svm/web"
)

var listAllInstallable bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a specified version of Apache Spark",
	Long: `Installs a specified version of Apache Spark.

Example:
svm install 2.2.2-with-hadoop-2.7
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if listAllInstallable {
			return web.GetAllInstallableVersions()
		}
		if err := validateInstallArgs(args); err != nil {
			return err
		}
		sparkVersion := args[0]
		dirname, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		version := parsers.ParseSparkVersion(sparkVersion)
		url := parsers.GetURLFromVersion(sparkVersion)
		fmt.Printf("Fetching from: %s \n", url)
		resource := svmio.Resource{Url: url, Version: version, Home: dirname}
		if err := svmio.DownloadFile(resource); err != nil {
			return err
		}
		if err = svmio.UnzipTar(resource); err != nil {
			return err
		}
		return svmio.RenameUnzipped(resource)
	},
}

func validateInstallArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("no version specified. Use `svm install --list` to view all downloadable versions \n")
	}
	if len(args) > 1 {
		return errors.New("multiple versions not supported. Please only specify one version \n")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&listAllInstallable, "list", "l", false, "Lists all installable versions")
}
