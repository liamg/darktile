// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

type FontLike interface {
	GetTextureWidth() float32
	GetTextureHeight() float32
}
