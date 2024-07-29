package utils

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Info struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

func TestBuildInfo(t *testing.T) {
	tests := []struct {
		Name string
		Data Info
		Want string
	}{
		{
			Name: "get build info with default values",
			Data: Info{
				BuildVersion: "",
				BuildDate:    "",
				BuildCommit:  "",
			},
			Want: `Build version: N/A
Build date: N/A
Build commit: N/A
`,
		},
		{
			Name: "get build info with default values",
			Data: Info{
				BuildVersion: "1.2.3",
				BuildDate:    "03.07.2024",
				BuildCommit:  "Ololo",
			},
			Want: `Build version: 1.2.3
Build date: 03.07.2024
Build commit: Ololo
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			stdout, err := captureOutput(
				func() error {
					BuildInfo(tt.Data.BuildVersion, tt.Data.BuildDate, tt.Data.BuildCommit)
					return nil
				})
			assert.NoError(t, err)
			assert.Equal(t, stdout, tt.Want)
		})
	}
}

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	if err != nil {
		return "", err
	}

	os.Stdout = orig
	err = w.Close()
	if err != nil {
		return "", err
	}

	out, _ := io.ReadAll(r)
	return string(out), err
}
