package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strings"
)

var (
	reInitialState = regexp.MustCompile(`initial state: ((?:\.|#)+)`)
	reProduction   = regexp.MustCompile(`((?:\.|#){5})\s+=>\s+(\.|#)`)
)

const (
	PlantChar   = '#'
	NoPlantChar = '.'
)

type Sequence string
type ProductionMap map[Sequence]bool

type Row struct {
	plants []byte
	origin int
}

func (r *Row) adjustRow() {
	fp := strings.Index(string(r.plants), string(PlantChar))
	lp := strings.LastIndex(string(r.plants), string(PlantChar))
	plantLen := lp - fp + 1
	adjusted := make([]byte, plantLen+8)
	for i := range adjusted {
		adjusted[i] = NoPlantChar
	}
	copy(adjusted[4:], r.plants[fp:])
	r.plants = adjusted

	r.origin += 4 - fp
}

func (r *Row) String() string {
	return string(r.plants)
}

func (r *Row) PlantCount() int {
	count := 0
	for _, p := range r.plants {
		if p == PlantChar {
			count++
		}
	}
	return count
}

func (r *Row) PotsWithPlants() []int {
	ids := make([]int, 0, r.PlantCount())
	for p := range r.plants {
		if r.plants[p] == PlantChar {
			ids = append(ids, p-r.origin)
		}
	}
	return ids
}

func (r *Row) Generate(pm ProductionMap) {
	r.adjustRow()
	ng := make([]byte, len(r.plants))
	ng[0], ng[1], ng[len(ng)-1], ng[len(ng)-2] = NoPlantChar, NoPlantChar, NoPlantChar, NoPlantChar
	for i := 2; i < len(r.plants)-2; i++ {
		seq := Sequence(r.plants[i-2 : i+3])
		if pm[seq] {
			ng[i] = PlantChar
		} else {
			ng[i] = NoPlantChar
		}
	}
	r.plants = ng
}

func main() {
	pm, row := getData("../data.txt")

	generations := make(map[string]int)
	origins := make(map[string]int)

	numGenerations := 50000000000
	skipped := false

	for i := 0; i < numGenerations; i++ {
		generation := row.String()
		if _, e := generations[generation]; e && !skipped {
			lastOrigin := origins[generation]
			currentOrigin := row.origin
			lastGen := generations[generation]
			currentGen := i
			if lastOrigin-1 == currentOrigin && lastGen+1 == currentGen {
				// at this point it's just going to keep repeating and shifting to the right

				left := numGenerations - currentGen - 10
				row.origin -= left
				i += left

				skipped = true
			}
		}
		origins[generation] = row.origin
		generations[generation] = i
		row.Generate(pm)
	}

	fmt.Println(row.PlantCount())
	pots := row.PotsWithPlants()
	fmt.Println(pots)

	answer := 0
	for _, p := range pots {
		answer += p
	}
	fmt.Println(answer)
}

func getData(filename string) (ProductionMap, *Row) {
	lines, _ := file.GetLines(filename)
	matches := reInitialState.FindStringSubmatch(lines[0])
	row := &Row{plants: []byte(matches[1])}
	pm := make(ProductionMap)
	for _, line := range lines[2:] {
		matches = reProduction.FindStringSubmatch(line)
		plant := false
		if matches[2][0] == PlantChar {
			plant = true
		}
		pm[Sequence(matches[1])] = plant

	}
	return pm, row
}
