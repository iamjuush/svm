package parsers

import (
	"fmt"
	"regexp"
	"strings"
)

type version struct {
	sparkVersion string
	fullVersion  string
}

// Extract versioning info from the input string for svm install. Few cases to consider -
// case 1: 3.2.0
// case 2: 3.2.0-hadoop2.7
// case 3: 3.2.0-without-hadoop
// case 4: 2.4.0-without-hadoop-scala-2.12.tgz
func parseSparkVersion(versionString string) version {
	var rCase1 = regexp.MustCompile(`^\d.\d.\d$`)
	var rCase2 = regexp.MustCompile(`^\d.\d.\d-hadoop\d.\d$`)
	var rCase3 = regexp.MustCompile(`^\d.\d.\d-without-hadoop$`)
	var rCase4 = regexp.MustCompile(`^\d.\d.\d-without-hadoop-scala-\d.\d{1,2}`)

	splitVersionString := strings.Split(versionString, "-")
	var fullVersion string
	switch {
	case rCase1.MatchString(versionString):
		fullVersion = fmt.Sprintf("spark-%s", versionString)
	case rCase2.MatchString(versionString):
		fullVersion = fmt.Sprintf("spark-%s-bin-%s", splitVersionString[0], splitVersionString[1])
	case rCase3.MatchString(versionString):
		fullVersion = fmt.Sprintf("spark-%s-bin-%s", splitVersionString[0], strings.Join(splitVersionString[1:], "-"))
	case rCase4.MatchString(versionString):
		fullVersion = fmt.Sprintf("spark-%s-bin-%s", splitVersionString[0], strings.Join(splitVersionString[1:], "-"))
	}

	sparkVersion := fmt.Sprintf("spark-%s", splitVersionString[0])

	return version{sparkVersion: sparkVersion, fullVersion: fullVersion}
}

func ParseSparkFilename(name string) string {
	name = strings.TrimSuffix(name, ".tgz")
	name = strings.TrimPrefix(name, "spark-")
	name = strings.Replace(name, "-bin", "", -1)
	return name
}

func GetURLFromVersion(version string) string {
	parsedVersion := parseSparkVersion(version)
	return fmt.Sprintf("https://archive.apache.org/dist/spark/%s/%s.tgz", parsedVersion.sparkVersion, parsedVersion.fullVersion)
}
