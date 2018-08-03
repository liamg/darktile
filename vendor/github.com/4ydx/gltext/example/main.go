package main

import (
	"fmt"
	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/math/fixed"
	"math"
	"os"
	"runtime"
	"time"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic("glfw error")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	if useStrictCoreProfile {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	}
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Opengl version", version)

	// code from here
	gltext.IsDebug = true

	var font *v41.Font
	config, err := gltext.LoadTruetypeFontConfig("fontconfigs", "font_1_honokamin")
	if err == nil {
		font, err = v41.NewFont(config)
		if err != nil {
			panic(err)
		}
		fmt.Println("Font loaded from disk...")
	} else {
		fd, err := os.Open("font/font_1_honokamin.ttf")
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		// Japanese character ranges
		// http://www.rikai.com/library/kanjitables/kanji_codes.unicode.shtml
		runeRanges := make(gltext.RuneRanges, 0)
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 32, High: 128})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3000, High: 0x3030})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3040, High: 0x309f})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x30a0, High: 0x30ff})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x4e00, High: 0x9faf})
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 0xff00, High: 0xffef})

		scale := fixed.Int26_6(24)
		runesPerRow := fixed.Int26_6(128)
		config, err = gltext.NewTruetypeFontConfig(fd, scale, runeRanges, runesPerRow)
		if err != nil {
			panic(err)
		}
		config.Name = "font_1_honokamin"

		err = config.Save("fontconfigs")
		if err != nil {
			panic(err)
		}
		font, err = v41.NewFont(config)
		if err != nil {
			panic(err)
		}
	}

	width, height := window.GetSize()
	font.ResizeWindow(float32(width), float32(height))

	scaleMin, scaleMax := float32(1.0), float32(1.1)
	text := v41.NewText(font, scaleMin, scaleMax)
	str := "大好き"
	for _, s := range str {
		fmt.Printf("%c: %d\n", s, rune(s))
	}
	text.SetString(str)
	text.SetColor(mgl32.Vec3{1, 1, 1})
	text.FadeOutPerFrame = 0.01

	i := 0
	start := time.Now()

	gl.ClearColor(0.4, 0.4, 0.4, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Position can be set freely
		text.SetPosition(mgl32.Vec2{0, float32(i)})
		i++
		if i > 200 {
			i = -200
		}
		text.Draw()

		// just for illustrative purposes
		// i imagine that user interaction of some sort will trigger these rather than a moment in time

		// fade out
		if math.Floor(time.Now().Sub(start).Seconds()) == 5 {
			text.BeginFadeOut()
		}

		// show text
		if math.Floor(time.Now().Sub(start).Seconds()) == 10 {
			text.Show()
		}

		// hide
		if math.Floor(time.Now().Sub(start).Seconds()) == 15 {
			text.Hide()
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
	text.Release()
	font.Release()
}
