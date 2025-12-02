package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	countPasses := uint(0)
	countZeros := uint(0)
	currPos := uint(50)
	for _, line := range lines {

		dir, amount, err := parseLine(line)
		if err != nil {
			log.Fatalf("error parsing line (%s): %s", line, err)
		}

		countPasses += amount / 100

		amount %= 100

		prevPos := currPos
		switch dir {
		case 'R':
			currPos = (currPos + amount) % 100

			if prevPos > currPos && prevPos != 0 && currPos != 0 {
				countPasses++
			}

		case 'L':
			if amount > currPos {
				currPos = 100 - (amount - currPos)
			} else {
				currPos -= amount
			}

			if prevPos < currPos && prevPos != 0 && currPos != 0 {
				countPasses++
			}
		}

		if currPos == 0 {
			countZeros++
		}

		fmt.Printf("dir: %s, amount: %d, prevPos: %d, curPos: %d, pass: %d\n", string(dir), amount, prevPos, currPos, countPasses)
	}

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", countZeros, countPasses+countZeros)
}

type Direction byte

const (
	Left  Direction = 'L'
	Right Direction = 'R'
)

func parseLine(line string) (Direction, uint, error) {

	if len(line) < 2 {
		return 0, 0, fmt.Errorf("malformed input")
	}

	dir := line[0]
	if dir != 'L' && dir != 'R' {
		return 0, 0, fmt.Errorf("invalid direction: %s", string(dir))
	}

	amount, err := strconv.ParseUint(line[1:], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid number (%s): %w", line[1:], err)
	}

	return Direction(dir), uint(amount), nil
}
