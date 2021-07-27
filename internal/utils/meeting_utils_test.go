package utils

import (
	"errors"
	model "github.com/ozoncp/ocp-meeting-api/internal/models"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSplitToBulks(t *testing.T) {

	var now = time.Now()

	sourceSlice := []model.Meeting{
		{1, 1, "", now, now},
		{2, 2, "", now, now},
		{3, 3, "", now, now},
		{4, 4, "", now, now},
		{5, 5, "", now, now},
	}

	expectedSlice := [][]model.Meeting{
		{
			{1, 1, "", now, now},
			{2, 2, "", now, now},
			{3, 3, "", now, now},
		},
		{
			{4, 4, "", now, now},
			{5, 5, "", now, now},
		},
	}

	result := SplitToBulks(sourceSlice, 3)

	if !reflect.DeepEqual(expectedSlice, result) {
		t.Errorf("Incorerect result: %v. Expected: %v", result, expectedSlice)
	}
}

func TestSliceToMap(t *testing.T) {

	var now = time.Now()

	sourceSlice := []model.Meeting{
		{1, 1, "", now, now},
		{2, 2, "", now, now},
		{3, 3, "", now, now},
	}

	expectedMap := map[uint64]model.Meeting{
		1: {1, 1, "", now, now},
		2: {2, 2, "", now, now},
		3: {3, 3, "", now, now},
	}

	result, _ := SliceToMap(sourceSlice)

	if !reflect.DeepEqual(expectedMap, result) {
		t.Errorf("Incorerect result: %v. Expected: %v", result, expectedMap)
	}
}

func TestSliceToMapError(t *testing.T) {

	var sourceSlice []model.Meeting
	var expectedErr = errors.New("incorrect slice of meetings")
	_, err := SliceToMap(sourceSlice)

	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("Expected error: %v", expectedErr.Error())
	}
}

func TestCheckFile(t *testing.T) {

	var fileName = "test"

	if err := CheckFile(fileName, 3); err != nil {
		t.Error(err.Error())
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0755)

	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	defer os.Remove(fileName)
}
