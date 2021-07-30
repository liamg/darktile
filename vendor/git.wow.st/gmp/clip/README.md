# Clip

A tiny library to access the MacOS clipboard (a.k.a. NSPasteboard).

```go
go get git.wow.st/gmp/clip
```

## API:

```go
package clip

// Clear clears the general pasteboard
func Clear()

// Set puts a string on the pasteboard, returning true if successful
func Set(string) bool

// Get retrieves the string currently on the pasteboard.
func Get() string
```

