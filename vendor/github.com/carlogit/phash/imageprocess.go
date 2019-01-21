package phash

import (
	"image"
)

func getImageMatrix(img image.Image) [][]float64 {
	xSize := img.Bounds().Max.X
	ySize := img.Bounds().Max.Y

	vals := make([][]float64, xSize)

	for x := 0; x < xSize; x++ {
		vals[x] = make([]float64, ySize)
		for y := 0; y < ySize; y++ {
			vals[x][y] = getXYValue(img, x, y)
		}
	}

	return vals
}

func getXYValue(img image.Image, x int, y int) float64 {
	_, _, b, _ := img.At(x, y).RGBA()
	return float64(b)
}
