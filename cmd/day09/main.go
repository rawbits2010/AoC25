package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	coords, err := parseCoords(lines)
	if err != nil {
		log.Fatal(err)
	}

	resultP1 := largestArea(coords)
	//resultP1 := part1BruteForce(coords)

	var resultP2 int

	fmt.Printf("Result - Part 1: %s, Part 2: %d\n", resultP1, resultP2)
}

func parseCoords(lines []string) ([]Coords, error) {

	if len(lines) == 0 {
		return nil, fmt.Errorf("no data provided")
	}

	coords := make([]Coords, len(lines))
	for i := 0; i < len(lines); i++ {

		x, y, err := getCoords(lines[i])
		if err != nil {
			return nil, fmt.Errorf("error parsing line at %d (%s) :%w", 0, lines[0], err)
		}

		coords[i] = Coords{x, y}
	}

	return coords, nil
}

func getCoords(line string) (string, string, error) {
	nums := strings.Split(line, ",")
	if len(nums) != 2 {
		return "", "", fmt.Errorf("invalid format")
	}
	return nums[0], nums[1], nil
}

type Coords struct {
	x, y string
}

func part1BruteForce(coords []Coords) string {
	areaMax := "0"
	lastI := 0
	lastJ := 0
	for i := 0; i < len(coords); i++ {
		for j := i; j < len(coords); j++ {
			first := coords[i]
			second := coords[j]

			var smallX, bigX string
			if checkIfSmaller(first.x, second.x) {
				smallX = first.x
				bigX = second.x
			} else {
				smallX = second.x
				bigX = first.x
			}

			var smallY, bigY string
			if checkIfSmaller(first.y, second.y) {
				smallY = first.y
				bigY = second.y
			} else {
				smallY = second.y
				bigY = first.y
			}

			area := calcArea(smallX, smallY, bigX, bigY)
			if checkIfSmaller(areaMax, area) {
				areaMax = area
				lastI = i
				lastJ = j
			}
		}
	}

	fmt.Println("Point 1 - X: ", coords[lastI].x, ", Y: ", coords[lastI].y)
	fmt.Println("Point 2 - X: ", coords[lastJ].x, ", Y: ", coords[lastJ].y)

	return areaMax
}

func largestArea(coords []Coords) string {

	xMinMax, yMinMax := findMinMax(coords)
	midCoords := calcMidCoord(xMinMax, yMinMax)
	sortedCoords := sortCoords(midCoords, coords)
	// NOTE: it is possible to not have points in a quadrant!
	// This method might fail on that case.

	tlSet := filterToCorner(sortedCoords.tl, Coords{xMinMax.min, yMinMax.min}, midCoords, filterTopLeft)
	brSet := filterToCorner(sortedCoords.br, Coords{xMinMax.max, yMinMax.max}, midCoords, filterBottomRight)
	trSet := filterToCorner(sortedCoords.tr, Coords{xMinMax.max, yMinMax.min}, midCoords, filterTopRight)
	blSet := filterToCorner(sortedCoords.bl, Coords{xMinMax.min, yMinMax.max}, midCoords, filterBottomLeft)
	sets := make([][]Coords, 4)
	sets[0] = tlSet
	sets[1] = brSet
	sets[2] = trSet
	sets[3] = blSet

	areaMax := "0"
	for sI := 0; sI < len(sets); sI++ {
		for sJ := sI; sJ < len(sets); sJ++ {
			firstSet := sets[sI]
			secondSet := sets[sJ]

			for i := 0; i < len(firstSet); i++ {
				for j := 0; j < len(secondSet); j++ {
					firstCoords := firstSet[i]
					secondCoords := secondSet[j]

					var smallX, bigX string
					if checkIfSmaller(firstCoords.x, secondCoords.x) {
						smallX = firstCoords.x
						bigX = secondCoords.x
					} else {
						smallX = secondCoords.x
						bigX = firstCoords.x
					}

					var smallY, bigY string
					if checkIfSmaller(firstCoords.y, secondCoords.y) {
						smallY = firstCoords.y
						bigY = secondCoords.y
					} else {
						smallY = secondCoords.y
						bigY = firstCoords.y
					}

					area := calcArea(smallX, smallY, bigX, bigY)
					if checkIfSmaller(areaMax, area) {
						areaMax = area
					}
				}
			}
		}
	}

	return areaMax
}

