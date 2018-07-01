package gui

import (
	"fmt"
	"math"
	"os"
	"runtime"
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
	cells      [][]Cell
	cols       int
	rows       int
	capslock   bool
}

func New(config config.Config, terminal *terminal.Terminal, logger *zap.SugaredLogger) *GUI {

	//logger.
	return &GUI{
		config:   config,
		logger:   logger,
		width:    600,
		height:   300,
		terminal: terminal,
		cells:    [][]Cell{},
	}
}

// inspired by https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

// can only be called on OS thread
func (gui *GUI) resize(w *glfw.Window, width int, height int) {

	if width == gui.width && height == gui.height {
		return
	}

	gui.logger.Debugf("GUI resize to %dx%d", width, height)

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
	text.Release()

	gui.cols = int(math.Floor(float64(float32(width) / gui.charWidth)))
	gui.rows = int(math.Floor(float64(float32(height) / gui.charHeight)))

	if err := gui.terminal.SetSize(gui.cols, gui.rows); err != nil {
		gui.logger.Errorf("Failed to resize terminal to %d cols, %d rows: %s", gui.cols, gui.rows, err)
	}

	gui.createTexts()
}

func (gui *GUI) getTermSize() (int, int) {
	return gui.cols, gui.rows
}

// checks if the terminals cells have been updated, and updates the text objects if needed - only call on OS thread
func (gui *GUI) updateTexts() {

	// runtime.LockOSThread() ?

	//gui.logger.Debugf("Updating texts...")

	cols, rows := gui.getTermSize()

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {

			c, err := gui.terminal.GetCellAtPos(terminal.Position{Line: row, Col: col})

			if err != nil || c == nil {
				gui.cells[row][col].Hide()
				continue
			}

			if c.IsHidden() {

				gui.cells[row][col].Hide()

				// debug
				//gui.texts[row][col].SetColor(c.GetColourVec())
				//gui.texts[row][col].SetString("?")
				//gui.texts[row][col].Show()
				// end debug
				continue
			}

			gui.cells[row][col].SetColour(c.GetColour())
			gui.cells[row][col].SetRune(c.GetRune())
			gui.cells[row][col].Show()

		}
	}
}

// builds text objects - only call on OS thread
func (gui *GUI) createTexts() {

	cols, rows := gui.getTermSize()

	cells := [][]Cell{}
	for row := 0; row < rows; row++ {

		if len(cells) <= row {
			cells = append(cells, []Cell{})
		}
		for col := 0; col < cols; col++ {
			if len(cells[row]) <= col {

				x := ((float32(col) * gui.charWidth) - (float32(gui.width) / 2)) + (gui.charWidth / 2)
				y := -(((float32(row) * gui.charHeight) - (float32(gui.height) / 2)) + (gui.charHeight / 2))

				cells[row] = append(cells[row], NewCell(gui.font, x, y, gui.charWidth, gui.charHeight))
			}
		}
	}

	gui.cells = cells

	gui.updateTexts()
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

	gui.window.SetFramebufferSizeCallback(gui.resize)
	gui.window.SetKeyCallback(gui.key)
	gui.window.SetCharCallback(gui.char)
	w, h := gui.window.GetSize()
	gui.resize(gui.window, w, h)

	gl.Viewport(0, 0, int32(gui.width), int32(gui.height))

	gui.logger.Debugf("Starting pty read handling...")

	updateChan := make(chan bool, 1024)

	gui.terminal.OnUpdate(func() {
		updateChan <- true
	})
	go gui.terminal.Read()

	text := v41.NewText(gui.font, 1.0, 1.1)
	text.SetString("")
	text.SetColor(mgl32.Vec3{1, 0, 0})
	text.SetPosition(mgl32.Vec2{0, 0})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	//gl.Disable(gl.MULTISAMPLE)
	// stop smoothing fonts
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	updateRequired := false

	gui.logger.Debugf("Starting render...")

	gl.ClearColor(0.1, 0.1, 0.1, 1.0)

	for !gui.window.ShouldClose() {

		updateRequired = false

	CheckUpdate:
		for {
			select {
			case <-updateChan:
				updateRequired = true
			case <-ticker.C:
				text.SetString(fmt.Sprintf("%dx%d", gui.cols, gui.rows))
				updateRequired = true
			default:
				break CheckUpdate
			}
		}

		gl.UseProgram(program)

		if updateRequired {

			gui.updateTexts()

			// Render the string.
			gui.window.SetTitle(gui.terminal.GetTitle())

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			cols, rows := gui.getTermSize()

			for row := 0; row < rows; row++ {
				for col := 0; col < cols; col++ {
					gui.cells[row][col].Draw()
				}
			}

			text.Draw()
		}

		glfw.PollEvents()
		if updateRequired {
			gui.window.SwapBuffers()
		}
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
