package consolelog_test

import (
	"testing"

	"gitlabFileScanner/internal/infrastructure/consolelog"
)

func TestLogger_Info_NoPanic(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	// These should not panic
	l.Info("simple message")
	l.Info("formatted: %s", "value")
	l.Info("")
	l.Info("many args: %s %d %v", "str", 42, true)
}

func TestLogger_Error_NoPanic(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	l.Error("error message")
	l.Error("error with detail: %v", 42)
	l.Error("")
}

func TestLogger_Warn_NoPanic(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	l.Warn("warning message")
	l.Warn("warning with detail: %s", "detail")
	l.Warn("")
}
