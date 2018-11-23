# loginshell

A golang library to get the login shell of the current user.

## Supporting GOOS

Need PR for the other OS!

- Linux
- Darwin

## Usage

    package main

    import (
        "github.com/riywo/loginshell"
        "fmt"
    )
    
    func main() {
        shell, err := loginshell.Shell()
        if err != nil { panic(err) }
        fmt.Printf("%s", shell)
    }

## License

MIT

## Author

Ryosuke IWANAGA a.k.a. riywo