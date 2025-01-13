package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"strings"
)

func main() {
	ss := getData("../data.txt")

	pairs := common.GetPairSets(ss)

nextPair:
	for _, pair := range pairs {
		if len(pair[0]) == len(pair[1]) {
			var diff []int
			for i := 0; i < len(pair[0]); i++ {
				if pair[0][i] != pair[1][i] {
					diff = append(diff, i)
				}
				if len(diff) > 1 {
					continue nextPair
				}
			}
			if len(diff) == 1 {
				chars := pair[0][0:diff[0]] + pair[0][diff[0]+1:]
				fmt.Println(chars)
				break
			}
		}
	}

}

func getData(filename string) []string {
	lines, _ := file.GetLines(filename)
	ss := make([]string, len(lines))
	for i, line := range lines {
		tokens := strings.Fields(line)
		ss[i] = tokens[0]
	}
	return ss
}
