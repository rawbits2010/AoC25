package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	invalidIds, err := findInvalidId(lines[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(invalidIds)

	sumInvalidIds, err := sumInvalidIds(invalidIds)
	if err != nil {
		log.Fatalf("error summing invalid Ids: %s", err)
	}

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", sumInvalidIds, 0)
}

func sumInvalidIds(invalidIds []string) (int, error) {
	sum := 0
	for _, id := range invalidIds {
		idNum, err := strconv.Atoi(id)
		if err != nil {
			return 0, fmt.Errorf("error converting invalid id (%s): %w", id, err)
		}
		sum += idNum
	}
	return sum, nil
}

func findInvalidId(line string) ([]string, error) {

	invalidIds := make([]string, 0, 100)

	idRanges := strings.Split(line, ",")
	if len(idRanges) <= 0 {
		return nil, fmt.Errorf("no id ranges found")
	}

	for _, idRange := range idRanges {

		limits := strings.Split(idRange, "-")
		if len(limits) < 2 {
			return nil, fmt.Errorf("invalid range found (%s)", idRange)
		}
		rangeStartStr := limits[0]
		rangeEndStr := limits[1]

		numToTestStr := rangeStartStr

		// next even digits
		if len(numToTestStr)%2 == 1 {
			numToTest := int(math.Pow10(len(rangeStartStr)))
			numToTestStr = strconv.Itoa(numToTest)
		}

		digitCount := len(numToTestStr)
		firstPartStr := numToTestStr[:digitCount/2]
		firstPart, err := strconv.Atoi(firstPartStr)
		if err != nil {
			return nil, fmt.Errorf("weird number (%s): %w", firstPartStr, err)
		}

		// the smallest second part we can have is >= then the second part
		// of the start of the range, so let' skip to that
		secondPartStr := numToTestStr[digitCount/2:]
		if firstPartStr != secondPartStr && checkIfSmaller(firstPartStr, secondPartStr) {
			firstPart++
			firstPartStr = strconv.Itoa(firstPart)
		}

		for {

			// test for validity
			newNumToTestStr := firstPartStr + firstPartStr
			if newNumToTestStr == rangeEndStr || checkIfSmaller(newNumToTestStr, rangeEndStr) {
				// within range
				invalidIds = append(invalidIds, newNumToTestStr)
			} else {
				// out of range
				break
			}

			// create next number
			firstPart++
			firstPartStr = strconv.Itoa(firstPart)

		}

	}

	return invalidIds, nil
}

func checkIfSmaller(toTest, limit string) bool {

	limitDigitCount := len(limit)
	toTestDigitCount := len(toTest)

	if limitDigitCount < toTestDigitCount {
		return false
	}
	if limitDigitCount > toTestDigitCount {
		return true
	}

	for i := 0; i < limitDigitCount; i++ {
		if toTest[i] > limit[i] {
			return false
		}
		if toTest[i] < limit[i] {
			return true
		}
	}

	return false
}
