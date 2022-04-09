package voxblox

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// readPointsFromCSV reads XYZ points from a CSV file.
func readPointsFromCSV(fileName string) []Point {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	var points []Point
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// Put the point in the list.
		x, _ := strconv.ParseFloat(rec[0], 64)
		y, _ := strconv.ParseFloat(rec[1], 64)
		z, _ := strconv.ParseFloat(rec[2], 64)
		points = append(points, Point{x, y, z})
	}
	return points
}

func TestStepICP(t *testing.T) {
	sourcePoints := readPointsFromCSV("../testdata/source.csv")
	targetPoints := readPointsFromCSV("../testdata/target.csv")
	assert.Equal(t, len(sourcePoints), len(targetPoints))

	transform := [6]float64{0.01, 0.05, 0.01, 0.001, 0.001, 0.001}
	err := stepICP(sourcePoints, targetPoints, &transform)
	assert.NoError(t, err)
	assert.InDelta(t, -0.03887, transform[0], 0.0001)
	assert.InDelta(t, -0.03588, transform[1], 0.0001)
	assert.InDelta(t, 0.00710, transform[2], 0.0001)
	assert.InDelta(t, 0.00047, transform[3], 0.0001)
	assert.InDelta(t, 0.00052, transform[4], 0.0001)
	assert.InDelta(t, 0.00050, transform[5], 0.0001)

	err = stepICP(sourcePoints, targetPoints, &transform)
	assert.NoError(t, err)
	assert.InDelta(t, -0.03887, transform[0], 0.0001)
	assert.InDelta(t, -0.07911, transform[1], 0.0001)
	assert.InDelta(t, 0.00565, transform[2], 0.0001)
	assert.InDelta(t, 0.00021, transform[3], 0.0001)
	assert.InDelta(t, 0.00028, transform[4], 0.0001)
	assert.InDelta(t, 0.00025, transform[5], 0.0001)
}

func TestGetGradient(t *testing.T) {
	tsdfLayer := NewTsdfLayer(0.05, 16)
	_ = tsdfLayer
}

func TestAddNormalizedPointInfo(t *testing.T) {
	point := Point{2.66496301, 1.69732249, 1.48568201}
	pointNormal := Point{-0.340126932, -0.911325812, -0.231946185}
	infoVector := [6]float64{kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon}
	addNormalizedPointInfo(point, pointNormal, &infoVector)
	assert.InEpsilon(t, 0.231373653, infoVector[0], kEpsilon)
	assert.InEpsilon(t, 1.66103041, infoVector[1], kEpsilon)
	assert.InEpsilon(t, 0.107599065, infoVector[2], kEpsilon)
	assert.InEpsilon(t, 3.97628975, infoVector[3], kEpsilon)
	assert.InEpsilon(t, 1.274863, infoVector[4], kEpsilon)
	assert.InEpsilon(t, 12.4632406, infoVector[5], kEpsilon)
}

func TestMatchPoints(t *testing.T) {
	// TODO
}