// NOTE: 1 needs to be lower then 2
func calcArea(x1, y1, x2, y2 string) string {
	xDiff := substract(x1, x2)
	xDiff = increment(xDiff)
	yDiff := substract(y1, y2)
	yDiff = increment(yDiff)
	return multiply(xDiff, yDiff)
}

func filterToCorner(coords []Coords, corner Coords, origMid Coords, filterFn filterCoordsFn) []Coords {

	mid := Coords{
		x: origMid.x,
		y: origMid.y,
	}
	coordsToCheck := coords
	for {
		var tempCoords []Coords
		tempCoords, mid = filterFn(corner, mid, coordsToCheck)
		if len(tempCoords) == 0 {
			break
		}
		coordsToCheck = tempCoords
		if len(tempCoords) == 1 {
			break
		}
	}

	//Point 1 - X:  82919 , Y:  85891
	//Point 2 - X:  12644 , Y:  18423

	return coordsToCheck
}

type filterCoordsFn func(Coords, Coords, []Coords) ([]Coords, Coords)

/*
func findClosestCoords(coords []Coords, to Coords) (Coords, bool) {

		minDist := distanceSquared(coords[0].x, coords[0].y, to.x, to.y)
		theCoords := coords[0]
		twins := false

		for i := 1; i < len(coords); i++ {
			tempDist := distanceSquared(coords[i].x, coords[i].y, to.x, to.y)
			if checkIfSmaller(tempDist, minDist) {
				minDist = tempDist
				theCoords = coords[i]
				twins = false
			} else if tempDist == minDist {
				// NOTE: if equal then there is a mirror point. need to check!
				twins = true
			}
		}

		return theCoords, twins
	}
*/
func calcMidCoord(xMinMax, yMinMax MinMax) Coords {

	calcMid := func(minmax MinMax) string {
		diff := substract(minmax.min, minmax.max)
		diffHalf, _ := divide(diff, "2")
		return add(minmax.min, diffHalf)
	}
	xMid := calcMid(xMinMax)
	yMid := calcMid(yMinMax)

	return Coords{xMid, yMid}
}

func filterTopLeft(corner Coords, mid Coords, coords []Coords) ([]Coords, Coords) {

	mid.x, _ = divide(mid.x, "2")
	mid.y, _ = divide(mid.y, "2")

	tl := make([]Coords, 0, len(coords)/4)
	for i := 0; i < len(coords); i++ {
		if checkIfSmaller(coords[i].x, mid.x) {
			if checkIfSmaller(coords[i].y, mid.y) {
				tl = append(tl, coords[i])
			}
		}
	}

	return tl, mid
}
func filterTopRight(corner Coords, mid Coords, coords []Coords) ([]Coords, Coords) {

	mid.x, _ = divide(mid.x, "2")
	mid.y, _ = divide(mid.y, "2")

	tempMid := Coords{}
	tempMid.x = substract(mid.x, corner.x)
	tempMid.y = mid.y

	tr := make([]Coords, 0, len(coords)/4)
	for i := 0; i < len(coords); i++ {
		if !checkIfSmaller(coords[i].x, tempMid.x) {
			if checkIfSmaller(coords[i].y, tempMid.y) {
				tr = append(tr, coords[i])
			}
		}
	}
	return tr, mid
}
func filterBottomLeft(corner Coords, mid Coords, coords []Coords) ([]Coords, Coords) {

	mid.x, _ = divide(mid.x, "2")
	mid.y, _ = divide(mid.y, "2")

	tempMid := Coords{}
	tempMid.x = mid.x
	tempMid.y = substract(mid.y, corner.y)

	bl := make([]Coords, 0, len(coords)/4)
	for i := 0; i < len(coords); i++ {
		if checkIfSmaller(coords[i].x, tempMid.x) {
			if !checkIfSmaller(coords[i].y, tempMid.y) {
				bl = append(bl, coords[i])
			}
		}
	}
	return bl, mid
}
func filterBottomRight(corner Coords, mid Coords, coords []Coords) ([]Coords, Coords) {

	mid.x, _ = divide(mid.x, "2")
	mid.y, _ = divide(mid.y, "2")

	tempMid := Coords{}
	tempMid.x = substract(mid.x, corner.x)
	tempMid.y = substract(mid.y, corner.y)

	br := make([]Coords, 0, len(coords)/4)
	for i := 0; i < len(coords); i++ {
		if !checkIfSmaller(coords[i].x, tempMid.x) {
			if !checkIfSmaller(coords[i].y, tempMid.y) {
				br = append(br, coords[i])
			}
		}
	}
	return br, mid
}

