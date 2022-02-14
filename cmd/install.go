/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/spf13/cobra"
	"strings"
	"svm/io"
	"svm/parsers"
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
			err := getAllInstallableVersions()
			if err != nil {
				return err
			}
			return nil
		}
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

func getAllInstallableVersions() error {
	const sparkURL = "https://archive.apache.org/dist/spark/"
	fmt.Printf("Getting list of installable versions from %s\n", sparkURL)
	response, err := soup.Get(sparkURL)
	if err != nil {
		return errors.New("cannot connect to https://archive.apache.org/dist/spark, check internet connection")
	}
	links := getSiteInfo(response)
	for _, link := range links {
		if strings.HasPrefix(link.Text(), "spark") {
			subResponse, err := soup.Get(fmt.Sprintf("%s/%s", sparkURL, link.Text()))
			if err != nil {
				return errors.New("cannot connect to https://archive.apache.org/dist/spark, check internet connection")
			}
			subLinks := getSiteInfo(subResponse)
			for _, subLink := range subLinks {
				if strings.HasSuffix(subLink.Text(), ".tgz") {
					println(parsers.ParseSparkFilename(subLink.Text()))
				}
			}
		}
	}
	return nil
}

func getSiteInfo(resp string) []soup.Root {
	mainSite := soup.HTMLParse(resp)
	links := mainSite.FindAll("a")
	return links
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&listAllInstallable, "list", "l", false, "Lists all installable versions")
}
