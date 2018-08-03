[![Go Report Card](https://goreportcard.com/badge/github.com/nullboundary/glfont)](https://goreportcard.com/report/github.com/nullboundary/glfont)
 
    Name    : glfont Library                      
    Author  : Noah Shibley, http://socialhardware.net                       
    Date    : June 16th 2016                                 
    Notes   : A modern opengl text rendering library for golang
    Dependencies:   freetype, go-gl, glfw

***
# Function List:

#### func  LoadFont

```go
func LoadFont(file string, scale int32, windowWidth int, windowHeight int) (*Font, error)
```
LoadFont loads the specified font at the given scale.

#### func  LoadTrueTypeFont

```go
func LoadTrueTypeFont(program uint32, r io.Reader, scale int32, low, high rune, dir Direction) (*Font, error)
```
LoadTrueTypeFont builds a set of textures based on a ttf files gylphs

#### func (*Font) Printf

```go
func (f *Font) Printf(x, y float32, scale float32, fs string, argv ...interface{}) error
```
Printf draws a string to the screen, takes a list of arguments like printf

#### func (*Font) SetColor

```go
func (f *Font) SetColor(red float32, green float32, blue float32, alpha float32)
```
SetColor allows you to set the text color to be used when you draw the text

#### func (f *Font) UpdateResolution

```go
func (f *Font) UpdateResolution(windowWidth int, windowHeight int)
```
UpdateResolution is needed when the viewport is resized

#### func (f *Font) Width

```go
func (f *Font) Width(scale float32, fs string, argv ...interface{}) float32
```
Width returns the width of a piece of text in pixels

***

# Example:

```go

package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/nullboundary/glfont"
)

const windowWidth = 1920
const windowHeight = 1080

func init() {
	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, _ := glfw.CreateWindow(int(windowWidth), int(windowHeight), "glfontExample", glfw.GetPrimaryMonitor(), nil)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	
	if err := gl.Init(); err != nil { 
		panic(err)
	}

	//load font (fontfile, font scale, window width, window height
	font, err := glfont.LoadFont("Roboto-Light.ttf", int32(52), windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

     //set color and draw text
		font.SetColor(1.0, 1.0, 1.0, 1.0) //r,g,b,a font color
		font.Printf(100, 100, 1.0, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.") //x,y,scale,string,printf args

		window.SwapBuffers()
		glfw.PollEvents()

	}
}
```

#### Contributors

* [kivutar](https://github.com/kivutar)
