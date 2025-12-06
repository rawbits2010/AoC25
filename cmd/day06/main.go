package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	// Part 1

	aochs, err := createAOChs(lines[len(lines)-1])
	if err != nil {
		log.Fatalf("error parsing operations line: %s", err)
	}

	aochs, err = doTheMath(lines[:len(lines)-1], aochs)
	if err != nil {
		log.Fatalf("error processing operations: %s", err)
	}

	resultPart1 := sumResults(aochs)

	// Part 2
	aochs, err = createAOChsP2(lines[len(lines)-1])
	if err != nil {
		log.Fatalf("error parsing operations line: %s", err)
	}

	aochs, err = doTheMathP2(lines[:len(lines)-1], aochs)
	if err != nil {
		log.Fatalf("error processing operations: %s", err)
	}

	resultPart2 := sumResults(aochs)

	fmt.Printf("Result - Part 1: %s, Part 2: %s\n", resultPart1, resultPart2)
}

func sumResults(aochs []AOCh) string {
	sum := "0"
	for _, chain := range aochs {
		sum = add(sum, chain.Accumulator)
	}
	return sum
}

//
//-Part 2----

func doTheMathP2(valueLines []string, aochs []AOCh) ([]AOCh, error) {

	expOperandSize := len(valueLines)

	if expOperandSize == 0 {
		return nil, fmt.Errorf("no value lines provided")
	}

	for i := 0; i < len(aochs); i++ {
		for j := 0; j < len(aochs[i].Operands); j++ {
			aochs[i].Operands[j] = make([]byte, expOperandSize)
		}
	}

	// read the operands
	for lineIdx, line := range valueLines {
		for i := 0; i < len(aochs); i++ {

			expNumDigits := len(aochs[i].Operands)

			if len(line) < expNumDigits {
				return nil, fmt.Errorf("malformed value line at %d - mismatched number of values (at val idx %d)", lineIdx, i)
			}

			val := line[:expNumDigits]
			line = line[expNumDigits:]

			for dIdx := 0; dIdx < expNumDigits; dIdx++ {
				aochs[i].Operands[dIdx][lineIdx] = val[dIdx]
			}

			if len(line) > 0 {
				line = line[1:] // +1 whitespace separator
			}
		}
		if len(line) > 0 {
			log.Printf("WARNING: malformed value line - extra characters in line at %d (%s)\n", lineIdx, line)
		}
	}

	cleanStr := func(val []byte) string {
		numStr := string(val)
		numStr = strings.TrimSpace(numStr)
		return numStr
	}

	// do math
	for i := 0; i < len(aochs); i++ {

		for opsIdx, operand := range aochs[i].Operands {

			val := cleanStr(operand)

			if opsIdx == 0 {
				aochs[i].Accumulator = val
				continue
			}

			switch aochs[i].Op {
			case Addition:
				aochs[i].Accumulator = add(aochs[i].Accumulator, val)
			case Multiplication:
				aochs[i].Accumulator = multiply(aochs[i].Accumulator, val)
			default:
				return nil, fmt.Errorf("unsupported operation found (%s)", aochs[i].Op)
			}
		}
	}

	return aochs, nil
}

func createAOChsP2(opsLine string) ([]AOCh, error) {

	if len(opsLine) == 0 {
		return nil, fmt.Errorf("no operatoions found")
	}

	ops := make([]Operation, 0, 100)
	digitCount := make([]int, 0, 100)

	currCount := 0
	for _, char := range opsLine {
		switch string(char) {
		case string(Addition):
			ops = append(ops, Addition)
			digitCount = append(digitCount, currCount-1)
		case string(Multiplication):
			ops = append(ops, Multiplication)
			digitCount = append(digitCount, currCount-1)
		case " ", "\t":
			currCount++
			continue
		default:
			return nil, fmt.Errorf("invalid operation found (%s)", string(char))
		}
		currCount = 1
	}
	digitCount = append(digitCount, currCount)
	digitCount = digitCount[1:]

	aoch := make([]AOCh, len(ops))
	for opIdx, op := range ops {
		aoch[opIdx].Op = op
		aoch[opIdx].Operands = make([][]byte, digitCount[opIdx])
	}

	return aoch, nil
}

//
//-Part 1----

func doTheMath(valueLines []string, aochs []AOCh) ([]AOCh, error) {

	if len(valueLines) == 0 {
		return nil, fmt.Errorf("no value lines provided")
	}

	// filling in starting values
	values, err := parseValueLine(valueLines[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing line (%d): %w", 0, err)
	}

	if len(values) != len(aochs) {
		return nil, fmt.Errorf("mismatched number of values (%d to %d ops)", len(values), len(aochs))
	}

	for valIdx, val := range values {
		aochs[valIdx].Accumulator = val
	}

	// do operations line by line
	for i := 1; i < len(valueLines); i++ {

		values, err = parseValueLine(valueLines[i])
		if err != nil {
			return nil, fmt.Errorf("error parsing line (%d): %w", i, err)
		}

		if len(values) != len(aochs) {
			return nil, fmt.Errorf("mismatched number of values (%d to %d ops)", len(values), len(aochs))
		}

		for valIdx, val := range values {
			switch aochs[valIdx].Op {
			case Addition:
				aochs[valIdx].Accumulator = add(aochs[valIdx].Accumulator, val)
			case Multiplication:
				aochs[valIdx].Accumulator = multiply(aochs[valIdx].Accumulator, val)
			default:
				return nil, fmt.Errorf("unsupported operation found (%s)", aochs[valIdx].Op)
			}
		}

	}

	return aochs, nil
}

func parseValueLine(valueLine string) ([]string, error) {

	if len(valueLine) == 0 {
		return nil, fmt.Errorf("no values found")
	}

	values := strings.Fields(valueLine)

	// NOTE: you should check for malformed numbers here

	return values, nil
}

func createAOChs(opsLine string) ([]AOCh, error) {

	if len(opsLine) == 0 {
		return nil, fmt.Errorf("no operatoions found")
	}

	ops := strings.Fields(opsLine)

	aoch := make([]AOCh, len(ops))
	for opIdx, op := range ops {

		switch op {
		case string(Addition):
			aoch[opIdx].Op = Addition
		case string(Multiplication):
			aoch[opIdx].Op = Multiplication
		default:
			return nil, fmt.Errorf("invalid operation found (%s)", op)
		}
	}

	return aoch, nil
}

type Operation string

const (
	Addition       Operation = "+"
	Multiplication Operation = "*"
)

// Aritmetic Operation Chain
type AOCh struct {
	Op          Operation
	Accumulator string
	Operands    [][]byte
}

const zeroDigit = '0'

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
