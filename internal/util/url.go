package util

import (
	"errors"
	"strings"
)

func ParamsFromUrl(url string) (metricName string, metricValue string, err error) {
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		err = errors.New("Unsupported url format")

		return
	}

	metricName = parts[3]
	metricValue = parts[4]

	return
}
