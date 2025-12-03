package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	sumPart1 := 0
	sumPart2 := 0
	for _, line := range lines {

		numStr := joltageSearch(line, 2)

		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.Fatalf("error converting jolts (%s): %s", numStr, err)
		}

		sumPart1 += num

		numStr = joltageSearch(line, 12)

		num, err = strconv.Atoi(numStr)
		if err != nil {
			log.Fatalf("error converting jolts (%s): %s", numStr, err)
		}

		sumPart2 += num
	}

	fmt.Printf("Result - Part 1:  %d, Part 2: %d\n", sumPart1, sumPart2)
}

func joltageSearch(line string, digitCount int) string {

	nums := make([]byte, digitCount)
	idx := -1
	for i := 1; i <= digitCount; i++ {

		num, newIdx := findLargestNumber(line[idx+1 : len(line)-(digitCount-i)])
		idx = newIdx + idx + 1

		nums[i-1] = num
	}

	return string(nums)
}

func findLargestNumber(bank string) (byte, int) {
	max := byte(0)
	idx := 0
	for i := 0; i < len(bank); i++ {
		if bank[i] > max {
			max = bank[i]
			idx = i
		}
	}

	//fmt.Println(bank, "-", max, "/", idx)

	return max, idx
}
