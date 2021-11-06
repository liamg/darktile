package hinters

import (
	"syscall"
)

func getUptime() int64 {
	sysInfo := &syscall.Sysinfo_t{}
	_ = syscall.Sysinfo(sysInfo)
	return sysInfo.Uptime
}
