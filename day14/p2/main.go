package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"strconv"
	"strings"
)

type Recipes struct {
	scores []byte
	elves  []int
}

func NewRecipes() *Recipes {
	r := &Recipes{}
	r.scores = []byte{3, 7}
	r.elves = []int{0, 1}
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
	if len(r.scores) > i {
		bs := make([]byte, 0, i)
		for s := len(r.scores) - i; s < len(r.scores); s++ {
			bs = append(bs, byte(r.scores[s])+'0')
		}
		return string(bs)
	}
	return string(r.scores)
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
	return created
}

func main() {
	recipes := NewRecipes()

	created := recipes.Len()

	lookingFor := "503761"

	for {
		created += recipes.CreateNew()
		last := recipes.Last(len(lookingFor) + 1)
		if index := strings.Index(last, lookingFor); index != -1 {
			recipesBefore := recipes.Len() - len(lookingFor)
			if index == 0 {
				recipesBefore -= 1
			}
			fmt.Println(recipesBefore)
			break
		}
	}

}
