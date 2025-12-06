package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {
	lines := inputhandler.ReadInput()

	idRanges, lastRangeIdx, err := readRanges(lines)

	if err != nil {
		log.Fatalf("error processing ranges: %s", err)
	}

	freshCount := countFreshIngredients(idRanges, lines[lastRangeIdx:])

	idCount := sumFreshIds(idRanges)

	fmt.Printf("Result - Part 1: %d, Part 2: %s\n", freshCount, idCount)
}

type IDRange struct {
	start       string
	end         string
	isRedundant bool
}

func readRanges(lines []string) ([]IDRange, int, error) {

	if len(lines) == 0 {
		return nil, -1, fmt.Errorf("empty database provided")
	}

	idRanges := make([]IDRange, 0, 169)

	for currIdx, line := range lines {

		if len(line) == 0 {
			return idRanges, currIdx - 1, nil
		}

		limits := strings.Split(line, "-")
		if len(limits) == 1 {
			if len(idRanges) != 0 {
				return nil, 0, fmt.Errorf("malformed data found at %d (%s)", currIdx, line)
			} else {
				return []IDRange{}, 0, nil
			}
		} else if len(limits) > 2 {
			return nil, 0, fmt.Errorf("malformed data found at %d (%s)", currIdx, line)
		}

		idRanges = append(idRanges,
			IDRange{
				start: limits[0],
				end:   limits[1],
			})
	}

	return idRanges, len(lines) - 1, nil
}

func countFreshIngredients(idRanges []IDRange, ingredientIds []string) int {

	freshCount := 0
	for _, id := range ingredientIds {

		if len(id) == 0 {
			continue
		}

		for _, idRange := range idRanges {

			if checkIfSmaller(id, idRange.start) {
				continue
			}

			if id == idRange.end || checkIfSmaller(id, idRange.end) {
				freshCount++
				break
			}
		}
	}

	return freshCount
}

func sumFreshIds(idRanges []IDRange) string {

	idRanges = filterOverlaps(idRanges)

	idCount := ""
	for _, idRange := range idRanges {
		if idRange.isRedundant {
			continue
		}

		diff := substract(idRange.start, idRange.end)
		count := increment(diff)

		idCount = add(count, idCount)

		fmt.Printf("rEnd: %s, rStart: %s, diff: %s, count: %s, idCount: %s\n", idRange.end, idRange.start, diff, count, idCount)
	}

	return idCount
}

func filterOverlaps(idRanges []IDRange) []IDRange {

	for i := 0; i < len(idRanges); i++ {

		start := idRanges[i].start
		end := idRanges[i].end
		isRedundant := false
		for j := 0; j < len(idRanges); j++ {
			if j == i {
				continue
			}
			if idRanges[j].isRedundant {
				continue
			}

			startIsHigher := !checkIfSmaller(start, idRanges[j].start)
			endIsLower := end == idRanges[j].end || checkIfSmaller(end, idRanges[j].end)

			if startIsHigher && endIsLower {
				isRedundant = true
				break
			}

			if startIsHigher {
				if start == idRanges[j].end || checkIfSmaller(start, idRanges[j].end) {
					start = increment(idRanges[j].end)
				}
			}

			if endIsLower {
				if !checkIfSmaller(end, idRanges[j].start) {
					end = decrement(idRanges[j].start)
				}
			}
		}

		if isRedundant {
			idRanges[i].isRedundant = true
		} else {
			idRanges[i].start = start
			idRanges[i].end = end
		}
	}

	return idRanges
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

const zeroDigit = '0'
const nineDigit = '9'

// substract substracts 'this' from 'from' on the actual strings
// without converting to integers.
// NOTE: won't validate numbers and is not handle negative
// numbers or remainders.
func substract(this, from string) string {

	remDigits := make([]byte, len(from))

	digitCountDiff := len(from) - len(this)
	overflow := 0
	var i int
	for i = len(from) - 1; i >= 0; i-- {

		var diff int
		if i >= digitCountDiff {
			diff = int(from[i]) - (int(this[i-digitCountDiff]) + overflow)
		} else {
			if overflow != 0 {
				diff = int(from[i]) - (zeroDigit + overflow)
			} else {
				diff = int(from[i]) - zeroDigit
			}
		}

		if diff < 0 {
			diff = 10 + diff
			overflow = 1
		} else {
			overflow = 0
		}

		numByte := diff + zeroDigit
		remDigits[i] = byte(numByte)
	}

	numStr := string(remDigits[i+1:])
	numStr = strings.TrimLeft(numStr, "0")
	if len(numStr) == 0 {
		return "0"
	}

	return numStr
}

// add adds 'this' to 'to' on the actual strings without converting
// to integers.
// NOTE: won't validate numbers and is not handle negative
// numbers or sums.
func add(this, to string) string {

	var longerNum string
	var shorterNum string
	if len(this) >= len(to) {
		longerNum = this
		shorterNum = to
	} else {
		longerNum = to
		shorterNum = this
	}

	sumDigits := make([]byte, len(longerNum)+1)

	digitCountDiff := len(longerNum) - len(shorterNum)
	overflow := 0
	for i := len(longerNum) - 1; i >= 0; i-- {

		var sum int
		if i >= digitCountDiff {
			sum = int(longerNum[i]-zeroDigit) + int(shorterNum[i-digitCountDiff]-zeroDigit) + overflow
		} else {
			sum = int(longerNum[i]-zeroDigit) + overflow
		}

		if sum > 9 {
			sum -= 10
			overflow = 1
		} else {
			overflow = 0
		}

		sumDigits[i+1] = byte(sum + zeroDigit)
	}

	sumDigits[0] = byte(zeroDigit + overflow)

	numStr := string(sumDigits)
	numStr = strings.TrimLeft(numStr, "0")
	if len(numStr) == 0 {
		return "0"
	}

	return numStr
}

// increment increments 'this' by 1 on the actual string
// without converting to integers.
// NOTE: won't validate numbers and not works for negative numbers.
func increment(this string) string {

	incDigits := make([]byte, len(this)+1)

	overflow := 0
	for i := len(this) - 1; i >= 0; i-- {

		var num int
		if i == len(this)-1 {
			num = int(this[i]) + 1 + overflow
		} else {
			num = int(this[i]) + overflow
		}

		if num > nineDigit {
			num -= 10
			overflow = 1
		} else {
			overflow = 0
		}

		incDigits[i+1] = byte(num)
	}

	incDigits[0] = byte(zeroDigit + overflow)

	numStr := string(incDigits)
	numStr = strings.TrimLeft(numStr, "0")

	return numStr
}

// decrement decrements 'this' by 1 on the actual string
// without converting to integers.
// NOTE: won't validate numbers and not works for negative numbers.
func decrement(this string) string {

	decDigits := make([]byte, len(this))

	overflow := 0
	for i := len(this) - 1; i >= 0; i-- {

		var num int
		if i == len(this)-1 {
			num = int(this[i]) - (1 + overflow)
		} else {
			num = int(this[i]) - overflow
		}

		if num < zeroDigit {
			num += 10
			overflow = 1
		} else {
			overflow = 0
		}

		decDigits[i] = byte(num)
	}

	numStr := string(decDigits)
	if len(numStr) == 1 {
		return numStr
	}

	numStr = strings.TrimLeft(numStr, "0")
	return numStr
}
