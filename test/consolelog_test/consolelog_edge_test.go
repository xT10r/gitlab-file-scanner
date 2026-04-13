package consolelog_test

import (
	"testing"

	"gitlabFileScanner/internal/infrastructure/consolelog"
)

func TestLogger_Info_WithPercentInMsg(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	// Эти вызовы не должны вызывать panic
	l.Info("100% done")
	l.Info("test %s %d without args")
	l.Info("test %s %d", "value", 42)
}

func TestLogger_Error_WithPercentInMsg(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	l.Error("failed at 50%")
	l.Error("error %s", "detail")
}

func TestLogger_Warn_WithPercentInMsg(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	l.Warn("disk 90% full")
	l.Warn("deprecated %s", "func")
}

func TestLogger_Debug_WithPercentInMsg(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	l.Debug("progress: 75%")
	l.Debug("request %d", 123)
}

func TestLogger_Constants(t *testing.T) {
	t.Parallel()

	if consolelog.LevelDebug != "DEBUG: " {
		t.Errorf("LevelDebug = %q, want %q", consolelog.LevelDebug, "DEBUG: ")
	}
	if consolelog.LevelInfo != "INFO: " {
		t.Errorf("LevelInfo = %q, want %q", consolelog.LevelInfo, "INFO: ")
	}
	if consolelog.LevelWarn != "WARN: " {
		t.Errorf("LevelWarn = %q, want %q", consolelog.LevelWarn, "WARN: ")
	}
	if consolelog.LevelError != "ERROR: " {
		t.Errorf("LevelError = %q, want %q", consolelog.LevelError, "ERROR: ")
	}
}
