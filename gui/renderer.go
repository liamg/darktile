package gui

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/liamg/aminal/buffer"
	"github.com/liamg/aminal/config"
	"github.com/liamg/aminal/glfont"
	"image"
	"math"
)

type OpenGLRenderer struct {
	areaWidth        int
	areaHeight       int
	areaX            int
	areaY            int
	cellWidth        float32
	cellHeight       float32
	termCols         uint
	termRows         uint
	cellPositions    map[[2]uint][2]float32
	config           *config.Config
	colourAttr       uint32
	program          uint32
	textureMap       map[*image.RGBA]uint32
	fontMap          *FontMap
	backgroundColour [3]float32

	rectRenderer *rectangleRenderer
}

type line struct {
	vao        uint32
	vbo        uint32
	cv         uint32
	colourAttr uint32
	colour     [3]float32
	points     []float32
	prog       uint32
}

func (r *OpenGLRenderer) CellWidth() float32 {
	return r.cellWidth
}

func (r *OpenGLRenderer) CellHeight() float32 {
	return r.cellHeight
}

func (r *OpenGLRenderer) newLine(x1 float32, y1 float32, x2 float32, y2 float32, dash float32, colourAttr uint32) *line {

	l := &line{}

	halfAreaWidth := float32(r.areaWidth / 2)
	halfAreaHeight := float32(r.areaHeight / 2)

	x1 = (x1 - halfAreaWidth) / halfAreaWidth
	y1 = -(y1 - (halfAreaHeight)) / halfAreaHeight
	x2 = (x2 - halfAreaWidth) / halfAreaWidth
	y2 = -(y2 - (halfAreaHeight)) / halfAreaHeight

	var xgap float32
	var tan float32
	if x2-x1 != 0 {
		tan = (y2 - y1) / (x2 - x1)
		xgap = dash / float32(math.Cos(math.Atan(float64(tan)))) / halfAreaWidth
	}

	l.points = []float32{
		x1, y1, 0,
	}

	if xgap == 0 {
		l.points = append(l.points, x2, y2, 0)
	} else {
		end := x1 + xgap
		for {
			var y float32
			if end >= x2 {
				end = x2
				y = y2
			} else {
				y = y1 + tan*(end-x1)
			}
			l.points = append(l.points, end, y, 0)
			start := end + xgap
			if start >= x2 {
				break
			}
			y = y1 + tan*(start-x1)
			l.points = append(l.points, start, y, 0)
			end = start + xgap
		}
	}

	l.colourAttr = colourAttr
	l.prog = r.program

	// SHAPE
	gl.GenBuffers(1, &l.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(l.points), gl.Ptr(&l.points[0]), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &l.vao)
	gl.BindVertexArray(l.vao)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, l.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	// colour
	gl.GenBuffers(1, &l.cv)

	return l
}

func (l *line) Draw() {
	gl.UseProgram(l.prog)
	gl.BindVertexArray(l.vao)
	gl.DrawArrays(gl.LINES, 0, int32(len(l.points)/3))
}

