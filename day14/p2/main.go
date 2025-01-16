package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"strconv"
)

type Recipes struct {
	scores []byte
	elves  []int
	//lengths []int
}

func NewRecipes() *Recipes {
	r := &Recipes{}
	r.scores = []byte{3, 7}
	r.elves = []int{0, 1}
	//r.lengths = []int{2}
	return r
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}

func (r *Recipes) Len() int {
	return len(r.scores)
}

func (r *Recipes) Last(i int) string {
	bs := make([]byte, 0, i)
	for s := len(r.scores) - i; s < len(r.scores); s++ {
		bs = append(bs, byte(r.scores[s])+'0')
	}
	return string(bs)
}

func (r *Recipes) CreateNew() int {
	created := 0
	elfValues := make([]byte, len(r.elves))
	for i, e := range r.elves {
		elfValues[i] = r.scores[e]
	}
	pairs := common.GetPairSets(elfValues)
	for _, pair := range pairs {
		val := fmt.Sprintf("%d", pair[0]+pair[1])
		r.scores = append(r.scores, byte(atoi(string(val[0]))))
		created++
		if len(val) > 1 {
			r.scores = append(r.scores, byte(atoi(string(val[1]))))
			created++
		}
	}
	for e := range r.elves {
		advance := int(elfValues[e] + 1)
		for a := 0; a < advance; a++ {
			i := r.elves[e] + 1
			if i == len(r.scores) {
				i = 0
			}
			r.elves[e] = i
		}
	}
	//r.lengths = append(r.lengths, len(r.scores))
	return created
}

func main() {
	recipes := NewRecipes()

	created := recipes.Len()
	stopAfter := 503761
	needAfter := 10
	for created < stopAfter {
		created += recipes.CreateNew()
	}

	curLen := recipes.Len()

	l := curLen + needAfter - (created - stopAfter)
	for recipes.Len() < l {
		recipes.CreateNew()
	}

	fmt.Println(recipes.Last(needAfter + (recipes.Len() - l))[0:needAfter])
}
