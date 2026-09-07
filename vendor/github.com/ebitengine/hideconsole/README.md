# HideConsole

Package hideconsole is a utility package to hide a console automatically even without `-ldflags "-Hwindowsgui"` on Windows.

On non-Windows, this package does nothing.

## Usage

Import this package with a blank identifier.

```go
import (
    _ "github.com/ebitengine/hideconsole"
)
```