type Quadrants struct {
	tl, tr, bl, br []Coords
}

func sortCoords(mid Coords, coords []Coords) Quadrants {

	quad := Quadrants{}
	quad.tl = make([]Coords, 0, len(coords)/4)
	quad.tr = make([]Coords, 0, len(coords)/4)
	quad.bl = make([]Coords, 0, len(coords)/4)
	quad.br = make([]Coords, 0, len(coords)/4)

	for i := 0; i < len(coords); i++ {
		if checkIfSmaller(coords[i].x, mid.x) {
			if checkIfSmaller(coords[i].y, mid.y) {
				quad.tl = append(quad.tl, coords[i])
			} else {
				quad.bl = append(quad.bl, coords[i])
			}
		} else {
			if checkIfSmaller(coords[i].y, mid.y) {
				quad.tr = append(quad.tr, coords[i])
			} else {
				quad.br = append(quad.br, coords[i])
			}
		}
	}

	return quad
}

type MinMax struct {
	min, max string
}

func findMinMax(coords []Coords) (MinMax, MinMax) {

	xMin := coords[0].x
	yMin := coords[0].y
	xMax := coords[0].x
	yMax := coords[0].y
	for i := 1; i < len(coords); i++ {

		if checkIfSmaller(coords[i].x, xMin) {
			xMin = coords[i].x
		}
		if checkIfSmaller(coords[i].y, yMin) {
			yMin = coords[i].y
		}

		if coords[i].x != xMax && !checkIfSmaller(coords[i].x, xMax) {
			xMax = coords[i].x
		}
		if coords[i].y != yMax && !checkIfSmaller(coords[i].y, yMax) {
			yMax = coords[i].y
		}
	}

	return MinMax{xMin, xMax}, MinMax{yMin, yMax}
}

func distanceSquared(x1, y1, x2, y2 string) string {

	calcDiffSqr := func(c1, c2 string) string {
		if c1 != c2 {
			var diff string
			if checkIfSmaller(c1, c2) {
				diff = substract(c1, c2)
			} else {
				diff = substract(c2, c1)
			}
			return multiply(diff, diff)
		} else {
			return "0"
		}
	}

	xDiffSqr := calcDiffSqr(x1, x2)
	yDiffSqr := calcDiffSqr(y1, y2)

	return add(xDiffSqr, yDiffSqr)
}

// divide divides 'this with 'with' on the actual strings
// without converting to integers.
// NOTE: won't validate numbers and is not handle negative
// numbers or remainders.
func divide(this, with string) (string, string) {

	dividendDigitCount := len(this)
	divisorDigitCount := len(with)

	quotient := "0"
	if dividendDigitCount < divisorDigitCount {
		return quotient, strings.Clone(this)
	}

	currPos := divisorDigitCount
	dividend := this[:currPos]
	for {
		prevDivisor := ""
		divisor := strings.Clone(with)
		divCount := "0"
		for {
			if checkIfSmaller(dividend, divisor) {
				divisor = prevDivisor
				break
			}
			divCount = increment(divCount)
			if dividend == divisor {
				break
			}
			prevDivisor = divisor
			divisor = add(with, divisor)
		}

		quotient += divCount
		if divCount != "0" {
			dividend = substract(divisor, dividend)
		}

		if currPos == dividendDigitCount {
			break
		}
		if dividend == "0" {
			dividend = this[currPos : currPos+1]
		} else {
			dividend += this[currPos : currPos+1]
		}
		currPos++
	}

	return strings.TrimLeft(quotient, "0"), dividend
}

func multiply(this, with string) string {

	if with == "0" {
		return "0"
	}

	newNum := strings.Clone(this)
	for {

		if with == "1" {
			break
		}

		newNum = add(newNum, this)
		with = decrement(with)
	}

	return newNum
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
