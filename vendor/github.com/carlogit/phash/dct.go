package phash

import (
	"math"
)

type dctPoint struct {
	xMax, yMax int
	xScales, yScales [2]float64
}

func (point *dctPoint) initializeScaleFactors() {
	point.xScales = [2]float64{ 1.0 / math.Sqrt(float64(point.xMax)), math.Sqrt(2.0 / float64(point.xMax))}
	point.yScales = [2]float64 { 1.0 / math.Sqrt(float64(point.yMax)), math.Sqrt(2.0 / float64(point.yMax))}
}

func (point *dctPoint) calculateValue(imageData [][]float64, x, y int) float64 {
	sum := float64(0.0)
	for i := 0; i < point.xMax; i++ {
		for j := 0; j < point.yMax; j++ {
			imageValue := float64(imageData[i][j])
			firstCosine := math.Cos(float64((1 + (2 * i)) * x) * math.Pi / float64(2 * point.xMax))
			secondCosine := math.Cos(float64((1 + (2 * j)) * y) * math.Pi / float64(2 * point.yMax))
			sum += (imageValue * firstCosine * secondCosine)
		}
	}
	return sum * point.getScaleFactor(x, y)
}

func (point *dctPoint) getScaleFactor(x, y int) float64 {
	xScaleFactor := point.xScales[1]
	if x == 0 {
		xScaleFactor = point.xScales[0]
	}
	yScaleFactor := point.yScales[1]
	if y == 0 {
		yScaleFactor = point.yScales[0]
	}
	
	return xScaleFactor * yScaleFactor
}
	
// getDCTMatrix Generates a DCT matrix from a given matrix.
// This is done using the Discrete Cosine Transformation (DCT) type-II algorithm.
func getDCTMatrix(matrix [][]float64) [][]float64 {
	xMax := len(matrix)
	yMax := len(matrix[0])

	dctPoint := &dctPoint{xMax:xMax, yMax:yMax}
	dctPoint.initializeScaleFactors()
	dctMatrix := make([][]float64, xMax)
	for x := 0; x < xMax; x++ {
		dctMatrix[x] = make([]float64, yMax)
		for y := 0; y < yMax; y++ {
			dctMatrix[x][y] = dctPoint.calculateValue(matrix, x, y)
		}
	}

	return dctMatrix
}
