package config_test

import (
	"os"
	"testing"

	"pomodoro_cli/internal/config"
	"pomodoro_cli/internal/errors"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errCode errors.ErrorCode
	}{
		{
			name:    "valid minutes",
			args:    []string{"program", "25", "5"},
			wantErr: false,
		},
		{
			name:    "valid hours",
			args:    []string{"program", "1", "1", "-h"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = tt.args
			_, err := config.ParseArgs()

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
