package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGauge(t *testing.T) {
	tests := []struct {
		name string
		data map[string]string
		want map[string]float64
		err  string
	}{
		{
			name: "Test fail gauge not found",
			data: map[string]string{"test": "10.64"},
			want: map[string]float64{"test2": 0},
			err:  "gauge not found",
		},
		{
			name: "Test can save gauge",
			data: map[string]string{"test": "10.64"},
			want: map[string]float64{"test": 10.64},
			err:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewStorage()
			for k, v := range test.data {
				err := storage.AddGauge(k, v)
				assert.NoError(t, err)
			}

			for k, v := range test.want {
				res, err := storage.GetGauge(k)
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
	tests := []struct {
		name string
		data map[string]string
		want map[string]int64
		err  string
	}{
		{
			name: "Test fail counter not found",
			data: map[string]string{"test": "10"},
			want: map[string]int64{"test2": 0},
			err:  "counter not found",
		},
		{
			name: "Test can save counter",
			data: map[string]string{"test": "10"},
			want: map[string]int64{"test": 10},
			err:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewStorage()
			for k, v := range test.data {
				err := storage.AddCounter(k, v)
				assert.NoError(t, err)
			}

			for k, v := range test.want {
				res, err := storage.GetCounter(k)
				if test.err != "" {
					assert.Equal(t, test.err, err.Error())
				} else {
					assert.Equal(t, v, res)
				}
			}
		})
	}
}
