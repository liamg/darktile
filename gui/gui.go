package gui

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"unsafe"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kbinani/screenshot"
	"github.com/liamg/aminal/buffer"
	"github.com/liamg/aminal/config"
	"github.com/liamg/aminal/platform"
	"github.com/liamg/aminal/terminal"
	"github.com/liamg/aminal/version"
	"go.uber.org/zap"
)

type GUI struct {
	window            *glfw.Window
	logger            *zap.SugaredLogger
	config            *config.Config
	terminal          *terminal.Terminal
	width             int //window width in pixels
	height            int //window height in pixels
	appliedWidth      int
	appliedHeight     int
	resizeCache       *ResizeCache // resize cache formed by resizeToTerminal()
	dpiScale          float32
	fontMap           *FontMap
	fontScale         float32
	renderer          *OpenGLRenderer
	colourAttr        uint32
	mouseDown         bool
	mouseDownModifier glfw.ModifierKey
	overlay           overlay
	terminalAlpha     float32
	showDebugInfo     bool
	keyboardShortcuts map[config.UserAction]*config.KeyCombination
	resizeLock        *sync.Mutex
	handCursor        *glfw.Cursor
	arrowCursor       *glfw.Cursor
	defaultCell       *buffer.Cell

	prevLeftClickX                  uint16
	prevLeftClickY                  uint16
	prevMotionTX                    int
	prevMotionTY                    int
	leftClickTime                   time.Time
	leftClickCount                  int // number of clicks in a serie - single click, double click, or triple click
	mouseMovedAfterSelectionStarted bool
	internalResize                  bool
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type ResizeCache struct {
	Width  int
	Height int
	Cols   uint
	Rows   uint
}

func (g *GUI) GetMonitor() *glfw.Monitor {
	if g.window == nil {
		panic("to determine current monitor the window must be set")
	}
	monitors := glfw.GetMonitors()

	if len(monitors) == 1 {
		return glfw.GetPrimaryMonitor()
	}

	x, y := g.window.GetPos()
	w, h := g.window.GetSize()
	var currentMonitor *glfw.Monitor
	bestMatch := 0

	for _, monitor := range monitors {
		mode := monitor.GetVideoMode()
		mx, my := monitor.GetPos()
		overlap := Max(0, Min(x+w, mx+mode.Width)-Max(x, mx)) *
			Max(0, Min(y+h, my+mode.Height)-Max(y, my))
		if bestMatch < overlap {
			bestMatch = overlap
			currentMonitor = monitor
		}
	}

	if currentMonitor == nil {
		// Monitor couldn't be found (xrandr scaling?) - default to primary
		return glfw.GetPrimaryMonitor()
	}

	return currentMonitor
}

// SetDPIScale sets the GUI DPI scale from user configuration (if set)
// or by calculating it from the monitor's configuration and size.
func (g *GUI) SetDPIScale() {
	if g.config.DPIScale > 0 {
		g.dpiScale = g.config.DPIScale
		return
	}
	g.dpiScale = g.calculateDpiScale()
}

// calculateDpiScale determines the DPI scale in comparison with "standard" monitor's DPI values.
func (g *GUI) calculateDpiScale() float32 {
	const standardDpi = 96
	const mmPerInch = 25.4

	m := g.GetMonitor()
	widthMM, _ := m.GetPhysicalSize()
	if widthMM == 0 {
		return 1.0
	}

	monitorDpi := float32(m.GetVideoMode().Width) / (float32(widthMM) / mmPerInch)
	return monitorDpi / standardDpi
}

func New(config *config.Config, terminal *terminal.Terminal, logger *zap.SugaredLogger) (*GUI, error) {
	shortcuts, err := config.KeyMapping.GenerateActionMap()
	if err != nil {
		return nil, err
	}

	return &GUI{
		config:            config,
		logger:            logger,
		width:             800,
		height:            600,
		appliedWidth:      0,
		appliedHeight:     0,
		dpiScale:          1,
		terminal:          terminal,
		fontScale:         10.0,
		terminalAlpha:     1,
		keyboardShortcuts: shortcuts,
		resizeLock:        &sync.Mutex{},
		internalResize:    false,
	}, nil
}

// inspired by https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

func (gui *GUI) scale() float32 {
	pw, _ := gui.window.GetFramebufferSize()
	ww, _ := gui.window.GetSize()
	return float32(ww) / float32(pw)
}

// can only be called on OS thread
func (gui *GUI) resizeToTerminal(newCols uint, newRows uint) {

	if gui.window.GetAttrib(glfw.Iconified) != 0 {
		return
	}

	gui.resizeLock.Lock()
	defer gui.resizeLock.Unlock()

	cols, rows := gui.renderer.GetTermSize()
	if cols == newCols && rows == newRows {
		return
	}

	gui.logger.Debugf("Initiating GUI resize to columns=%d rows=%d", newCols, newRows)

	gui.logger.Debugf("Calculating size...")
	width, height := gui.renderer.GetRectangleSize(newCols, newRows)

	roundedWidth := int(math.Ceil(float64(width)))
	roundedHeight := int(math.Ceil(float64(height)))

	gui.resizeCache = &ResizeCache{roundedWidth, roundedHeight, newCols, newRows}

	gui.logger.Debugf("Resizing window to %dx%d", roundedWidth, roundedHeight)
	gui.internalResize = true
	gui.window.SetSize(roundedWidth, roundedHeight) // will trigger resize()
	gui.internalResize = false
}

func (gui *GUI) generateDefaultCell(reverse bool) {
	color := gui.config.ColourScheme.Background
	if reverse {
		color = gui.config.ColourScheme.Foreground
	}
	cell := buffer.NewBackgroundCell(color)
	gui.renderer.backgroundColour = color
	gui.defaultCell = &cell
	gl.ClearColor(
		color[0],
		color[1],
		color[2],
		1.0,
	)
}

func (gui *GUI) getCursorBg(cell *buffer.Cell) (bg [3]float32) {
	if gui.config.ColourScheme.Cursor != cell.Bg() {
		bg = gui.config.ColourScheme.Cursor
	} else {
		bg = cell.Fg()
	}
	return bg
}

func (gui *GUI) getCursorFg(cell *buffer.Cell) (fg [3]float32) {
	fg = cell.Bg()
	return fg
}

// can only be called on OS thread
func (gui *GUI) resize(w *glfw.Window, width int, height int) {

	if gui.window.GetAttrib(glfw.Iconified) != 0 {
		return
	}

	if gui.appliedWidth == width && gui.appliedHeight == height {
		return
	}

	if gui.internalResize == false {
		gui.resizeLock.Lock()
		defer gui.resizeLock.Unlock()
	}

	gui.logger.Debugf("Initiating GUI resize to %dx%d", width, height)

	gui.width = width
	gui.height = height
	gui.appliedWidth = width
	gui.appliedHeight = height

	gui.logger.Debugf("Updating font resolutions...")
	gui.loadFonts()

	gui.logger.Debugf("Setting renderer area...")
	gui.renderer.SetArea(0, 0, gui.width, gui.height)

	if gui.resizeCache != nil && gui.resizeCache.Width == width && gui.resizeCache.Height == height {
		gui.logger.Debugf("No need to resize internal terminal!")
	} else {
		gui.logger.Debugf("Calculating size in cols/rows...")
		cols, rows := gui.renderer.GetTermSize()
		gui.logger.Debugf("Resizing internal terminal...")
		if err := gui.terminal.SetSize(cols, rows); err != nil {
			gui.logger.Errorf("Failed to resize terminal to %d cols, %d rows: %s", cols, rows, err)
		}
	}

	gui.resizeCache = nil

	gui.logger.Debugf("Setting viewport size...")
	gl.Viewport(0, 0, int32(gui.width), int32(gui.height))

	gui.terminal.SetCharSize(gui.renderer.cellWidth, gui.renderer.cellHeight)

	gui.logger.Debugf("Resize complete!")

	gui.redraw()
	gui.window.SwapBuffers()
}

func (gui *GUI) getTermSize() (uint, uint) {
	if gui.renderer == nil {
		return 0, 0
	}
	return gui.renderer.GetTermSize()
}

func (gui *GUI) Close() {
	gui.window.SetShouldClose(true)
}

func (gui *GUI) Render() error {

	gui.logger.Debugf("Creating window...")
	var err error
	gui.window, err = gui.createWindow()
	gui.SetDPIScale()
	gui.window.SetSize(int(float32(gui.width)*gui.dpiScale),
		int(float32(gui.height)*gui.dpiScale))
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
	if err := gui.loadFonts(); err != nil {
		return fmt.Errorf("Failed to load font: %s", err)
	}

	titleChan := make(chan bool, 1)
	resizeChan := make(chan bool, 1)
	reverseChan := make(chan bool, 1)

	gui.renderer = NewOpenGLRenderer(gui.config, gui.fontMap, 0, 0, gui.width, gui.height, gui.colourAttr, program)

	gui.window.SetFramebufferSizeCallback(gui.resize)
	gui.window.SetKeyCallback(gui.key)
	gui.window.SetCharCallback(gui.char)
	gui.window.SetScrollCallback(gui.glfwScrollCallback)
	gui.window.SetMouseButtonCallback(gui.mouseButtonCallback)
	gui.window.SetCursorPosCallback(gui.mouseMoveCallback)
	gui.window.SetRefreshCallback(func(w *glfw.Window) {
		gui.terminal.SetDirty()
	})
	gui.window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused {
			gui.terminal.SetDirty()
		}
	})
	gui.window.SetPosCallback(gui.windowPosChangeCallback)
	glfw.SetMonitorCallback(gui.monitorChangeCallback)

	gui.generateDefaultCell(false)

	{
		w, h := gui.window.GetFramebufferSize()
		gui.resize(gui.window, w, h)
	}

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
	gl.Disable(gl.DEPTH_TEST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gui.terminal.AttachTitleChangeHandler(titleChan)
	gui.terminal.AttachResizeHandler(resizeChan)
	gui.terminal.AttachReverseHandler(reverseChan)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for {
			<-ticker.C
			gui.logger.Sync()
		}
	}()

	gui.terminal.SetProgram(program)

	latestVersion := ""

	go func() {
		r, err := version.GetNewerRelease()
		if err == nil && r != nil {
			latestVersion = r.TagName
			gui.terminal.SetDirty()
		}
	}()

	startTime := time.Now()
	showMessage := true

	for !gui.window.ShouldClose() {

		forceRedraw := false

		select {
		case <-titleChan:
			gui.window.SetTitle(gui.terminal.GetTitle())
		case <-resizeChan:
			cols, rows := gui.terminal.GetSize()
			gui.resizeToTerminal(uint(cols), uint(rows))
		case reverse := <-reverseChan:
			gui.generateDefaultCell(reverse)
			forceRedraw = true
		default:
			// this is more efficient than glfw.PollEvents()
			glfw.WaitEventsTimeout(0.02) // up to 50fps on no input, otherwise higher
		}

		if gui.terminal.CheckDirty() || forceRedraw {

			gui.redraw()

			if gui.showDebugInfo {
				gui.textbox(2, 2, fmt.Sprintf(`Cursor:      %d,%d
View Size:   %d,%d
Buffer Size: %d lines
`,
					gui.terminal.GetLogicalCursorX(),
					gui.terminal.GetLogicalCursorY(),
					gui.terminal.ActiveBuffer().ViewWidth(),
					gui.terminal.ActiveBuffer().ViewHeight(),
					gui.terminal.ActiveBuffer().Height(),
				),
					[3]float32{1, 1, 1},
					[3]float32{0.8, 0, 0},
				)
			}

			if showMessage {
				if latestVersion != "" && time.Since(startTime) < time.Second*10 && gui.terminal.ActiveBuffer().RawLine() == 0 {
					time.AfterFunc(time.Second, gui.terminal.SetDirty)
					_, h := gui.terminal.GetSize()
					var msg string
					if version.Version == "" {
						msg = "You are using a development build of Aminal."
					} else {
						msg = fmt.Sprintf("Version %s of Aminal is now available.", strings.Replace(latestVersion, "v", "", -1))
					}
					gui.textbox(
						2,
						uint16(h-3),
						fmt.Sprintf("%s (%d)", msg, 10-int(time.Since(startTime).Seconds())),
						[3]float32{1, 1, 1},
						[3]float32{0, 0.5, 0},
					)
				} else {
					showMessage = false
				}
			}

			gui.SwapBuffers()
		}

	}

	gui.logger.Debugf("Stopping render...")
	return nil

}

