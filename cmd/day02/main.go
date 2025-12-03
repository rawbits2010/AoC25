package main

import (
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	invalidIdsPart1, err := findInvalidId(lines[0], true)
	if err != nil {
		log.Fatal(err)
	}

	sumInvalidIdsPart1, err := sumInvalidIds(invalidIdsPart1)
	if err != nil {
		log.Fatalf("error summing invalid Ids: %s", err)
	}

	invalidIdsPart2, err := findInvalidId(lines[0], false)
	if err != nil {
		log.Fatal(err)
	}

	sumInvalidIdsPart2, err := sumInvalidIds(invalidIdsPart2)
	if err != nil {
		log.Fatalf("error summing invalid Ids: %s", err)
	}

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", sumInvalidIdsPart1, sumInvalidIdsPart2)
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

func findInvalidId(line string, part1 bool) ([]string, error) {

	resultInvalidIds := make([]string, 0, 100)

	idRanges := strings.Split(line, ",")
	if len(idRanges) <= 0 {
		return nil, fmt.Errorf("no id ranges found")
	}

	for _, idRange := range idRanges {

		invalidIds := make([]string, 0, 100)

		limits := strings.Split(idRange, "-")
		if len(limits) < 2 {
			return nil, fmt.Errorf("invalid range found (%s)", idRange)
		}
		rangeStartStr := limits[0]
		rangeEndStr := limits[1]

		digitCount := len(rangeEndStr)
		sectionCounts := make([]uint, 0, digitCount)

		if part1 {
			sectionCounts = append(sectionCounts, 2)
		} else {
			for i := 2; i < digitCount; i++ {
				sectionCounts = append(sectionCounts, uint(i))
			}
			sectionCounts = append(sectionCounts, uint(digitCount))
		}

		for _, sectionCount := range sectionCounts {
			tmpInvalidIds, err := findRepeatingDigits(rangeStartStr, rangeEndStr, sectionCount)
			if err != nil {
				return nil, err
			}

			// filter duplicates
			for _, id := range tmpInvalidIds {
				if !slices.Contains(invalidIds, id) {
					invalidIds = append(invalidIds, id)
				}
			}
		}

		resultInvalidIds = append(resultInvalidIds, invalidIds...)
	}

	return resultInvalidIds, nil
}

func findRepeatingDigits(rangeStartStr, rangeEndStr string, sectionCount uint) ([]string, error) {

	numToTestStr := rangeStartStr

	// increase digit number to be divisable by section count
	for {
		if len(numToTestStr)%int(sectionCount) == 0 {
			break
		}
		numToTest := int(math.Pow10(len(numToTestStr)))
		numToTestStr = strconv.Itoa(numToTest)
	}

	// early bail
	if checkIfSmaller(rangeEndStr, numToTestStr) {
		return []string{}, nil
	}

	digitCount := len(numToTestStr) / int(sectionCount)
	firstPartStr := numToTestStr[:digitCount]
	firstPart, err := strconv.Atoi(firstPartStr)
	if err != nil {
		return nil, fmt.Errorf("weird number (%s): %w", firstPartStr, err)
	}

	// one-time range start check by comparing the 1st part to rest
	for i := 1; i < int(sectionCount); i++ {
		nextPartStr := numToTestStr[digitCount*i : (digitCount * (i + 1))]
		if nextPartStr == firstPartStr {
			continue
		}
		if checkIfSmaller(nextPartStr, firstPartStr) {
			break
		}
		firstPart++
		firstPartStr = strconv.Itoa(firstPart)
		break
	}

	invalidIds := make([]string, 0, 100)
	for {

		// test for validity
		newNumToTestStr := firstPartStr
		for i := 1; i < int(sectionCount); i++ {
			newNumToTestStr += firstPartStr
		}
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
