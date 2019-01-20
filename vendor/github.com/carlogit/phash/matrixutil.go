package phash

func reduceMatrix(dctMatrix [][]float64, size int) [][]float64 {
	newMatrix := make([][]float64, size)
	for x := 0; x < size; x++ {
		newMatrix[x] = make([]float64, size)
		for y := 0; y < size; y++ {
			newMatrix[x][y] = dctMatrix[x][y]
		}
	}
	
	return newMatrix
}


func calculateMeanValue(dctMatrix [][]float64) float64 {
	var total float64
	var xSize = len(dctMatrix)
	var ySize = len(dctMatrix[0])

	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {
			total += dctMatrix[x][y]
		}
	}

	total -= dctMatrix[0][0]

	avg := total / float64((xSize * ySize) - 1)

	return avg
}
