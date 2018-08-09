package gui

import (
	"fmt"
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
	window     *glfw.Window
	logger     *zap.SugaredLogger
	config     config.Config
	terminal   *terminal.Terminal
	width      int //window width in pixels
	height     int //window height in pixels
	font       *glfont.Font
	fontScale  int32
	renderer   Renderer
	colourAttr uint32
}

func New(config config.Config, terminal *terminal.Terminal, logger *zap.SugaredLogger) *GUI {

	//logger.
	return &GUI{
		config:    config,
		logger:    logger,
		width:     600,
		height:    300,
		terminal:  terminal,
		fontScale: 15.0,
	}
}

// inspired by https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

// can only be called on OS thread
func (gui *GUI) resize(w *glfw.Window, width int, height int) {

	gui.logger.Debugf("GUI resize to %dx%d", width, height)

	gui.width = width
	gui.height = height

	if gui.font != nil {
		gui.font.UpdateResolution((width), (height))
	}

	gl.Viewport(0, 0, int32(gui.width), int32(gui.height))

	gui.renderer.SetArea(0, 0, gui.width, gui.height)

	cols, rows := gui.renderer.GetTermSize()

	if err := gui.terminal.SetSize(cols, rows); err != nil {
		gui.logger.Errorf("Failed to resize terminal to %d cols, %d rows: %s", cols, rows, err)
	}

}

func (gui *GUI) getTermSize() (int, int) {
	if gui.renderer == nil {
		return 0, 0
	}
	return gui.renderer.GetTermSize()
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
	if err := gui.loadFont("./fonts/envypn-15.ttf"); err != nil {
		//if err := gui.loadFont("./fonts/Roboto.ttf"); err != nil {
		return fmt.Errorf("Failed to load font: %s", err)
	}

	changeChan := make(chan bool, 1)
	titleChan := make(chan bool, 1)

	gui.renderer = NewOpenGLRenderer(gui.config, gui.font, gui.fontScale, 0, 0, gui.width, gui.height, gui.colourAttr, program)

	gui.window.SetFramebufferSizeCallback(gui.resize)
	gui.window.SetKeyCallback(gui.key)
	gui.window.SetCharCallback(gui.char)
	gui.window.SetRefreshCallback(func(w *glfw.Window) {
		select {
		case changeChan <- true:
		default:
		}
	})
	gui.window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused {
			select {
			case changeChan <- true:
			default:
			}
		}
	})
	w, h := gui.window.GetSize()
	gui.resize(gui.window, w, h)

	gui.logger.Debugf("Starting pty read handling...")

	go func() {
		err := gui.terminal.Read()
		if err != nil {
			gui.logger.Errorf("Read from pty failed: %s", err)
		}
		gui.Close()
	}()

	gui.logger.Debugf("Starting render...")

	gl.UseProgram(program)

	// stop smoothing fonts

	//gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	//glfw.SwapInterval(1)

	gl.ClearColor(
		gui.config.ColourScheme.Background[0],
		gui.config.ColourScheme.Background[1],
		gui.config.ColourScheme.Background[2],
		1.0,
	)

	gui.terminal.AttachTitleChangeHandler(titleChan)
	gui.terminal.AttachDisplayChangeHandler(changeChan)

	frames := 0
	frameCount := 0
	fps := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for !gui.window.ShouldClose() {

		select {

		case <-changeChan:
			frames = 2
			gui.logger.Sync()
		case <-titleChan:
			gui.window.SetTitle(gui.terminal.GetTitle())
		case <-ticker.C:
			fps = frameCount
			frameCount = 0
		default:
		}

		gl.UseProgram(program)

		if gui.config.Rendering.AlwaysRepaint || frames > 0 {

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			cols, rows := gui.getTermSize()

			for row := 0; row < rows; row++ {
				for col := 0; col < cols; col++ {
					gui.renderer.DrawCell(gui.terminal.GetCell(col, row), col, row)
				}
			}

			cx := int(gui.terminal.GetLogicalCursorX())
			cy := int(gui.terminal.GetLogicalCursorY())
			gui.renderer.DrawCursor(cx, cy, gui.config.ColourScheme.Cursor)

			_ = fps
			/*
				gui.font.SetColor(1, 0.5, 0.5, 0.5)
				fpsData := ""
				if gui.config.Rendering.AlwaysRepaint {
					fpsData = fmt.Sprintf("%d FPS | %d,%d", fps, gui.terminal.GetLogicalCursorX(), gui.terminal.GetLogicalCursorY())
				}
				gui.font.Print(10, float32(gui.height-20), 1.5, fmt.Sprintf("%s", fpsData))
			*/
		}

		if gui.config.Rendering.AlwaysRepaint || frames > 0 {
			gui.window.SwapBuffers()
			frameCount++
			frames--
		}

		glfw.WaitEventsTimeout(0.02) // up to 50fps on no input, otherwise higher
	}

	gui.logger.Debugf("Stopping render...")
	return nil

}

func (gui *GUI) loadFont(path string) error {
	font, err := glfont.LoadFont(path, gui.fontScale, gui.width, gui.height)
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
	glfw.WindowHint(glfw.ContextVersionMajor, 3) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
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
