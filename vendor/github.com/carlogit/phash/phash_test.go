package phash

import (
	"os"
	"testing"
)

func TestGetHash(t *testing.T) {
	file1 := openFile("testdata/soccerball.jpg")
	defer file1.Close()

	file2 := openFile("testdata/soccerball (copy).jpg")
	defer file2.Close()

	phash1, _ := GetHash(file1)
	phash2, _ := GetHash(file2)

	expectedHash := "1101100110001110001101010011100000011010000000010111110111101000"
	if phash1 != expectedHash {
		t.Errorf("phash => %s, want %s", phash1, expectedHash)
	}

	if phash1 != phash2 {
		t.Errorf("phashes for files %s and %s must be the same, but they are different", file1.Name(), file2.Name())
	}
}

func TestSimilarImages(t *testing.T) {
	file1 := openFile("testdata/soccerball.jpg")
	defer file1.Close()

	file2 := openFile("testdata/soccerball (scaled down).jpg")
	defer file2.Close()

	file3 := openFile("testdata/soccerball (cropped).jpg")
	defer file3.Close()

	file4 := openFile("testdata/soccerball (modifications).jpg")
	defer file4.Close()

	file5 := openFile("testdata/soccerball (perspective).jpg")
	defer file5.Close()

	phash1, _ := GetHash(file1)
	phash2, _ := GetHash(file2)
	phash3, _ := GetHash(file3)
	phash4, _ := GetHash(file4)
	phash5, _ := GetHash(file5)

	distance := GetDistance(phash1, phash2)
	verifyDistanceInRange(t, file2.Name(), distance, 0, 1)

	distance = GetDistance(phash1, phash3)
	verifyDistanceInRange(t, file3.Name(), distance, 0, 1)

	distance = GetDistance(phash1, phash4)
	verifyDistanceInRange(t, file4.Name(), distance, 1, 3)

	distance = GetDistance(phash1, phash5)
	verifyDistanceInRange(t, file5.Name(), distance, 2, 5)
}

func TestGetDistance(t *testing.T) {
	var distancetests = []struct {
	hash1    string
	hash2    string
	distance int
	}{
		{"0010011100100010000001000101001000101110100101110", "0010011100100010000001000101001000101110100101110", 0},
		{"0010011100100000000001000101001000101110100101110", "0010011100100010000001000101001000101110100101111", 2},
		{"1111111111111111111111111111111111111111111111111", "0000000000000000000000000000000000000000000000000", 49},
	}
	
	for _, distancetest := range distancetests {
		distance := GetDistance(distancetest.hash1, distancetest.hash2)
		if distance != distancetest.distance {
			t.Errorf("distance between %s and %s => %d, want %d", distancetest.hash1, distancetest.hash2, distance, distancetest.distance)
		}
	}
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	return file
}

func verifyDistanceInRange(t *testing.T, comparedImageName string, distance, minDistance, maxDistance int) {
	if distance < minDistance || distance > maxDistance {
		t.Errorf("distance with %s => %d, want value between %d and %d", comparedImageName, distance, minDistance, maxDistance)
	}
}