func (l *line) setColour(colour [3]float32) {
	if l.colour == colour {
		return
	}

	c := make([]float32, len(l.points))

	for i := 0; i < len(c); i += 3 {
		c[i] = colour[0]
		c[i+1] = colour[1]
		c[i+2] = colour[2]
	}

	gl.UseProgram(l.prog)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.cv)
	gl.BufferData(gl.ARRAY_BUFFER, len(c)*4, gl.Ptr(c), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(l.colourAttr)
	gl.VertexAttribPointer(l.colourAttr, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	l.colour = colour
}

func (l *line) Free() {
	gl.DeleteVertexArrays(1, &l.vao)
	gl.DeleteBuffers(1, &l.vbo)
	gl.DeleteBuffers(1, &l.cv)

	l.vao = 0
	l.vbo = 0
	l.cv = 0
}

func NewOpenGLRenderer(config *config.Config, fontMap *FontMap, areaX int, areaY int, areaWidth int, areaHeight int, colourAttr uint32, program uint32) (*OpenGLRenderer, error) {
	rectRenderer, err := newRectangleRenderer()
	if err != nil {
		return nil, err
	}
	r := &OpenGLRenderer{
		areaWidth:     areaWidth,
		areaHeight:    areaHeight,
		areaX:         areaX,
		areaY:         areaY,
		cellPositions: map[[2]uint][2]float32{},
		config:        config,
		colourAttr:    colourAttr,
		program:       program,
		textureMap:    map[*image.RGBA]uint32{},
		fontMap:       fontMap,
		rectRenderer:  rectRenderer,
	}
	r.SetArea(areaX, areaY, areaWidth, areaHeight)
	return r, nil
}

// This method ensures that all OpenGL resources are deleted correctly
func (r *OpenGLRenderer) Free() {
	for _, tex := range r.textureMap {
		gl.DeleteTextures(1, &tex)
	}
	r.textureMap = map[*image.RGBA]uint32{}

	r.fontMap.Free()

	gl.DeleteProgram(r.program)
	r.program = 0

	if r.rectRenderer != nil {
		r.rectRenderer.Free()
		r.rectRenderer = nil
	}
}

func (r *OpenGLRenderer) GetTermSize() (uint, uint) {
	return r.termCols, r.termRows
}

func (r *OpenGLRenderer) SetArea(areaX int, areaY int, areaWidth int, areaHeight int) {
	r.areaWidth = areaWidth
	r.areaHeight = areaHeight
	r.areaX = areaX
	r.areaY = areaY
	f := r.fontMap.DefaultFont()
	_, r.cellHeight = f.MaxSize()
	r.cellWidth, _ = f.Size("X")
	//= f.LineHeight()   // includes vertical padding
	r.termCols = uint(math.Floor(float64(float32(r.areaWidth) / r.cellWidth)))
	r.termRows = uint(math.Floor(float64(float32(r.areaHeight) / r.cellHeight)))
}

func (r *OpenGLRenderer) ConvertCoordinates(col uint, row uint) (float32, float32) {
	left := float32(float32(col) * r.cellWidth)
	top := float32(float32(row) * r.cellHeight)

	return left, top
}

func (r *OpenGLRenderer) DrawCursor(col uint, row uint, colour config.Colour) {
	left, top := r.ConvertCoordinates(col, row)
	r.rectRenderer.render(left, top, r.cellWidth, r.cellHeight, colour)
}

func (r *OpenGLRenderer) DrawCellBg(cell buffer.Cell, col uint, row uint, colour *config.Colour, force bool) {

	var bg [3]float32

	if colour != nil {
		bg = *colour
	} else {
		bg = cell.Bg()
	}

	if bg != r.backgroundColour || force {
		left, top := r.ConvertCoordinates(col, row)
		r.rectRenderer.render(left, top, r.cellWidth, r.cellHeight, bg)
	}
}

func (r *OpenGLRenderer) getUndelineThickness() float32 {
	thickness := r.cellHeight / 16
	if thickness < 1 {
		thickness = 1
	}
	return thickness
}

// DrawUnderline draws a line under 'span' characters starting at (col, row)
func (r *OpenGLRenderer) DrawUnderline(span int, col uint, row uint, colour [3]float32) {
	//calculate coordinates
	x := float32(float32(col) * r.cellWidth)
	y := (float32(row+1))*r.cellHeight + r.fontMap.DefaultFont().MinY()*0.25

	thickness := r.getUndelineThickness()
	r.rectRenderer.render(x, y, r.cellWidth*float32(span), thickness, colour)
}

func (r *OpenGLRenderer) DrawLinkLine(span int, col uint, row uint, colour [3]float32) {
	//calculate coordinates
	x := float32(float32(col) * r.cellWidth)
	y := (float32(row+1))*r.cellHeight + r.fontMap.DefaultFont().MinY()*0.5
	line := r.newLine(x, y, x+r.cellWidth*float32(span), y, r.cellWidth/4, r.colourAttr)

	line.setColour(colour)
	line.Draw()

	line.Free()
}

func (r *OpenGLRenderer) DrawCellText(text string, col uint, row uint, alpha float32, colour [3]float32, bold bool) {

	var f *glfont.Font
	if bold {
		f = r.fontMap.BoldFont()
	} else {
		f = r.fontMap.DefaultFont()
	}

	f.SetColor(colour[0], colour[1], colour[2], alpha)

	x := float32(r.areaX) + float32(col)*r.cellWidth
	y := float32(r.areaY) + (float32(row+1) * r.cellHeight) + f.MinY()

	f.Print(x, y, text)
}

func (r *OpenGLRenderer) DrawCellImage(cell buffer.Cell, col uint, row uint) {

	img := cell.Image()

	if img == nil {
		return
	}

	ix := float32(col) * r.cellWidth
	iy := float32(r.areaHeight) - (float32(row+1) * r.cellHeight)
	iy -= float32(cell.Image().Bounds().Size().Y)
	gl.UseProgram(r.program)

	var tex uint32

	tex, ok := r.textureMap[img]
	if !ok {
		gl.Enable(gl.TEXTURE_2D)
		gl.GenTextures(1, &tex)
		gl.BindTexture(gl.TEXTURE_2D, tex)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(img.Bounds().Size().X),
			int32(img.Bounds().Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(img.Pix),
		)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.Disable(gl.TEXTURE_2D)

		gl.Disable(gl.BLEND)

		r.textureMap[img] = tex
	}

	var w = float32(img.Bounds().Size().X)
	var h = float32(img.Bounds().Size().Y)

	var readFboId uint32
	gl.GenFramebuffers(1, &readFboId)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, readFboId)

	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D, tex, 0)
	gl.BlitFramebuffer(0, 0, int32(w), int32(h),
		int32(ix), int32(iy), int32(ix+w), int32(iy+h),
		gl.COLOR_BUFFER_BIT, gl.LINEAR)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.DeleteFramebuffers(1, &readFboId)
}
