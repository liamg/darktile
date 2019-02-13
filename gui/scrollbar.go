package gui

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	scrollbarVertexShaderSource = `
		#version 330 core
		layout (location = 0) in vec2 position;
		uniform vec2 resolution;

		void main() {
			// convert from window coordinates to GL coordinates
			vec2 glCoordinates = ((position / resolution) * 2.0 - 1.0) * vec2(1, -1);

			gl_Position = vec4(glCoordinates, 0.0, 1.0);
		}` + "\x00"

	scrollbarFragmentShaderSource = `
		#version 330 core
		uniform vec4 inColor;
		out vec4 outColor;
		void main() {
			outColor = inColor;
		}` + "\x00"

	BorderVertexValuesCount = 16
	ArrowsVertexValuesCount = 24
)

var (
	scrollbarColor_Bg           = [3]float32{float32(241) / float32(255), float32(241) / float32(255), float32(241) / float32(255)}
	scrollbarColor_ThumbNormal  = [3]float32{float32(193) / float32(255), float32(193) / float32(255), float32(193) / float32(255)}
	scrollbarColor_ThumbHover   = [3]float32{float32(168) / float32(255), float32(168) / float32(255), float32(168) / float32(255)}
	scrollbarColor_ThumbClicked = [3]float32{float32(120) / float32(255), float32(120) / float32(255), float32(120) / float32(255)}

	scrollbarColor_ButtonNormalBg = [3]float32{float32(241) / float32(255), float32(241) / float32(255), float32(241) / float32(255)}
	scrollbarColor_ButtonNormalFg = [3]float32{float32(80) / float32(255), float32(80) / float32(255), float32(80) / float32(255)}

	scrollbarColor_ButtonHoverBg = [3]float32{float32(210) / float32(255), float32(210) / float32(255), float32(210) / float32(255)}
	scrollbarColor_ButtonHoverFg = [3]float32{float32(80) / float32(255), float32(80) / float32(255), float32(80) / float32(255)}

	scrollbarColor_ButtonDisabledBg = [3]float32{float32(241) / float32(255), float32(241) / float32(255), float32(241) / float32(255)}
	scrollbarColor_ButtonDisabledFg = [3]float32{float32(163) / float32(255), float32(163) / float32(255), float32(163) / float32(255)}

	scrollbarColor_ButtonClickedBg = [3]float32{float32(120) / float32(255), float32(120) / float32(255), float32(120) / float32(255)}
	scrollbarColor_ButtonClickedFg = [3]float32{float32(255) / float32(255), float32(255) / float32(255), float32(255) / float32(255)}
)

type scrollbarPart int

const (
	None scrollbarPart = iota
	UpperArrow
	UpperSpace // the space between upper arrow and thumb
	Thumb
	BottomSpace // the space between thumb and bottom arrow
	BottomArrow
)

type ScreenRectangle struct {
	left, top     float32 // upper left corner in pixels relative to the window (in pixels)
	right, bottom float32
}

func (sr *ScreenRectangle) width() float32 {
	return sr.right - sr.left
}

func (sr *ScreenRectangle) height() float32 {
	return sr.bottom - sr.top
}

func (sr *ScreenRectangle) isInside(x float32, y float32) bool {
	return x >= sr.left && x < sr.right &&
		y >= sr.top && y < sr.bottom
}

type scrollbar struct {
	program                   uint32
	vbo                       uint32
	vao                       uint32
	uniformLocationResolution int32
	uniformLocationInColor    int32

	isDirty bool

	position            ScreenRectangle // relative to the window's top left corner, in pixels
	positionUpperArrow  ScreenRectangle // relative to the control's top left corner
	positionBottomArrow ScreenRectangle
	positionThumb       ScreenRectangle

	scrollPosition    int
	maxScrollPosition int

	thumbIsDragging           bool
	startedDraggingAtPosition int     // scrollPosition when the dragging was started
	startedDraggingAtThumbTop float32 // sb.positionThumb.top when the dragging was started
	offsetInThumbY            float32 // y offset inside the thumb of the dragging point
	scrollPositionDelta       int

	upperArrowIsDown  bool
	bottomArrowIsDown bool

	upperArrowFg  []float32
	upperArrowBg  []float32
	bottomArrowFg []float32
	bottomArrowBg []float32
	thumbColor    []float32
}

