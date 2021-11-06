//go:build cgo && (!linux || unix)
package hinters

/*
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <sys/timespec.h>


time_t getuptime() {
	struct timespec tp;
	clock_gettime(CLOCK_UPTIME, &tp);
	return tp.tv_sec;
}
*/
import "C"

func getUptime() int64 {
	time := C.getuptime()
	return int64(time)
}
