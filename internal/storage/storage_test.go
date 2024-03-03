package storage

import "testing"

func TestMemStorage(t *testing.T) {
	tests := []struct {
		name string
		data map[string]string
		want map[string]string
	}{
		{
			name: "Test value doesn't exist",
			data: map[string]string{"test": "test"},
			want: map[string]string{"test2": ""},
		},
		{
			name: "Test success get value",
			data: map[string]string{"test": "test"},
			want: map[string]string{"test": "test"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewStorage()
			for k, v := range test.data {
				storage.AddMetric(k, v)
			}

			for k, v := range test.want {
				if res, _ := storage.GetMetric(k); res != v {
					t.Errorf("fail")
				}
			}
		})
	}
}
