package gui

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/4ydx/gltext"
	v41 "github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"gitlab.com/liamg/raft/config"
	"gitlab.com/liamg/raft/terminal"
	"go.uber.org/zap"
	"golang.org/x/image/math/fixed"
)

type GUI struct {
	window     *glfw.Window
	logger     *zap.SugaredLogger
	config     config.Config
	font       *v41.Font
	terminal   *terminal.Terminal
	width      int
	height     int
	charWidth  float32
	charHeight float32
	texts      [][]*v41.Text
	textLock   sync.Mutex
}

func New(config config.Config, terminal *terminal.Terminal, logger *zap.SugaredLogger) *GUI {

	//logger.
	return &GUI{
		config:   config,
		logger:   logger,
		width:    600,
		height:   300,
		terminal: terminal,
		texts:    [][]*v41.Text{},
	}
}

// inspired by https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

func (gui *GUI) SetSize(w int, h int) {
	gui.window.SetSize(w, h)
	gui.resize(gui.window, w, h)
}

func (gui *GUI) resize(w *glfw.Window, width int, height int) {
	gui.width = width
	gui.height = height
	if gui.font != nil {
		gui.font.ResizeWindow(float32(width), float32(height))
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	scaleMin, scaleMax := float32(1.0), float32(1.1)
	text := v41.NewText(gui.font, scaleMin, scaleMax)
	text.SetString("A")
	gui.charWidth, gui.charHeight = text.Width(), text.Height()

	cols := int(math.Floor(float64(float32(width) / gui.charWidth)))
	rows := int(math.Floor(float64(float32(height) / gui.charHeight)))

	if err := gui.terminal.SetSize(cols, rows); err != nil {
		gui.logger.Errorf("Failed to resize terminal to %d cols, %d rows: %s", cols, rows, err)
	}

	gui.createTexts()
}

func (gui *GUI) getTermSize() (int, int) {
	scaleMin, scaleMax := float32(1.0), float32(1.1)
	text := v41.NewText(gui.font, scaleMin, scaleMax)
	text.SetString("A")
	return gui.width / int(text.Width()), gui.height / int(text.Height())
}

// checks if the terminals cells have been updated, and updates the text objects if needed
func (gui *GUI) updateTexts() {
	gui.textLock.Lock()
	defer gui.textLock.Unlock()

	cols, rows := gui.getTermSize()

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {

			r, err := gui.terminal.GetRuneAtPos(terminal.Position{Row: row, Col: col})
			if err != nil {
				gui.logger.Errorf("Failed to read rune: %s", err)
			}
			if r > 0 {
				gui.texts[row][col].SetString(string(r))
				gui.texts[row][col].SetColor(mgl32.Vec3{1, 1, 1})
				// @todo set colour
			}
		}
	}
}

// builds text objects
func (gui *GUI) createTexts() {
	gui.textLock.Lock()
	defer gui.textLock.Unlock()
	scaleMin, scaleMax := float32(1.0), float32(1.1)

	cols, rows := gui.getTermSize()

	texts := [][]*v41.Text{}
	for row := 0; row < rows; row++ {

		if len(texts) <= row {
			texts = append(texts, []*v41.Text{})
		}
		for col := 0; col < cols; col++ {
			if len(texts[row]) <= col {
				text := v41.NewText(gui.font, scaleMin, scaleMax)

				if row < len(gui.texts) {
					if col < len(gui.texts[row]) {
						text.SetString(gui.texts[row][col].String)
					}
				}

				text.SetColor(mgl32.Vec3{1, 1, 1})

				x := ((float32(col) * gui.charWidth) - (float32(gui.width) / 2)) + (gui.charWidth / 2)
				y := -(((float32(row) * gui.charHeight) - (float32(gui.height) / 2)) + (gui.charHeight / 2))

				text.SetPosition(mgl32.Vec2{x, y})
				texts[row] = append(texts[row], text)
			}
		}
	}

	gui.texts = texts

}

func (gui *GUI) Render() error {

	gui.logger.Debugf("Locking OS thread...")
	runtime.LockOSThread()

	gui.logger.Debugf("Creating window...")
	var err error
	gui.window, err = gui.createWindow(gui.width, gui.height)
	if err != nil {
		return fmt.Errorf("Failed to create window: %s", err)
	}
	defer glfw.Terminate()

	gui.logger.Debugf("Initialising OpenGL and creating program...")
	program, err := gui.createProgram()
	if err != nil {
		return fmt.Errorf("Failed to initialise OpenGL: %s", err)
	}

	gui.logger.Debugf("Loading font...")
	//gui.font, err = gui.loadFont("/usr/share/fonts/nerd-fonts-complete/ttf/Roboto Mono Nerd Font Complete.ttf", 12)
	if err := gui.loadFont("./fonts/CamingoCode-Regular.ttf", 12); err != nil {
		return fmt.Errorf("Failed to load font: %s", err)
	}

	gui.resize(gui.window, gui.width, gui.height)

	gui.window.SetFramebufferSizeCallback(gui.resize)
	w, h := gui.window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))

	gui.logger.Debugf("Starting pty read handling...")
	gui.terminal.OnUpdate(func() {
		gui.updateTexts()
	})
	go gui.terminal.Read()

	scaleMin, scaleMax := float32(1.0), float32(1.1)
	text := v41.NewText(gui.font, scaleMin, scaleMax)
	text.SetString("")
	text.SetColor(mgl32.Vec3{1, 0, 0})
	text.SetPosition(mgl32.Vec2{0, 0})

	frames := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	//gl.Disable(gl.MULTISAMPLE)
	// stop smoothing fonts
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gui.logger.Debugf("Starting render...")
	for !gui.window.ShouldClose() {

		select {
		case <-ticker.C:
			text.SetString(fmt.Sprintf("%d fps | %d, %d", frames, gui.terminal.GetPosition().Row, gui.terminal.GetPosition().Col))
			frames = 0
		default:
		}

		gl.UseProgram(program)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render the string.
		gui.window.SetTitle(gui.terminal.GetTitle())

		gui.textLock.Lock()
		cols, rows := gui.getTermSize()

		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				gui.texts[row][col].Draw()
			}
		}
		gui.textLock.Unlock()

		text.Draw()

		glfw.PollEvents()
		gui.window.SwapBuffers()

		frames++
	}

	gui.logger.Debugf("Stopping render...")
	return nil

}

func (gui *GUI) loadFont(path string, scale int32) error {

	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	runeRanges := make(gltext.RuneRanges, 0)
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 32, High: 127})
	/*runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3000, High: 0x3030})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3040, High: 0x309f})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x30a0, High: 0x30ff})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x4e00, High: 0x9faf})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0xff00, High: 0xffef})
	*/

	runesPerRow := fixed.Int26_6(128)
	conf, err := gltext.NewTruetypeFontConfig(fd, fixed.Int26_6(scale), runeRanges, runesPerRow)
	if err != nil {
		return err
	}

	font, err := v41.NewFont(conf)
	if err != nil {
		return err
	}
	font.ResizeWindow(float32(gui.width), float32(gui.height))
	gui.font = font
	return nil
}

func (gui *GUI) createWindow(width int, height int) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Terminal", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	return window, nil
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func (gui *GUI) createProgram() (uint32, error) {
	if err := gl.Init(); err != nil {
		return 0, fmt.Errorf("Failed to initialise OpenGL: %s", err)
	}
	gui.logger.Infof("OpenGL version %s", gl.GoStr(gl.GetString(gl.VERSION)))

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog, nil
}
