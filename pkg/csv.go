package pkg

import (
	"driver/internal"
	"driver/tools/zap"
	"encoding/csv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func GetLocationsFromCSVWithGivenRange(start, stop int) ([]internal.Model, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "./..")

	file, err := os.Open(filepath.Join(basepath, "./docs/Coordinates.csv"))
	if err != nil {
		zap.Logger.Error("Error while opening the file" + err.Error())
		return nil, err
	}
	defer file.Close()

	// create csv reader
	reader := csv.NewReader(file)

	// skip lines until we reach the start
	for i := 0; i < start; i++ {
		_, err := reader.Read()
		if err != nil {
			zap.Logger.Error("Error while reading the file" + err.Error())
			return nil, err
		}
	}

	// read lines until we reach the stop
	var locations []internal.Model
	for i := start; i < stop; i++ {
		record, err := reader.Read()
		if err != nil {
			// if we reach the end of the file, return the locations
			if err == io.EOF {
				break
			}
			zap.Logger.Error("Error while reading the file" + err.Error())
			return nil, err
		}

		// convert string to float
		lat, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			zap.Logger.Error("Error while converting string to float" + err.Error())
			return nil, err
		}

		// convert string to float
		lng, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			zap.Logger.Error("Error while converting string to float" + err.Error())
			return nil, err
		}

		locations = append(locations, internal.Model{
			ID: primitive.NewObjectID(),
			Location: internal.Location{
				Type: "Point",
				Coordinates: []float64{
					lng,
					lat,
				},
			},
		})
	}

	return locations, nil
}