func (gui *GUI) redraw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	lines := gui.terminal.GetVisibleLines()
	lineCount := int(gui.terminal.ActiveBuffer().ViewHeight())
	colCount := int(gui.terminal.ActiveBuffer().ViewWidth())
	cx := uint(gui.terminal.GetLogicalCursorX())
	cy := uint(gui.terminal.GetLogicalCursorY()) + uint(gui.terminal.GetScrollOffset())
	var colour *config.Colour
	for y := 0; y < lineCount; y++ {
		if y < len(lines) {
			cells := lines[y].Cells()
			for x := 0; x < colCount; x++ {

				cursor := false
				if gui.terminal.Modes().ShowCursor {
					cursor = cx == uint(x) && cy == uint(y)
				}

				if gui.terminal.ActiveBuffer().InSelection(uint16(x), uint16(y)) {
					colour = &gui.config.ColourScheme.Selection
				} else {
					colour = nil
				}

				cell := gui.defaultCell
				if colour != nil || cursor || x < len(cells) {

					if x < len(cells) {
						cell = &cells[x]
						if cell.Image() != nil {
							gui.renderer.DrawCellImage(*cell, uint(x), uint(y))
							continue
						}
					}

					if cursor {
						var bgColour config.Colour = gui.getCursorBg(cell)
						colour = &bgColour
					}

					gui.renderer.DrawCellBg(*cell, uint(x), uint(y), colour, false)
				}

			}
		}
	}
	for y := 0; y < lineCount; y++ {

		if y < len(lines) {

			var builder strings.Builder
			bold := false
			dim := false
			col := 0
			colour := [3]float32{0, 0, 0}
			cells := lines[y].Cells()

			for x := 0; x < colCount; x++ {
				if x < len(cells) {
					cell := cells[x]

					cursor := false
					if gui.terminal.Modes().ShowCursor {
						cursor = cx == uint(x) && cy == uint(y)
					}

					var newFg [3]float32
					if cursor {
						newFg = gui.getCursorFg(&cell)
					} else {
						newFg = cell.Fg()
					}

					if builder.Len() > 0 && (cell.Attr().Dim != dim || cell.Attr().Bold != bold || colour != newFg) {
						var alpha float32 = 1.0
						if dim {
							alpha = 0.5
						}
						gui.renderer.DrawCellText(builder.String(), uint(col), uint(y), alpha, colour, bold)
						col = x
						builder.Reset()
					}
					dim = cell.Attr().Dim
					colour = newFg
					bold = cell.Attr().Bold
					r := cell.Rune()
					if r == 0 {
						r = ' '
					}
					builder.WriteRune(r)
				}
			}
			if builder.Len() > 0 {
				var alpha float32 = 1.0
				if dim {
					alpha = 0.5
				}
				gui.renderer.DrawCellText(builder.String(), uint(col), uint(y), alpha, colour, bold)
			}
		}

	}
	// underlines
	for y := 0; y < lineCount; y++ {

		if y < len(lines) {

			span := 0
			colour := [3]float32{0, 0, 0}
			cells := lines[y].Cells()

			var x int

			for x = 0; x < colCount && x < len(cells); x++ {
				cell := cells[x]
				if span > 0 && (!cell.Attr().Underline || colour != cell.Fg()) {
					gui.renderer.DrawUnderline(span, uint(x-span), uint(y), colour)
					span = 0
				}

				colour = cell.Fg()
				if cell.Attr().Underline {
					span++
				}
			}
			if span > 0 {
				gui.renderer.DrawUnderline(span, uint(x-span), uint(y), colour)
			}
		}

	}
	gui.renderOverlay()
}

