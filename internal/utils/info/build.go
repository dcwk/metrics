package info

import "fmt"

func BuildInfo(buildVersion string, buildDate string, buildCommit string) {
	fmt.Printf("Build version: %s\n", strValOrDefault(buildVersion, "N/A"))
	fmt.Printf("Build date: %s\n", strValOrDefault(buildDate, "N/A"))
	fmt.Printf("Build commit: %s\n", strValOrDefault(buildCommit, "N/A"))
}

func strValOrDefault(val string, defaultVal string) string {
	if val == "" {
		val = defaultVal
	}

	return val
}
