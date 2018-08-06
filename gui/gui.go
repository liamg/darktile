package gui

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/liamg/glfont"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"gitlab.com/liamg/raft/config"
	"gitlab.com/liamg/raft/terminal"
	"go.uber.org/zap"
)

type GUI struct {
	window          *glfw.Window
	logger          *zap.SugaredLogger
	config          config.Config
	terminal        *terminal.Terminal
	width           int //window width in pixels
	height          int //window height in pixels
	charWidth       float32
	charHeight      float32
	cells           [][]Cell
	cols            int
	rows            int
	colourAttr      uint32
	font            *glfont.Font
	fontScale       int32
	verticalPadding float32
}

func New(config config.Config, terminal *terminal.Terminal, logger *zap.SugaredLogger) *GUI {

	//logger.
	return &GUI{
		config:          config,
		logger:          logger,
		width:           600,
		height:          300,
		terminal:        terminal,
		cells:           [][]Cell{},
		fontScale:       12.0,
		verticalPadding: 6.0,
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
		gui.font.UpdateResolution((width), (height))
	}

	gl.Viewport(0, 0, int32(gui.width), int32(gui.height))

	gui.charWidth, gui.charHeight = gui.font.Width(1, "A"), gui.font.Height(1, "A")+gui.verticalPadding

	gui.cols = int(math.Floor(float64(float32(width) / gui.charWidth)))
	gui.rows = int(math.Floor(float64(float32(height) / gui.charHeight)))

	//fmt.Printf("%#v\n", gui)
	//os.Exit(0)

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

			if err != nil {
				//gui.logger.Errorf("Failed to get cell: %s", err)
				gui.cells[row][col].Hide()
				continue
			}

			if c == nil || c.IsHidden() {
				gui.cells[row][col].Hide()
				continue
			}

			gui.cells[row][col].SetFgColour(c.GetFgColour())
			gui.cells[row][col].SetRune(c.GetRune())
			gui.cells[row][col].Show()

			if gui.terminal.IsCursorVisible() && gui.terminal.GetPosition().Col == col && gui.terminal.GetPosition().Line == row {
				gui.cells[row][col].SetBgColour(
					gui.config.ColourScheme.Cursor[0],
					gui.config.ColourScheme.Cursor[1],
					gui.config.ColourScheme.Cursor[2],
				)
			} else {
				gui.cells[row][col].SetBgColour(c.GetBgColour())
			}
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

				cells[row] = append(cells[row], gui.NewCell(gui.font, x, y, gui.charWidth, gui.charHeight, gui.colourAttr, gui.config.ColourScheme.DefaultBg))
			}
		}
	}

	gui.cells = cells

	gui.updateTexts()
}

func (gui *GUI) Close() {
	gui.window.SetShouldClose(true)
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

	gui.colourAttr = uint32(gl.GetAttribLocation(program, gl.Str("inColour\x00")))
	gl.BindFragDataLocation(program, 0, gl.Str("outColour\x00"))

	gui.logger.Debugf("Loading font...")
	//if err := gui.loadFont("/usr/share/fonts/nerd-fonts-complete/ttf/Roboto Mono Nerd Font Complete.ttf", 12); err != nil {
	if err := gui.loadFont("./fonts/Roboto.ttf"); err != nil {
		return fmt.Errorf("Failed to load font: %s", err)
	}

	gui.window.SetFramebufferSizeCallback(gui.resize)
	gui.window.SetKeyCallback(gui.key)
	gui.window.SetCharCallback(gui.char)
	w, h := gui.window.GetSize()
	gui.resize(gui.window, w, h)

	gui.logger.Debugf("Starting pty read handling...")

	updateChan := make(chan bool, 1024)

	gui.terminal.OnUpdate(func() {
		updateChan <- true
	})
	go func() {
		err := gui.terminal.Read()
		if err != nil {
			gui.logger.Errorf("Read from pty failed: %s", err)
		}
		gui.Close()
	}()

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	//gl.Disable(gl.MULTISAMPLE)
	// stop smoothing fonts
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	updateRequired := 0

	gui.logger.Debugf("Starting render...")

	gl.UseProgram(program)

	gl.ClearColor(
		gui.config.ColourScheme.DefaultBg[0],
		gui.config.ColourScheme.DefaultBg[1],
		gui.config.ColourScheme.DefaultBg[2],
		1.0,
	)

	for !gui.window.ShouldClose() {

		if updateRequired > 0 {

			updateRequired--

		} else {
		CheckUpdate:
			for {
				select {
				case <-updateChan:
					updateRequired = 2
				case <-ticker.C:
					/*
						text.SetString(
							fmt.Sprintf(
								"%dx%d@%d,%d",
								gui.cols,
								gui.rows,
								gui.terminal.GetPosition().Col,
								gui.terminal.GetPosition().Line,
							),
						)*/
					updateRequired = 2
				default:
					break CheckUpdate
				}
			}
		}

		gl.UseProgram(program)

		if updateRequired > 0 {

			gui.updateTexts()

			// Render the string.
			// @todo uncommentbut dont run all of the time... - perhaps use onTitleChange event from terminal?
			//gui.window.SetTitle(gui.terminal.GetTitle())

			//gl.ClearColor(0.5, 0.5, 0.5, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			cols, rows := gui.getTermSize()

			for row := 0; row < rows; row++ {
				for col := 0; col < cols; col++ {
					gui.cells[row][col].DrawBg()
				}
			}

			for row := 0; row < rows; row++ {
				for col := 0; col < cols; col++ {
					gui.cells[row][col].DrawText()
				}
			}

			gui.font.SetColor(1, 0.2, 0.2, 0.5)
			gui.font.Printf(100, 300, 1.5, "%#v", gui.terminal.GetPosition())

		}

		glfw.PollEvents()
		if updateRequired > 0 {
			gui.window.SwapBuffers()
		}
	}

	gui.logger.Debugf("Stopping render...")
	return nil

}

func (gui *GUI) loadFont(path string) error {
	font, err := glfont.LoadFont("./fonts/Roboto.ttf", gui.fontScale, gui.width, gui.height)
	if err != nil {
		return fmt.Errorf("LoadFont: %v", err)
	}
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

	gui.logger.Debugf("Compiling shaders...")

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog, nil
}
