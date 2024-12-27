package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"strings"
)

func main() {
	content, _ := file.GetContent("../data.txt")
	f := 0
	for _, df := range strings.Split(string(content), "\n") {
		f += common.StrToA(df)
	}
	fmt.Println(f)
}
