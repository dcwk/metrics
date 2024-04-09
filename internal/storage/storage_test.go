package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGauge(t *testing.T) {
	storage := NewStorage()
	tests := []struct {
		name string
		data map[string]float64
		want map[string]float64
		err  string
	}{
		{
			name: "Test fail gauge not found",
			data: map[string]float64{"test": 10.64},
			want: map[string]float64{"test2": 0},
			err:  "gauge not found",
		},
		{
			name: "Test can save gauge",
			data: map[string]float64{"test": 10.64},
			want: map[string]float64{"test": 10.64},
			err:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.data {
				err := storage.AddGauge(k, &v)
				assert.NoError(t, err)
			}

			for k, v := range test.want {
				res, err := storage.GetGauge(k, false)
				if test.err != "" {
					assert.Equal(t, test.err, err.Error())
				} else {
					assert.Equal(t, v, res)
				}
			}
		})
	}
}

func TestCounter(t *testing.T) {
	storage := NewStorage()
	tests := []struct {
		name string
		data map[string]int64
		want map[string]int64
		err  string
	}{
		{
			name: "Test fail counter not found",
			data: map[string]int64{"test": 10},
			want: map[string]int64{"test2": 0},
			err:  "counter not found",
		},
		{
			name: "Test can save counter",
			data: map[string]int64{"test": 10},
			want: map[string]int64{"test": 20},
			err:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.data {
				err := storage.AddCounter(k, &v)
				assert.NoError(t, err)
			}

			for k, v := range test.want {
				res, err := storage.GetCounter(k, false)
				if test.err != "" {
					assert.Equal(t, test.err, err.Error())
				} else {
					assert.Equal(t, v, res)
				}
			}
		})
	}
}

func TestGetMetrics(t *testing.T) {
	storage := NewStorage()
	tests := []struct {
		name string
		data map[string]float64
		want map[string]float64
		err  string
	}{
		{
			name: "Test fail gauge not found",
			data: map[string]float64{"test": 10.64},
			want: map[string]float64{"test2": 0},
			err:  "gauge not found",
		},
		{
			name: "Test can save gauge",
			data: map[string]float64{"test": 10.64},
			want: map[string]float64{"test": 10.64},
			err:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.data {
				err := storage.AddGauge(k, &v)
				assert.NoError(t, err)
			}

			for k, v := range test.want {
				res, err := storage.GetGauge(k, false)
				if test.err != "" {
					assert.Equal(t, test.err, err.Error())
				} else {
					assert.Equal(t, v, res)
				}
			}
		})
	}

	metricsList, _ := storage.GetJSONMetrics()
	assert.Equal(t, metricsList, `{"list":[{"id":"test","type":"gauge","value":10.64}]}`)
}