// Returns the vertical scrollbar width in pixels
func getDefaultScrollbarWidth() int {
	return 13
}

func createScrollbarProgram() (uint32, error) {
	vertexShader, err := compileShader(scrollbarVertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(scrollbarFragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog, nil
}

func newScrollbar() (*scrollbar, error) {
	prog, err := createScrollbarProgram()
	if err != nil {
		return nil, err
	}

	var vbo uint32
	var vao uint32

	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, (BorderVertexValuesCount+ArrowsVertexValuesCount)*4, nil, gl.DYNAMIC_DRAW) // only reserve data

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	result := &scrollbar{
		program:                   prog,
		vbo:                       vbo,
		vao:                       vao,
		uniformLocationResolution: gl.GetUniformLocation(prog, gl.Str("resolution\x00")),
		uniformLocationInColor:    gl.GetUniformLocation(prog, gl.Str("inColor\x00")),

		isDirty: false,

		position: ScreenRectangle{
			right:  0,
			bottom: 0,
			left:   0,
			top:    0,
		},

		scrollPosition:    0,
		maxScrollPosition: 0,

		thumbIsDragging:   false,
		upperArrowIsDown:  false,
		bottomArrowIsDown: false,
	}

	result.recalcElementPositions()
	result.resetElementColors(-1, -1) // (-1, -1) ensures that no part is hovered by the mouse

	return result, nil
}

func (sb *scrollbar) Free() {
	if sb.program != 0 {
		gl.DeleteProgram(sb.program)
		sb.program = 0
	}

	if sb.vbo != 0 {
		gl.DeleteBuffers(1, &sb.vbo)
		sb.vbo = 0
	}

	if sb.vao != 0 {
		gl.DeleteBuffers(1, &sb.vao)
		sb.vao = 0
	}
}

// Recalc positions of the scrollbar elements according to current
func (sb *scrollbar) recalcElementPositions() {
	arrowHeight := sb.position.width()

	sb.positionUpperArrow = ScreenRectangle{
		left:   0,
		top:    0,
		right:  sb.position.width(),
		bottom: arrowHeight,
	}

	sb.positionBottomArrow = ScreenRectangle{
		left:   sb.positionUpperArrow.left,
		top:    sb.position.height() - arrowHeight,
		right:  sb.positionUpperArrow.right,
		bottom: sb.position.height(),
	}
	thumbHeight := sb.position.width()
	thumbTop := arrowHeight
	if sb.maxScrollPosition != 0 {
		thumbTop += (float32(sb.scrollPosition) * (sb.position.height() - thumbHeight - arrowHeight*2)) / float32(sb.maxScrollPosition)
	}

	sb.positionThumb = ScreenRectangle{
		left:   2,
		top:    thumbTop,
		right:  sb.position.width() - 2,
		bottom: thumbTop + thumbHeight,
	}
}

func (sb *scrollbar) resize(gui *GUI) {
	sb.position.left = float32(gui.width) - float32(getDefaultScrollbarWidth())*gui.dpiScale
	sb.position.top = float32(0.0)
	sb.position.right = float32(gui.width)
	sb.position.bottom = float32(gui.height - 1)

	sb.recalcElementPositions()
	sb.isDirty = true
}

func (sb *scrollbar) render(gui *GUI) {
	var savedProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &savedProgram)
	defer gl.UseProgram(uint32(savedProgram))

	gl.UseProgram(sb.program)
	gl.Uniform2f(sb.uniformLocationResolution, float32(gui.width), float32(gui.height))
	gl.BindVertexArray(sb.vao)
	defer gl.BindVertexArray(0)

	// Draw background
	gl.Uniform4f(sb.uniformLocationInColor, scrollbarColor_Bg[0], scrollbarColor_Bg[1], scrollbarColor_Bg[2], 1.0)
	borderVertices := [...]float32{
		sb.position.left, sb.position.top,
		sb.position.right, sb.position.top,
		sb.position.right, sb.position.bottom,

		sb.position.right, sb.position.bottom,
		sb.position.left, sb.position.bottom,
		sb.position.left, sb.position.top,
	}
	gl.NamedBufferSubData(sb.vbo, 0, len(borderVertices)*4, gl.Ptr(&borderVertices[0]))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(borderVertices)/2))

	// Draw upper arrow
	// Upper arrow background
	gl.Uniform4f(sb.uniformLocationInColor, sb.upperArrowBg[0], sb.upperArrowBg[1], sb.upperArrowBg[2], 1.0)
	upperArrowBgVertices := [...]float32{
		sb.position.left + sb.positionUpperArrow.left, sb.position.top + sb.positionUpperArrow.top,
		sb.position.left + sb.positionUpperArrow.right, sb.position.top + sb.positionUpperArrow.top,
		sb.position.left + sb.positionUpperArrow.right, sb.position.top + sb.positionUpperArrow.bottom,

		sb.position.left + sb.positionUpperArrow.right, sb.position.top + sb.positionUpperArrow.bottom,
		sb.position.left + sb.positionUpperArrow.left, sb.position.top + sb.positionUpperArrow.bottom,
		sb.position.left + sb.positionUpperArrow.left, sb.position.top + sb.positionUpperArrow.top,
	}
	gl.NamedBufferSubData(sb.vbo, 0, len(upperArrowBgVertices)*4, gl.Ptr(&upperArrowBgVertices[0]))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(upperArrowBgVertices)/2))

	// Upper arrow foreground
	gl.Uniform4f(sb.uniformLocationInColor, sb.upperArrowFg[0], sb.upperArrowFg[1], sb.upperArrowFg[2], 1.0)
	upperArrowFgVertices := [...]float32{
		sb.position.left + sb.positionUpperArrow.left + sb.positionUpperArrow.width()/2.0, sb.position.top + sb.positionUpperArrow.top + sb.positionUpperArrow.height()/3.0,
		sb.position.left + sb.positionUpperArrow.left + sb.positionUpperArrow.width()*2.0/3.0, sb.position.top + sb.positionUpperArrow.top + sb.positionUpperArrow.height()/2.0,
		sb.position.left + sb.positionUpperArrow.left + sb.positionUpperArrow.width()/3.0, sb.position.top + sb.positionUpperArrow.top + sb.positionUpperArrow.height()/2.0,
	}
	gl.NamedBufferSubData(sb.vbo, 0, len(upperArrowFgVertices)*4, gl.Ptr(&upperArrowFgVertices[0]))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(upperArrowFgVertices)/2))

	// Draw bottom arrow
	// Bottom arrow background
	gl.Uniform4f(sb.uniformLocationInColor, sb.bottomArrowBg[0], sb.bottomArrowBg[1], sb.bottomArrowBg[2], 1.0)
	bottomArrowBgVertices := [...]float32{
		sb.position.left + sb.positionBottomArrow.left, sb.position.top + sb.positionBottomArrow.top,
		sb.position.left + sb.positionBottomArrow.right, sb.position.top + sb.positionBottomArrow.top,
		sb.position.left + sb.positionBottomArrow.right, sb.position.top + sb.positionBottomArrow.bottom,

		sb.position.left + sb.positionBottomArrow.right, sb.position.top + sb.positionBottomArrow.bottom,
		sb.position.left + sb.positionBottomArrow.left, sb.position.top + sb.positionBottomArrow.bottom,
		sb.position.left + sb.positionBottomArrow.left, sb.position.top + sb.positionBottomArrow.top,
	}
	gl.NamedBufferSubData(sb.vbo, 0, len(bottomArrowBgVertices)*4, gl.Ptr(&bottomArrowBgVertices[0]))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(bottomArrowBgVertices)/2))

	// Bottom arrow foreground
	gl.Uniform4f(sb.uniformLocationInColor, sb.bottomArrowFg[0], sb.bottomArrowFg[1], sb.bottomArrowFg[2], 1.0)
	bottomArrowFgVertices := [...]float32{
		sb.position.left + sb.positionBottomArrow.left + sb.positionBottomArrow.width()/3.0, sb.position.top + sb.positionBottomArrow.top + sb.positionBottomArrow.height()/2.0,
		sb.position.left + sb.positionBottomArrow.left + sb.positionBottomArrow.width()*2.0/3.0, sb.position.top + sb.positionBottomArrow.top + sb.positionBottomArrow.height()/2.0,
		sb.position.left + sb.positionBottomArrow.left + sb.positionBottomArrow.width()/2.0, sb.position.top + sb.positionBottomArrow.top + sb.positionBottomArrow.height()*2.0/3.0,
	}
	gl.NamedBufferSubData(sb.vbo, 0, len(bottomArrowFgVertices)*4, gl.Ptr(&bottomArrowFgVertices[0]))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(bottomArrowFgVertices)/2))

	// Draw thumb
	gl.Uniform4f(sb.uniformLocationInColor, sb.thumbColor[0], sb.thumbColor[1], sb.thumbColor[2], 1.0)
	thumbVertices := [...]float32{
		sb.position.left + sb.positionThumb.left, sb.position.top + sb.positionThumb.top,
		sb.position.left + sb.positionThumb.right, sb.position.top + sb.positionThumb.top,
		sb.position.left + sb.positionThumb.right, sb.position.top + sb.positionThumb.bottom,

		sb.position.left + sb.positionThumb.right, sb.position.top + sb.positionThumb.bottom,
		sb.position.left + sb.positionThumb.left, sb.position.top + sb.positionThumb.bottom,
		sb.position.left + sb.positionThumb.left, sb.position.top + sb.positionThumb.top,
	}
	gl.NamedBufferSubData(sb.vbo, 0, len(thumbVertices)*4, gl.Ptr(&thumbVertices[0]))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(thumbVertices)/2))

	sb.isDirty = false
}

