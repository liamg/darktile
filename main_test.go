// +build linux

package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/liamg/aminal/gui"
	"github.com/liamg/aminal/terminal"

	"github.com/carlogit/phash"
)

func terminate(msg string) int {
	defer fmt.Print(msg)
	return 1
}

func sleep() {
	time.Sleep(time.Second)
}

func hash(path string) string {
	image, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer image.Close()
	imageHash, err := phash.GetHash(image)
	if err != nil {
		panic(err)
	}
	return imageHash
}

func compareImages(img1 string, img2 string) {
	template := hash(img1)
	screen := hash(img2)
	distance := phash.GetDistance(template, screen)
	if distance != 0 {
		os.Exit(terminate(fmt.Sprintf("Screenshot doesn't match expected image. Distance of hashes difference: %d\n", distance)))
	}
}

func send(terminal *terminal.Terminal, cmd string) {
	terminal.Write([]byte(cmd))
}

func enter(terminal *terminal.Terminal) {
	terminal.Write([]byte("\n"))
}

func TestCursorMovement(t *testing.T) {

	testFunc := func(term *terminal.Terminal, g *gui.GUI) {
		sleep()
		send(term, "vttest\n")
		sleep()
		send(term, "1\n")
		sleep()

		if term.ActiveBuffer().Compare("vttest/test-cursor-movement-1") == false {
			t.Error(fmt.Sprint("ActiveBuffer doesn't match vttest template"))
		}

		enter(term)
		sleep()
		g.Screenshot ("test-cursor-movement-2.png")
		compareImages("vttest/test-cursor-movement-2.png", "test-cursor-movement-2.png")
		os.Exit(0)
	}

	initialize(testFunc)
}
