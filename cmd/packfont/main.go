package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Please specify font path and name")
		os.Exit(1)
	}

	path := os.Args[1]
	name := os.Args[2]

	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	packed, err := os.OpenFile(fmt.Sprintf("./internal/app/darktile/packed/%s.go", strings.ToLower(name)), os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		panic(err)
	}
	defer packed.Close()

	if _, err := packed.WriteString(fmt.Sprintf(`package packed

var %sTTF = []byte{
`, name)); err != nil {
		panic(err)
	}

	for i, b := range fontBytes {
		if _, err := packed.WriteString(fmt.Sprintf(" 0x%x,", b)); err != nil {
			panic(err)
		}
		if i > 0 && i%16 == 0 {
			if _, err := packed.WriteString("\n"); err != nil {
				panic(err)
			}
		}
	}
	if _, err := packed.WriteString("}\n"); err != nil {
		panic(err)
	}

}