func (sb *scrollbar) setPosition(max int, position int) {
	if max <= 0 {
		max = position
	}

	if position > max {
		position = max
	}

	sb.maxScrollPosition = max
	sb.scrollPosition = position

	sb.recalcElementPositions()
	sb.isDirty = true
}

func (sb *scrollbar) mouseHitTest(px float64, py float64) scrollbarPart {
	// convert to local coordinates
	mouseX := float32(px - float64(sb.position.left))
	mouseY := float32(py - float64(sb.position.top))

	result := None

	if sb.positionUpperArrow.isInside(mouseX, mouseY) {
		result = UpperArrow
	} else if sb.positionBottomArrow.isInside(mouseX, mouseY) {
		result = BottomArrow
	} else if sb.positionThumb.isInside(mouseX, mouseY) {
		result = Thumb
	} else {
		// construct UpperSpace
		pos := ScreenRectangle{
			left:   sb.positionThumb.left,
			top:    sb.positionUpperArrow.bottom,
			right:  sb.positionThumb.right,
			bottom: sb.positionThumb.top,
		}

		if pos.isInside(mouseX, mouseY) {
			result = UpperSpace
		}

		// now update it to be BottomSpace
		pos.top = sb.positionThumb.bottom
		pos.bottom = sb.positionBottomArrow.top
		if pos.isInside(mouseX, mouseY) {
			result = BottomSpace
		}
	}

	return result
}

