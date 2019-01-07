package w32

import (
	"testing"
)

var testProcess = "notepad.exe"
var wantCode = uint32(42)

func TestCreateProcess(t *testing.T) {

	pi, err := CreateProcessQuick(testProcess)
	if err != nil {
		t.Errorf("[!!] Failed to create %s: %s", testProcess, err)
	} else {
		t.Logf("[OK] Created process %s with handle 0x%x, PID %d", testProcess, pi.Process, pi.ProcessId)
	}

	err = TerminateProcess(pi.Process, wantCode)
	if err != nil {
		t.Errorf("[!!]Failed to terminate %s: %s", testProcess, err)
	} else {
		t.Logf("[OK] Called TerminateProcess on PID %d", pi.ProcessId)
	}

	err = WaitForSingleObject(pi.Process, 1000) // 1000ms
	if err != nil {
		t.Errorf("[!!] failed in WaitForSingleObject: %s", err)
	} else {
		t.Logf("[OK] WaitForSingleObject returned...")
	}

	// make sure we see the magic exit code we asked for
	code, err := GetExitCodeProcess(pi.Process)
	if err != nil {
		t.Errorf("[!!] Failed to get exit code for PID %d: %s", pi.ProcessId, err)
	} else {
		t.Logf("[OK] PID %d Exited with code %d", pi.ProcessId, code)
	}
	if code != 42 {
		t.Errorf("[!!] Unexpected exit code for PID %d - want %d, got %d", pi.ProcessId, wantCode, code)
	}

	CloseHandle(pi.Process)
	CloseHandle(pi.Thread)

}
