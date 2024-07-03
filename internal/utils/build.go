package utils

import "fmt"

func BuildInfo(buildVersion string, buildDate string, buildCommit string) {
	fmt.Println(fmt.Sprintf("Build version: %s", strValOrDefault(buildVersion, "N/A")))
	fmt.Println(fmt.Sprintf("Build date: %s", strValOrDefault(buildDate, "N/A")))
	fmt.Println(fmt.Sprintf("Build commit: %s", strValOrDefault(buildCommit, "N/A")))
}

func strValOrDefault(val string, defaultVal string) string {
	if val == "" {
		val = defaultVal
	}

	return val
}
