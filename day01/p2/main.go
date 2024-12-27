package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"strings"
)

func main() {
	content, _ := file.GetContent("../data.txt")

	fm := make(map[int]int)
	dfs := make([]int, 0, strings.Count(string(content), "\n")+1)
	for _, df := range strings.Split(string(content), "\n") {
		dfs = append(dfs, common.StrToA(string(df)))
	}

	f, i := 0, 0
	for {
		f += dfs[i]
		if fv, e := fm[f]; e {
			fm[f] = fv + 1
		} else {
			fm[f] = 1
		}
		if fm[f] == 2 {
			break
		}
		if i == len(dfs)-1 {
			i = 0
		} else {
			i++
		}
	}

	fmt.Println(f)
}
