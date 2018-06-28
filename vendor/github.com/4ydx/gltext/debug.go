package gltext

import (
	"fmt"
	"runtime"
)

var IsDebug = false

func DebugPrefix() string {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("DB: [%s:%d]", fn, line)
}

func TextDebug(message string) {
	if IsDebug {
		pc, fn, line, _ := runtime.Caller(1)
		fmt.Printf("[error] in %s[%s:%d] %s", runtime.FuncForPC(pc).Name(), fn, line, message)
	}
}

// PrintVBO prints the individual index locations as well as the texture locations
//
// (0,0) (x1,y1): This shows the layout of the runes.  There relative locations to one another can be seen here.
// - If called just after makeBufferData, the left-most x value will start at 0.
// - If called after centerTheData, all indices will have been shifted so that the entire text value is
//   centered around the screen's origin of (0,0).
//
// (U,V) (u1,v1) -> (x,y): The (x,y) values refer to pixel locations within the texture
// - Open the texture in an image editor and, using the upper left hand corner as (0,0)
//   move to the location (x,y).  This is where opengl will pinpoint your rune within the image.
func PrintVBO(vbo []float32, w, h float32) {
	if len(vbo)%4 != 0 {
		fmt.Println("VBO appears to have an incorrect size.  Should be a multiple of 4.")
	}
	// drawing a quad take 16 floats: (2 x,y + 2 u,v) * 4 indices
	for i := 0; i < len(vbo); i += 16 {
		fmt.Println("Quad")
		at := i
		fmt.Printf(
			"(0,0) (%.2f,%.2f); (U,V) (%f,%f) -> (%f,%f)\n",
			vbo[at], vbo[at+1], vbo[at+2], vbo[at+3], vbo[at+2]*w, vbo[at+3]*h,
		)
		at += 4
		fmt.Printf(
			"(1,0) (%.2f,%.2f); (U,V) (%f,%f) -> (%f,%f)\n",
			vbo[at], vbo[at+1], vbo[at+2], vbo[at+3], vbo[at+2]*w, vbo[at+3]*h,
		)
		at += 4
		fmt.Printf(
			"(1,1) (%.2f,%.2f); (U,V) (%f,%f) -> (%f,%f)\n",
			vbo[at], vbo[at+1], vbo[at+2], vbo[at+3], vbo[at+2]*w, vbo[at+3]*h,
		)
		at += 4
		fmt.Printf(
			"(0,1) (%.2f,%.2f); (U,V) (%f,%f) -> (%f,%f)\n",
			vbo[at], vbo[at+1], vbo[at+2], vbo[at+3], vbo[at+2]*w, vbo[at+3]*h,
		)
	}
}