func (gui *GUI) createWindow() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialise GLFW: %s", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	versions := [][2]int{
		{4, 6},
		{4, 5},
		{4, 4},
		{4, 3},
		{4, 2},
		{4, 1},
		{4, 0},
		{3, 3},
		{3, 2},
	}

	var window *glfw.Window

	for _, v := range versions {
		var err error
		window, err = gui.createWindowWithOpenGLVersion(v[0], v[1])
		if err == nil {
			break
		}
	}

	if window == nil {
		return nil, fmt.Errorf("failed to create window, please update your graphics drivers and try again")
	}

	window.SetSizeLimits(int(300*gui.dpiScale), int(150*gui.dpiScale), 10000, 10000)
	window.MakeContextCurrent()
	window.Show()
	window.Focus()

	return window, nil
}

func (gui *GUI) createWindowWithOpenGLVersion(major int, minor int) (*glfw.Window, error) {

	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)

	window, err := glfw.CreateWindow(int(float32(gui.width)*gui.dpiScale),
		int(float32(gui.height)*gui.dpiScale), "Aminal", nil, nil)
	if err != nil {
		e := err.Error()
		if i := strings.Index(e, ", got version "); i > -1 {
			v := strings.Split(strings.TrimSpace(e[i+14:]), ".")
			if len(v) == 2 {
				maj, mjErr := strconv.Atoi(v[0])
				if mjErr == nil {
					if min, miErr := strconv.Atoi(v[1]); miErr == nil {
						return gui.createWindowWithOpenGLVersion(maj, min)
					}
				}
			}
		}

		return nil, fmt.Errorf("failed to create window using OpenGL v%d.%d: %s", major, minor, err)
	}

	return window, nil
}

