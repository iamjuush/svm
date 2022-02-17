/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
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
			err := web.GetAllInstallableVersions()
			if err != nil {
				return err
			}
			return nil
		}
		if err := validateArgs(args); err != nil {
			return err
		}
		version := args[0]
		url := parsers.GetURLFromVersion(version)
		fmt.Printf("Fetching from: %s \n", url)
		resource := svmio.Resource{Filename: version + ".tgz", Url: url}
		err := svmio.DownloadFile(resource, version)
		if err != nil {
			return err
		}
		err = svmio.Untar(version)
		err = svmio.RenameUnzipped(version)
		if err != nil {
			return err
		}
		//TODO: Rethink the file namings during initial file creation (tmp), after download (.tgz), and final (folder w/ name as version)
		//TODO: Delete tar file
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
	installCmd.Flags().BoolVarP(&listAllInstallable, "list", "l", false, "Lists all installable versions")
}
