package utils

import (
	"errors"
	model "github.com/ozoncp/ocp-meeting-api/internal/models"
	"os"
)

func SplitToBulks(meetings []model.Meeting, butchSize uint) [][]model.Meeting {

	if len(meetings) == 0 || butchSize == 0 {
		return nil
	}

	var chunkCount = len(meetings) / int(butchSize)

	if (len(meetings) % int(butchSize)) > 0 {
		chunkCount++
	}

	var resultSlice = make([][]model.Meeting, chunkCount)

	for i := 0; i < chunkCount; i++ {

		if int(butchSize)*(i+1) > len(meetings) {
			resultSlice[i] = meetings[(int(butchSize) * i):len(meetings)]
			break
		}

		resultSlice[i] = meetings[(int(butchSize) * i):(int(butchSize) * (i + 1))]
	}

	return resultSlice
}

func SliceToMap(meetings []model.Meeting) (map[uint64]model.Meeting, error) {

	if len(meetings) == 0 {
		return nil, errors.New("incorrect slice of meetings")
	}

	result := make(map[uint64]model.Meeting, len(meetings))

	for _, value := range meetings {
		result[value.Id] = value
	}

	return result, nil
}

func CheckFile(path string, count uint) error {
	var err error
	var file *os.File

	check := func() {
		file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)

		if err != nil {
			return
		}

		defer func() {
			err = file.Close()
		}()
	}

	for ; count > 0; count-- {
		check()
	}

	return err
}
