// Package phash computes a phash string for a JPEG image and retrieves
// the hamming distance between two phash strings.
package phash

import (
	"io"

	"github.com/disintegration/imaging"
)

// GetHash returns a phash string for a JPEG image
func GetHash(reader io.Reader) (string, error) {
	image, err := imaging.Decode(reader)

	if err != nil {
		return "", err
	}
	
	image = imaging.Resize(image, 32, 32, imaging.Lanczos)
	image = imaging.Grayscale(image)
	
	imageMatrixData := getImageMatrix(image)	
	dctMatrix := getDCTMatrix(imageMatrixData)
	
	smallDctMatrix := reduceMatrix(dctMatrix, 8)
	dctMeanValue := calculateMeanValue(smallDctMatrix)
	return buildHash(smallDctMatrix, dctMeanValue), nil
}

// GetDistance returns the hamming distance between two phashes
func GetDistance(hash1, hash2 string) int {
	distance := 0
	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			distance++
		}
	}

	return distance
}

func buildHash(dctMatrix [][]float64, dctMeanValue float64) string {
	var hash string
	var xSize = len(dctMatrix)
	var ySize = len(dctMatrix[0])

	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {
			if dctMatrix[x][y] > dctMeanValue {
				hash += "1"
			} else {
				hash += "0"
			}
		}
	}

	return hash
}
