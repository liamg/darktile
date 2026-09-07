// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Ebitengine Authors

// Package hideconsole is a utility package to hide a console automatically
// even without `-ldflags "-Hwindowsgui"` on Windows.
//
// On non-Windows, this package does nothing.
//
// # Usage
//
// Import this package with a blank identifier:
//
//	import (
//		_ "github.com/ebitengine/hideconsole"
//	)
package hideconsole
