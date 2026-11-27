package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    config
		wantErr bool
	}{
		{
			name: "defaults",
			args: nil,
			want: config{serverURL: "http://localhost:8080", reportInterval: 10 * time.Second, pollInterval: 2 * time.Second},
		},
		{
			name: "custom values",
			args: []string{"-a=localhost:9090", "-r=5", "-p=1"},
			want: config{serverURL: "http://localhost:9090", reportInterval: 5 * time.Second, pollInterval: time.Second},
		},
		{
			name:    "unknown flag",
			args:    []string{"-x=1"},
			wantErr: true,
		},
		{
			name:    "non-positive interval",
			args:    []string{"-p=0"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFlags(tt.args)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
