package main

import (
	"fmt"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	countPart1, countPart2 := findMovableRolls(lines)

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", countPart1, countPart2)
}

const Roll = '@'
const Empty = '.'

func findMovableRolls(lines []string) (int, int) {

	countPart1 := -1
	var countPart2 int
	for {
		var removedCount int
		lines, removedCount = removeAccessible(lines)

		if countPart1 < 0 {
			countPart1 = removedCount
		}

		countPart2 += removedCount

		if removedCount == 0 {
			break
		}
	}

	return countPart1, countPart2
}

func removeAccessible(lines []string) ([]string, int) {

	tempLines := make([]string, len(lines))

	rollCount := 0
	for i := 0; i < len(lines); i++ {

		tempBytes := make([]byte, len(lines[i]))
		for j := 0; j < len(lines[i]); j++ {

			if lines[i][j] != Roll {
				tempBytes[j] = lines[i][j]
				continue
			}

			count := 0
			if i > 0 && j > 0 && lines[i-1][j-1] == Roll {
				count++
			}
			if i > 0 && lines[i-1][j] == Roll {
				count++
			}
			if i > 0 && j < len(lines[i])-1 && lines[i-1][j+1] == Roll {
				count++
			}
			if j > 0 && lines[i][j-1] == Roll {
				count++
			}
			if j < len(lines[i])-1 && lines[i][j+1] == Roll {
				count++
			}
			if i < len(lines)-1 && j > 0 && lines[i+1][j-1] == Roll {
				count++
			}
			if i < len(lines)-1 && lines[i+1][j] == Roll {
				count++
			}
			if i < len(lines)-1 && j < len(lines[i])-1 && lines[i+1][j+1] == Roll {
				count++
			}

			if count < 4 {
				tempBytes[j] = Empty
				rollCount++
			} else {
				tempBytes[j] = Roll
			}

		}

		tempLines[i] = string(tempBytes)
	}

	return tempLines, rollCount
}
