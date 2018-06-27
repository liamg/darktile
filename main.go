package main

import (
	"io"
	"os"

	"gitlab.com/liamg/terminal/pty"
)

func main() {

	pty, err := pty.NewPtyWithShell()
	if err != nil {
		panic(err)
	}

	go io.Copy(pty, os.Stdin)
	io.Copy(os.Stdout, pty)

	//	return pty, err
}
