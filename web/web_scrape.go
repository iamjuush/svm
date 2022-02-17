package web

import (
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"strings"
	"svm/parsers"
)

func GetAllInstallableVersions() error {
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
					println(parsers.SparkToSVMFilename(subLink.Text()))
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
