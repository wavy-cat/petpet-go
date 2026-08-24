package loggerfactory_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/loggerfactory"
	"go.uber.org/zap/zapcore"
)

const serviceLogHelperEnv = "PETPET_GO_LOGGERFACTORY_SERVICE_LOG_HELPER"

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		preset       config.LoggerPreset
		debugEnabled bool
		infoEnabled  bool
	}{
		{
			name: "disabled",
		},
		{
			name:        "production",
			preset:      config.ProdPreset,
			infoEnabled: true,
		},
		{
			name:         "development",
			preset:       config.DevPreset,
			debugEnabled: true,
			infoEnabled:  true,
		},
		{
			name:        "GCP",
			preset:      config.GCPPreset,
			infoEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger, err := loggerfactory.New(config.Logger{Preset: tt.preset})
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if logger == nil {
				t.Fatal("New() logger = nil")
			}
			if got := logger.Core().Enabled(zapcore.DebugLevel); got != tt.debugEnabled {
				t.Errorf("New() debug enabled = %t, want %t", got, tt.debugEnabled)
			}
			if got := logger.Core().Enabled(zapcore.InfoLevel); got != tt.infoEnabled {
				t.Errorf("New() info enabled = %t, want %t", got, tt.infoEnabled)
			}
		})
	}
}

func TestNewUnknownPreset(t *testing.T) {
	t.Parallel()

	logger, err := loggerfactory.New(config.Logger{Preset: "structured"})
	if err == nil {
		t.Fatal("New() error = nil, want unsupported preset error")
	}
	if logger != nil {
		t.Fatalf("New() logger = %v, want nil", logger)
	}
	if got, want := err.Error(), `unsupported logger preset "structured"`; got != want {
		t.Fatalf("New() error = %q, want %q", got, want)
	}
}

func TestNewWithService(t *testing.T) {
	t.Parallel()

	if os.Getenv(serviceLogHelperEnv) == "1" {
		writeServiceLog(t)

		return
	}

	output := serviceLogFromSubprocess(t)

	var entry map[string]any
	if err := json.NewDecoder(bytes.NewReader(output)).Decode(&entry); err != nil {
		t.Fatalf("decode log entry %q: %v", output, err)
	}
	if got, want := entry["message"], "ready"; got != want {
		t.Errorf("log message = %v, want %q", got, want)
	}
	if got, want := entry["service"], "petpet-api"; got != want {
		t.Errorf("log service = %v, want %q", got, want)
	}
}

func TestNewWithServiceUnknownPreset(t *testing.T) {
	t.Parallel()

	logger, err := loggerfactory.NewWithService(
		config.Logger{Preset: "structured"},
		"petpet-api",
	)
	if err == nil {
		t.Fatal("NewWithService() error = nil, want unsupported preset error")
	}
	if logger != nil {
		t.Fatalf("NewWithService() logger = %v, want nil", logger)
	}
	if got, want := err.Error(), `unsupported logger preset "structured"`; got != want {
		t.Fatalf("NewWithService() error = %q, want %q", got, want)
	}
}

func writeServiceLog(t *testing.T) {
	t.Helper()

	logger, err := loggerfactory.NewWithService(
		config.Logger{Preset: config.GCPPreset},
		"petpet-api",
	)
	if err != nil {
		t.Fatalf("NewWithService() error = %v", err)
	}
	if logger == nil {
		t.Fatal("NewWithService() logger = nil")
	}

	logger.Info("ready")
}

func serviceLogFromSubprocess(t *testing.T) []byte {
	t.Helper()

	// The executable is the trusted test binary, not user-controlled input.
	cmd := exec.CommandContext(
		t.Context(),
		os.Args[0],
		"-test.run=^TestNewWithService$",
	)
	cmd.Env = append(os.Environ(), serviceLogHelperEnv+"=1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("service log helper error = %v, output = %q", err, output)
	}

	return output
}