func (sb *scrollbar) isMouseInside(px float64, py float64) bool {
	return sb.position.isInside(float32(px), float32(py))
}

func (sb *scrollbar) mouseButtonCallback(g *GUI, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey, mouseX float64, mouseY float64) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			switch sb.mouseHitTest(mouseX, mouseY) {
			case UpperArrow:
				sb.upperArrowIsDown = true
				g.terminal.ScreenScrollUp(1)

			case UpperSpace:
				g.terminal.ScrollPageUp()

			case Thumb:
				sb.thumbIsDragging = true
				sb.startedDraggingAtPosition = sb.scrollPosition
				sb.startedDraggingAtThumbTop = sb.positionThumb.top
				sb.offsetInThumbY = float32(mouseY) - sb.position.top - sb.positionThumb.top
				sb.scrollPositionDelta = 0

			case BottomSpace:
				g.terminal.ScrollPageDown()

			case BottomArrow:
				sb.bottomArrowIsDown = true
				g.terminal.ScreenScrollDown(1)
			}
		} else if action == glfw.Release {
			if sb.thumbIsDragging {
				sb.thumbIsDragging = false
			}

			if sb.upperArrowIsDown {
				sb.upperArrowIsDown = false
			}

			if sb.bottomArrowIsDown {
				sb.bottomArrowIsDown = false
			}
		}

		sb.isDirty = true
	}

	sb.resetElementColors(mouseX, mouseY)
}

