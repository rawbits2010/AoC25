package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	splitterHitCount, err := processManifold(lines)
	if err != nil {
		log.Fatalf("error processing part 1: %s\n", err)
	}
	/*
		timelines, err := DFS(lines)
		if err != nil {
			log.Fatalf("error processing part 2: %s\n", err)
		}
	*/

	startIdx, err := getStartIdx(lines[0])
	if err != nil {
		log.Fatal(err)
	}
	timelines := timelinesCounter(lines, startIdx)
	sumTimelines := 0
	for _, beamCount := range timelines {
		sumTimelines += beamCount
	}

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", splitterHitCount, sumTimelines)
}

//
// Part 2 - take 2

func timelinesCounter(lines []string, startIdx int) []int {
	beamCounters := make([]int, len(lines[0]))
	beamCounters[startIdx] = 1
	for _, currLine := range lines {

		newBeamCounters := make([]int, len(lines[0]))
		for beamCounterIdx, beamCounter := range beamCounters {
			if beamCounter == 0 {
				continue
			}

			if currLine[beamCounterIdx] == Splitter {
				newBeamCounters[beamCounterIdx-1] += beamCounter
				newBeamCounters[beamCounterIdx+1] += beamCounter
			} else {
				newBeamCounters[beamCounterIdx] += beamCounter
			}
		}

		beamCounters = newBeamCounters
	}
	return beamCounters
}

//
// Part 2 - fail

type SplitterNode struct {
	lineIdx, charIdx    int
	leftBeam, rightBeam bool
}

func DFS(lines []string) (int, error) {

	if len(lines) == 0 {
		return 0, fmt.Errorf("no manifold area provided")
	}

	route := make([]SplitterNode, 0)

	startIdx, err := getStartIdx(lines[0])
	if err != nil {
		return 0, fmt.Errorf("no beam start (%s) present", string(Start))
	}

	timelines := 0
	for {
		var beamIdx int
		var currLineIdx int
		if len(route) == 0 {
			beamIdx = startIdx
			currLineIdx = 0
		} else {
			route, beamIdx, timelines = pickCurrentNode(route, timelines, len(lines[0]))
			if route == nil {
				break // hit the end
			}
			currNode := route[len(route)-1]
			currLineIdx = currNode.lineIdx
		}

		nextNode, isBeamExited, err := findNextNode(beamIdx, currLineIdx+1, lines)
		if err != nil {
			return 0, err
		}

		if isBeamExited {
			timelines++
			continue
		}

		route = append(route, nextNode)
	}

	return timelines, nil
}

func pickCurrentNode(route []SplitterNode, timelines int, manifoldWidth int) ([]SplitterNode, int, int) {

	var beamIdx int
	for {

		if len(route) == 0 {
			return nil, 0, timelines
		}

		currNode := &route[len(route)-1]
		if !currNode.leftBeam {

			beamIdx = currNode.charIdx - 1
			currNode.leftBeam = true

			if beamIdx < 0 {
				timelines++
				continue
			}
			break

		} else if !currNode.rightBeam {

			beamIdx = currNode.charIdx + 1
			currNode.rightBeam = true

			if beamIdx >= manifoldWidth {
				timelines++
				continue
			}
			break

		} else {
			route = route[:len(route)-1]
		}
	}

	return route, beamIdx, timelines
}

func findNextNode(beamIdx int, fromLineIdx int, lines []string) (SplitterNode, bool, error) {

	for lineIdx := fromLineIdx; lineIdx < len(lines); lineIdx++ {

		beamEndIdx, splitterCount, err := processSplitters([]int{beamIdx}, lines[lineIdx])
		if err != nil {
			return SplitterNode{}, false, fmt.Errorf("error processing manifold line %d", lineIdx)
		}

		if splitterCount == 0 {
			continue
		} else if splitterCount != 1 {
			return SplitterNode{}, false, fmt.Errorf("one tachyon beam can only hit ONE splitter at once")
		}

		nextNode := SplitterNode{
			lineIdx:   lineIdx,
			charIdx:   beamEndIdx[0] + 1,
			leftBeam:  false,
			rightBeam: false,
		}

		return nextNode, false, nil
	}

	return SplitterNode{}, true, nil
}

//
// Part 1

func processManifold(lines []string) (int, error) {

	if len(lines) == 0 {
		return 0, fmt.Errorf("no manifold area provided")
	}

	startIdx, err := getStartIdx(lines[0])
	if err != nil {
		return 0, err
	}

	splitterHitCount := 0
	beamEndIdxs := []int{startIdx}
	for lineIdx := 1; lineIdx < len(lines); lineIdx++ {

		var splittersHit int
		beamEndIdxs, splittersHit, err = processSplitters(beamEndIdxs, lines[lineIdx])
		if err != nil {
			return 0, fmt.Errorf("error processing line at %d", lineIdx)
		}

		splitterHitCount += splittersHit

	}

	return splitterHitCount, nil
}

const Start = 'S'
const Splitter = '^'

func getStartIdx(line string) (int, error) {
	idx := strings.Index(line, string(Start))
	if idx == -1 {
		return 0, fmt.Errorf("no start (%s) mark is present", string(Start))
	}
	return idx, nil
}

func processSplitters(beamEndIdx []int, line string) ([]int, int, error) {

	newBeamEndIdx := make([]int, 0, len(line))
	splitterHitCount := 0
	for _, beamIdx := range beamEndIdx {

		if len(line) <= beamIdx {
			return nil, 0, fmt.Errorf("invalid manifold width %d - tryed to check %d", len(line), beamIdx)
		}

		if line[beamIdx] == Splitter {
			splitterHitCount++

			leftBeamIdx := beamIdx - 1
			rightBeamIdx := beamIdx + 1

			if leftBeamIdx >= 0 {
				if len(newBeamEndIdx) > 0 {
					if newBeamEndIdx[len(newBeamEndIdx)-1] != leftBeamIdx {
						newBeamEndIdx = append(newBeamEndIdx, leftBeamIdx)
					}
				} else {
					newBeamEndIdx = append(newBeamEndIdx, leftBeamIdx)
				}
			}

			if rightBeamIdx < len(line) {
				newBeamEndIdx = append(newBeamEndIdx, rightBeamIdx)
			}
		} else {
			if len(newBeamEndIdx) > 0 {
				if newBeamEndIdx[len(newBeamEndIdx)-1] != beamIdx {
					newBeamEndIdx = append(newBeamEndIdx, beamIdx)
				}
			} else {
				newBeamEndIdx = append(newBeamEndIdx, beamIdx)
			}
		}

	}

	return newBeamEndIdx, splitterHitCount, nil
}
