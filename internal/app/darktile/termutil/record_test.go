package termutil

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/creack/pty"
	"github.com/stretchr/testify/require"
)

var (
	recordName = flag.String("record.name", "", "fixture name to write under testdata/sessions")
	recordCmd  = flag.String("record.cmd", "", "command to run under a pty and capture")
	recordCols = flag.Uint("record.cols", 80, "pty width used while recording")
	recordRows = flag.Uint("record.rows", 24, "pty height used while recording")
)

// TestRecordSession captures the raw bytes a real program writes to a pty and
// stores them as a golden fixture. It is a manual tool, not part of the suite:
//
//	go test ./internal/app/darktile/termutil \
//	  -run TestRecordSession \
//	  -record.name 256-colour -record.cmd './scripts/256-colour.sh'
//
// Only record programs whose output is deterministic. Anything that prints a
// path, a hostname, a timestamp or reacts to the environment will produce a
// fixture that fails on someone else's machine.
func TestRecordSession(t *testing.T) {
	if *recordCmd == "" {
		t.Skip("set -record.cmd and -record.name to capture a fixture")
	}
	require.NotEmpty(t, *recordName, "-record.name is required")

	cmd := exec.Command("sh", "-c", *recordCmd)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Cols: uint16(*recordCols),
		Rows: uint16(*recordRows),
	})
	require.NoError(t, err)
	defer func() { _ = ptmx.Close() }()

	// The read ends when the child closes the pty, which surfaces as EIO on
	// Linux rather than a clean EOF.
	captured, _ := io.ReadAll(ptmx)
	_ = cmd.Wait()

	require.NotEmpty(t, captured, "command produced no output")

	dir := filepath.Join("testdata", "sessions")
	require.NoError(t, os.MkdirAll(dir, 0o755))

	path := filepath.Join(dir, *recordName+".vt")
	require.NoError(t, os.WriteFile(path, captured, 0o644))

	t.Logf("captured %d bytes to %s — add it to goldenSessions, then run with -update", len(captured), path)
}