func (sb *scrollbar) mouseMoveCallback(g *GUI, px float64, py float64) {
	sb.resetElementColors(px, py)

	if sb.thumbIsDragging {
		py -= float64(sb.position.top)

		minThumbTop := sb.positionUpperArrow.bottom
		maxThumbTop := sb.positionBottomArrow.top - sb.positionThumb.height()

		newThumbTop := float32(py) - sb.offsetInThumbY

		newPositionDelta := int((float32(sb.maxScrollPosition) * (newThumbTop - minThumbTop - sb.startedDraggingAtThumbTop)) / (maxThumbTop - minThumbTop))

		if newPositionDelta > sb.scrollPositionDelta {
			scrollLines := newPositionDelta - sb.scrollPositionDelta
			g.logger.Debugf("old position: %d, new position delta: %d, scroll down %d lines", sb.scrollPosition, newPositionDelta, scrollLines)
			g.terminal.ScreenScrollDown(uint16(scrollLines))
			sb.scrollPositionDelta = newPositionDelta
		} else if newPositionDelta < sb.scrollPositionDelta {
			scrollLines := sb.scrollPositionDelta - newPositionDelta
			g.logger.Debugf("old position: %d, new position delta: %d, scroll up %d lines", sb.scrollPosition, newPositionDelta, scrollLines)
			g.terminal.ScreenScrollUp(uint16(scrollLines))
			sb.scrollPositionDelta = newPositionDelta
		}

		sb.recalcElementPositions()
		g.logger.Debugf("new thumbTop: %f, fact thumbTop: %f, position: %d", newThumbTop, sb.positionThumb.top, sb.scrollPosition)
	}

	sb.isDirty = true
}

func (sb *scrollbar) resetElementColors(mouseX float64, mouseY float64) {
	part := sb.mouseHitTest(mouseX, mouseY)

	if sb.scrollPosition == 0 {
		sb.upperArrowBg = scrollbarColor_ButtonDisabledBg[:]
		sb.upperArrowFg = scrollbarColor_ButtonDisabledFg[:]
	} else if sb.upperArrowIsDown {
		sb.upperArrowFg = scrollbarColor_ButtonClickedFg[:]
		sb.upperArrowBg = scrollbarColor_ButtonClickedBg[:]
	} else if part == UpperArrow {
		sb.upperArrowFg = scrollbarColor_ButtonHoverFg[:]
		sb.upperArrowBg = scrollbarColor_ButtonHoverBg[:]
	} else {
		sb.upperArrowFg = scrollbarColor_ButtonNormalFg[:]
		sb.upperArrowBg = scrollbarColor_ButtonNormalBg[:]
	}

	if sb.scrollPosition == sb.maxScrollPosition {
		sb.bottomArrowBg = scrollbarColor_ButtonDisabledBg[:]
		sb.bottomArrowFg = scrollbarColor_ButtonDisabledFg[:]
	} else if sb.bottomArrowIsDown {
		sb.bottomArrowFg = scrollbarColor_ButtonClickedFg[:]
		sb.bottomArrowBg = scrollbarColor_ButtonClickedBg[:]
	} else if part == BottomArrow {
		sb.bottomArrowFg = scrollbarColor_ButtonHoverFg[:]
		sb.bottomArrowBg = scrollbarColor_ButtonHoverBg[:]
	} else {
		sb.bottomArrowFg = scrollbarColor_ButtonNormalFg[:]
		sb.bottomArrowBg = scrollbarColor_ButtonNormalBg[:]
	}

	if sb.thumbIsDragging {
		sb.thumbColor = scrollbarColor_ThumbClicked[:]
	} else if part == Thumb {
		sb.thumbColor = scrollbarColor_ThumbHover[:]
	} else {
		sb.thumbColor = scrollbarColor_ThumbNormal[:]
	}
}

func (sb *scrollbar) cursorEnterCallback(g *GUI, entered bool) {
	if !entered {
		sb.resetElementColors(-1, -1) // (-1, -1) ensures that no part is hovered by the mouse
		sb.isDirty = true
	}
}
