package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParamsFromUrl(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		metricName  string
		metricValue string
	}{
		{
			name:        "Test can get params",
			url:         "/update/counter/someMetric/527",
			metricName:  "someMetric",
			metricValue: "527",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn, mv, err := ParamsFromUrl(tt.url)
			assert.NoError(t, err)
			assert.Equal(t, mn, tt.metricName)
			assert.Equal(t, mv, tt.metricValue)
		})
	}
}
