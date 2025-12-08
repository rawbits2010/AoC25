package main

import (
	"fmt"
	"log"
	"maps"
	"sort"
	"strconv"
	"strings"

	"github.com/rawbits2010/AoC25/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	coords, err := readCoords(lines)
	if err != nil {
		log.Fatalf("error reading coords: %s", err)
	}

	distances := calculateDistances(coords)

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].value < distances[j].value
	})

	//circuits := letsMakeContact(distances, len(coords), 40)
	circuits := letsMakeContact(distances, len(coords), 1000)
	/*
		for _, c := range circuits {
			c.Print()
		}
	*/
	sort.Slice(circuits, func(i, j int) bool {
		return circuits[i].Count() > circuits[j].Count()
	})

	resultP1 := circuits[0].Count() * circuits[1].Count() * circuits[2].Count()

	lastConnDist := letsMakeContactP2(distances, len(coords))
	resultP2 := coords[lastConnDist.b1Idx].x * coords[lastConnDist.b2Idx].x

	fmt.Printf("Result - Part 1: %d, Part 2: %d\n", resultP1, resultP2)
}

func letsMakeContactP2(distances []Distance, numBoxes int) Distance {

	circuits := make([]*Circuit, numBoxes)
	for i := 0; i < len(circuits); i++ {
		circuits[i] = NewCircuit()
		circuits[i].Add(i)
	}

	lastConnIdx := 0
	for dIdx, dist := range distances {

		c1Idx := 0
		for ; c1Idx < len(circuits); c1Idx++ {
			if circuits[c1Idx].Contains(dist.b1Idx) {
				break
			}
		}
		c2Idx := 0
		for ; c2Idx < len(circuits); c2Idx++ {
			if circuits[c2Idx].Contains(dist.b2Idx) {
				break
			}
		}

		if c1Idx == c2Idx {
			continue
		}

		lastConnIdx = dIdx

		tempBoxes := circuits[c2Idx].Members()
		for _, box := range tempBoxes {
			circuits[c1Idx].Add(box)
		}

		circuits = append(circuits[:c2Idx], circuits[c2Idx+1:]...)
	}

	return distances[lastConnIdx]
}

type Circuit struct {
	boxes map[int]bool
}

func NewCircuit() *Circuit {
	return &Circuit{
		boxes: make(map[int]bool),
	}
}

func (c Circuit) Contains(idx int) bool {
	_, ok := c.boxes[idx]
	return ok
}

func (c *Circuit) Add(idx int) {
	c.boxes[idx] = true
}

func (c Circuit) Count() int {
	return len(c.boxes)
}

func (c Circuit) Members() []int {
	temp := make([]int, len(c.boxes))
	i := 0
	for val := range maps.Keys(c.boxes) {
		temp[i] = val
		i++
	}
	return temp
}

func (c Circuit) Print() {
	for val := range maps.Keys(c.boxes) {
		fmt.Printf("%d, ", val)
	}
	fmt.Println()
}

func letsMakeContact(distances []Distance, numBoxes int, limit int) []*Circuit {

	circuits := make([]*Circuit, numBoxes)
	for i := 0; i < len(circuits); i++ {
		circuits[i] = NewCircuit()
		circuits[i].Add(i)
	}
	for dIdx, dist := range distances {
		_ = dIdx

		// NOTE: took me 3 hours to notice this part in the
		// puzzle description *sigh*
		if dIdx == limit {
			break
		}

		c1Idx := 0
		for ; c1Idx < len(circuits); c1Idx++ {
			if circuits[c1Idx].Contains(dist.b1Idx) {
				break
			}
		}
		c2Idx := 0
		for ; c2Idx < len(circuits); c2Idx++ {
			if circuits[c2Idx].Contains(dist.b2Idx) {
				break
			}
		}

		if c1Idx == c2Idx {
			continue
		}

		tempBoxes := circuits[c2Idx].Members()
		for _, box := range tempBoxes {
			circuits[c1Idx].Add(box)
		}

		circuits = append(circuits[:c2Idx], circuits[c2Idx+1:]...)
	}

	return circuits
}

func calculateDistances(coords []Coords) []Distance {

	distances := make([]Distance, 0, len(coords)*(len(coords)-1))
	for c1Idx := 0; c1Idx < len(coords); c1Idx++ {
		for c2Idx := c1Idx + 1; c2Idx < len(coords); c2Idx++ {
			dTemp := distanceSquared(coords[c1Idx], coords[c2Idx])
			distances = append(distances, Distance{c1Idx, c2Idx, dTemp})
		}
	}

	return distances
}

type Distance struct {
	b1Idx, b2Idx int
	value        int
}

func readCoords(lines []string) ([]Coords, error) {

	coords := make([]Coords, len(lines))
	for lineIdx, line := range lines {

		nums := strings.Split(line, ",")
		if len(nums) != 3 {
			return nil, fmt.Errorf("malformed input line at %d (%s)", lineIdx, line)
		}

		x, err := strconv.Atoi(nums[0])
		if err != nil {
			return nil, fmt.Errorf("invalid number in line %d (%s)", lineIdx, nums[0])
		}

		y, err := strconv.Atoi(nums[1])
		if err != nil {
			return nil, fmt.Errorf("invalid number in line %d (%s)", lineIdx, nums[1])
		}

		z, err := strconv.Atoi(nums[2])
		if err != nil {
			return nil, fmt.Errorf("invalid number in line %d (%s)", lineIdx, nums[2])
		}

		coords[lineIdx] = Coords{x, y, z}
	}

	return coords, nil
}

type Coords struct {
	x, y, z int
}

func distanceSquared(box1, box2 Coords) int {
	return (box1.x-box2.x)*(box1.x-box2.x) + (box1.y-box2.y)*(box1.y-box2.y) + (box1.z-box2.z)*(box1.z-box2.z)
}
