package dbg

import "sync"
import "time"
import "fmt"

var mu sync.Mutex
var msgs []string

func D(msg string) {
	now := time.Now().Format("15:04:05.99999999")
	mu.Lock()
	msgs = append(msgs, fmt.Sprint(now, " ", msg))
	mu.Unlock()
}

func Dump() {
	mu.Lock()
	for _, msg := range msgs {
		fmt.Println(msg)
	}
	msgs = nil
	mu.Unlock()
}

func RunDumper() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			Dump()
		}
	}()
}
