package phash

import (
	"testing"
)

func TestReduceMatrix(t *testing.T) {
	matrix := make([][]float64, 9)
	for i:=0; i < 9; i++ {
		matrix[i] = make([]float64, 9)
		for j:=0; j < 9; j++ {
			matrix[i][j] = float64(i * 10 + j)
		}
	}
	
	
	newMatrix := reduceMatrix(matrix, 2)

	if len(newMatrix) != 2 {
		t.Errorf("number of rows => %d, want %d", len(newMatrix), 2)
	}
	
	if len(newMatrix[0]) != 2 {
		t.Errorf("number of columns for row 1 => %d, want %d", len(newMatrix[0]), 2)
	}

	if len(newMatrix[1]) != 2 {
		t.Errorf("number of columns for row 2 => %d, want %d", len(newMatrix[1]), 2)
	}

	if newMatrix[0][0] != 0 {
		t.Errorf("value for [%d, %d] => %d, want %d", newMatrix[0][0], 0)
	}

	if newMatrix[0][1] != 1 {
		t.Errorf("value for [%d, %d] => %d, want %d", newMatrix[0][1], 1)
	}
	if newMatrix[1][0] != 10 {
		t.Errorf("value for [%d, %d] => %d, want %d", newMatrix[1][0], 10)
	}
	if newMatrix[1][1] != 11 {
		t.Errorf("value for [%d, %d] => %d, want %d", newMatrix[1][1], 11)
	}
}

func TestCalculateMeanValue(t *testing.T) {
	matrix := make([][]float64, 3)
	matrix[0] = make([]float64, 3)
	matrix[1] = make([]float64, 3)
	matrix[2] = make([]float64, 3)
	
	matrix[0][0] = 10
	matrix[0][1] = 20
	matrix[0][2] = 30
	
	matrix[1][0] = 5
	matrix[1][1] = 25
	matrix[1][2] = 45
	
	matrix[2][0] = 23
	matrix[2][1] = 34
	matrix[2][2] = 66
	
	meanValue := calculateMeanValue(matrix)
	
	if meanValue != 31 {
		t.Errorf("mean value for matrix => %d, want %d", meanValue, 31)
	}
}

