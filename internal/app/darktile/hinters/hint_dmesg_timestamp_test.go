package hinters

import (
	"testing"
	"time"

	"github.com/liamg/darktile/internal/app/darktile/termutil"
	"github.com/stretchr/testify/assert"
)

func Test_dmesg_hinter_resolves_timestamp(t *testing.T) {

	hinter := &DmesgTimestampHinter{}
	api := &TestAPI{}

	// force system uptime to a known time for consistency in the test
	sysStart = time.Date(2021, 5, 21, 1, 1, 1, 1, time.UTC)

	input := `[47028.800212] audit: type=1107 audit(1621609618.710:91): pid=1011 uid=103 auid=4294967295 ses=4294967295 subj=unconfined msg='apparmor="DENIED" operation="dbus_signal"  bus="system" path="/org/freedesktop/NetworkManager" interface="org.freedesktop.NetworkManager" member="PropertiesChanged" name=":1.12" mask="receive" pid=339538 label="snap.spotify.spotify" peer_pid=1012 peer_label="unconfined"`

	match, offset, length := hinter.Match(input, 4)

	assert.Equal(t, true, match)
	hinter.Activate(api, input[offset:offset+length], termutil.Position{}, termutil.Position{})

	assert.Equal(t, "Fri May 21 14:04:49 2021", api.highlighted)
}

func Test_dmesg_hinter_resolves_timestamp_then_clears(t *testing.T) {

	hinter := &DmesgTimestampHinter{}
	api := &TestAPI{}

	// force system uptime to a known time for consistency in the test
	sysStart = time.Date(2021, 5, 21, 1, 1, 1, 1, time.UTC)

	input := `[47028.800212] audit: type=1107 audit(1621609618.710:91): pid=1011 uid=103 auid=4294967295 ses=4294967295 subj=unconfined msg='apparmor="DENIED" operation="dbus_signal"  bus="system" path="/org/freedesktop/NetworkManager" interface="org.freedesktop.NetworkManager" member="PropertiesChanged" name=":1.12" mask="receive" pid=339538 label="snap.spotify.spotify" peer_pid=1012 peer_label="unconfined"`

	match, offset, length := hinter.Match(input, 4)

	assert.Equal(t, true, match)
	hinter.Activate(api, input[offset:offset+length], termutil.Position{}, termutil.Position{})

	assert.Equal(t, "Fri May 21 14:04:49 2021", api.highlighted)
	hinter.Deactivate(api)
	assert.Equal(t, "", api.highlighted)
}
