package w32

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

var forkFn = path.Join(os.TempDir(), "forktest.pid")

func TestFork(t *testing.T) {

	ppid := os.Getpid()
	t.Logf("[OK] I am PID %d", ppid)
	pid, err := Fork()
	if err != nil {
		t.Fatalf("[!!] Failed to fork. PID: %d: %s", pid, err)
	}

	if pid == 0 {
		// We can't log anything here because our stdout doesn't point
		// to the same console as our parent.
		//
		// This process won't show up in Task Manager, and os.Getpid() won't
		// work, I guess because we haven't told CSRSS we exist.
		f, _ := os.Create(forkFn)
		f.WriteString(fmt.Sprintf("%d", ppid))
		f.Close()
	} else {
		t.Logf("[OK] Forked child with PID %d", pid)
		t.Logf("[OK] Sleeping, then trying to read checkfile.")
		time.Sleep(2 * time.Second)
		raw, err := ioutil.ReadFile(forkFn)
		if err != nil {
			t.Fatalf("[!!] Failed to read PID checkfile: %s", err)
		}
		if string(raw) == strconv.Itoa(ppid) {
			t.Logf("[OK] Found PID checkfile - PID matches!")
		} else {
			t.Errorf("[!] Child reported PID %q vs %q!", string(raw), strconv.Itoa(ppid))
		}
		os.Remove(forkFn)

	}

}
