package utils

import (
	"errors"
)

//HardcodeSlice Slice for filtering
var HardcodeSlice = []int{1, 3, 5}

func SplitSlice(sourceSlice []int, chunkLen int) ([][]int, error) {

	if len(sourceSlice) < 1 {
		return nil, errors.New("incorrect source slice")
	}

	if chunkLen <= 0 {
		return nil, errors.New("incorrect chunk length")
	}

	var chunkCount = (len(sourceSlice) / chunkLen) + 1
	var resultSlice = make([][]int, chunkCount)

	for i := 0; i < chunkCount; i++ {

		if chunkLen*(i+1) > len(sourceSlice) {
			resultSlice[i] = sourceSlice[(chunkLen * i):len(sourceSlice)]
			break
		}

		resultSlice[i] = sourceSlice[(chunkLen * i):(chunkLen * (i + 1))]
	}

	return resultSlice, nil
}

func SwapKeyValue(source map[int]int) map[int]int {

	var result = make(map[int]int, len(source))

	for key, value := range source {
		result[value] = key
	}

	return result
}

func FilterIntSlice(sourceSlice []int) []int {

	var resultSlice = make([]int, 0)

	for _, value := range sourceSlice {
		if !contains(HardcodeSlice, value) {
			resultSlice = append(resultSlice, value)
		}
	}

	return resultSlice
}

func contains(source []int, element int) bool {
	for _, value := range source {
		if value == element {
			return true
		}
	}

	return false
}