func (gui *GUI) onDebugMessage(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	gui.logger.Infof("GL debug message: %s", message)
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func (gui *GUI) createProgram() (uint32, error) {
	if err := gl.Init(); err != nil {
		return 0, fmt.Errorf("failed to initialise OpenGL: %s", err)
	}
	gui.logger.Infof("OpenGL version %s", gl.GoStr(gl.GetString(gl.VERSION)))

	if gui.config.DebugMode {
		// This allows to catch some OpenGL errors
		gl.DebugMessageCallback(gui.onDebugMessage, nil)
		gl.Enable(gl.DEBUG_OUTPUT)
	}

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

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return prog, nil
}

func (gui *GUI) launchTarget(target string) {

	err := platform.LaunchTarget(target)
	if err != nil {
		gui.logger.Errorf("Failed to launch target %s: %s", target, err)
	}
}

func (gui *GUI) SwapBuffers() {
	UpdateNSGLContext(gui.window)
	gui.window.SwapBuffers()
}

func (gui *GUI) Screenshot(path string) {
	x, y := gui.window.GetPos()
	w, h := gui.window.GetSize()

	img, err := screenshot.CaptureRect(image.Rectangle{Min: image.Point{X: x, Y: y},
		Max: image.Point{X: x + w, Y: y + h}})
	if err != nil {
		panic(err)
	}
	file, _ := os.Create(path)
	defer file.Close()
	png.Encode(file, img)
}

func (gui *GUI) windowPosChangeCallback(w *glfw.Window, xpos int, ypos int) {
	gui.SetDPIScale()
}

func (gui *GUI) monitorChangeCallback(monitor *glfw.Monitor, event glfw.MonitorEvent) {
	gui.SetDPIScale()
}
