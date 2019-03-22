// +build linux

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/carlogit/phash"
	"github.com/liamg/aminal/config"
	"github.com/liamg/aminal/gui"
	"github.com/liamg/aminal/terminal"
)

var termRef *terminal.Terminal
var guiRef *gui.GUI

func terminate(msg string) int {
	defer fmt.Print(msg)
	return 1
}

func sleep(seconds ...int) {
	count := 1
	if len(seconds) > 0 {
		count = seconds[0]
	}
	time.Sleep(time.Duration(count) * time.Second)
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

func imagesAreEqual(expected string, actual string) int {
	expectedHash := hash(expected)
	actualHash := hash(actual)
	return phash.GetDistance(expectedHash, actualHash)
}

func compareImages(expected string, actual string) {
	distance := imagesAreEqual(expected, actual)
	if distance != 0 {
		os.Exit(terminate(fmt.Sprintf("Screenshot \"%s\" doesn't match expected image \"%s\". Distance of hashes difference: %d\n",
			actual, expected, distance)))
	}
}

func send(terminal *terminal.Terminal, cmd string) {
	err := terminal.Write([]byte(cmd))
	if err != nil {
		panic(err)
	}
}

func enter(terminal *terminal.Terminal) {
	err := terminal.Write([]byte("\n"))
	if err != nil {
		panic(err)
	}
}

func validateScreen(img string, waitForChange bool) {
	fmt.Printf("taking screenshot: %s and comparing...", img)

	err := guiRef.Screenshot(img)
	if err != nil {
		panic(err)
	}

	compareImages(strings.Join([]string{"vttest/", img}, ""), img)

	fmt.Printf("compare OK\n")

	enter(termRef)

	if waitForChange {
		fmt.Print("Waiting for screen change...")
		attempts := 10
		for {
			sleep()
			actualScren := "temp.png"
			err = guiRef.Screenshot(actualScren)
			if err != nil {
				panic(err)
			}
			distance := imagesAreEqual(actualScren, img)
			if distance != 0 {
				break
			}
			fmt.Printf(" %d", attempts)
			attempts--
			if attempts <= 0 {
				break
			}
		}
		fmt.Print("done\n")
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	go m.Run()

	for f := range tests {
		f()
	}
}

var tests = make(chan func())

func runMain(f func()) {
	complete := make(chan bool, 1)
	tests <- func() {
		f()
		complete <- true
	}
	<-complete
}

func TestCursorMovement(t *testing.T) {
	runMain(func() {

		testFunc := func(term *terminal.Terminal, g *gui.GUI) {
			termRef = term
			guiRef = g

			sleep()
			send(term, "vttest\n")
			sleep()
			send(term, "1\n")
			sleep()

			term.Lock()
			compareResult := term.ActiveBuffer().CompareViewLines("vttest/test-cursor-movement-1")
			term.Unlock()

			if compareResult == false {
				os.Exit(terminate(fmt.Sprintf("ActiveBuffer doesn't match vttest template vttest/test-cursor-movement-1")))
			}

			validateScreen("test-cursor-movement-1.png", true)
			validateScreen("test-cursor-movement-2.png", true)
			validateScreen("test-cursor-movement-3.png", true)
			validateScreen("test-cursor-movement-4.png", true)
			validateScreen("test-cursor-movement-5.png", true)
			validateScreen("test-cursor-movement-6.png", false)

			g.Close()
		}

		initialize(testFunc, testConfig())
	})
}

func TestScreenFeatures(t *testing.T) {
	runMain(func() {

		testFunc := func(term *terminal.Terminal, g *gui.GUI) {
			termRef = term
			guiRef = g

			sleep()
			send(term, "vttest\n")
			sleep()
			send(term, "2\n")
			sleep()

			validateScreen("test-screen-features-1.png", true)
			validateScreen("test-screen-features-2.png", true)
			validateScreen("test-screen-features-3.png", true)
			validateScreen("test-screen-features-4.png", true)
			validateScreen("test-screen-features-5.png", true)
			validateScreen("test-screen-features-6.png", true)
			validateScreen("test-screen-features-7.png", true)
			validateScreen("test-screen-features-8.png", true)
			validateScreen("test-screen-features-9.png", true)
			validateScreen("test-screen-features-10.png", true)
			validateScreen("test-screen-features-11.png", true)
			validateScreen("test-screen-features-12.png", true)
			validateScreen("test-screen-features-13.png", true)
			validateScreen("test-screen-features-14.png", true)
			validateScreen("test-screen-features-15.png", false)

			g.Close()
		}

		initialize(testFunc, testConfig())
	})
}

func TestSixel(t *testing.T) {
	runMain(func() {

		testFunc := func(term *terminal.Terminal, g *gui.GUI) {
			termRef = term
			guiRef = g

			sleep()
			send(term, "export PS1='> '\n")
			sleep()
			send(term, "clear\n")
			sleep()
			send(term, "cat example.sixel\n")
			sleep(10) // Displaying SIXEL graphics *sometimes* takes long time. Excessive synchronization???

			validateScreen("test-sixel.png", false)

			g.Close()
		}

		initialize(testFunc, testConfig())
	})
}

// Last Test should terminate main goroutine since it's infinity looped to execute others GUI tests in main goroutine
func TestExit(t *testing.T) {
	os.Exit(0)
}

func testConfig() *config.Config {
	c := config.DefaultConfig()

	// Force the scrollbar not showing when running unit tests
	c.ShowVerticalScrollbar = false

	// Use a vanilla shell on POSIX to help ensure consistency.
	if runtime.GOOS != "windows" {
		c.Shell = "/bin/sh"
	}

	return c
}
