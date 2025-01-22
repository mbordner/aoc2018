package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"strconv"
)

var (
	reNanoBot = regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>,\s+r=(\d+)`)
)

type IntNumber interface {
	int | int32 | int64
}

type Vector[T IntNumber] struct {
	X, Y, Z T
}

func (v Vector[T]) String() string {
	return fmt.Sprintf("{%d, %d, %d}", v.X, v.Y, v.Z)
}

func (v Vector[T]) Add(o Vector[T]) Vector[T] {
	return Vector[T]{X: v.X + o.X, Y: v.Y + o.Y, Z: v.Z + o.Z}
}

func (v Vector[T]) Dis(o Vector[T]) T {
	return abs(o.X-v.X) + abs(o.Y-v.Y) + abs(o.Z-v.Z)
}

type NanoBots[T IntNumber] []*NanoBot[T]
type NanoBot[T IntNumber] struct {
	pos    Vector[T]
	radius T
}

func main() {
	nanobots, strongest := getData[int64]("../data.txt")
	inRange := 0
	for _, bot := range nanobots {
		if strongest.pos.Dis(bot.pos) <= strongest.radius {
			inRange++
		}
	}
	fmt.Printf("nanobots in range: %d\n", inRange)
}

func getData[T IntNumber](filename string) (NanoBots[T], *NanoBot[T]) {
	lines, _ := file.GetLines(filename)
	bots := make(NanoBots[T], len(lines))
	maxSignalRadius := T(0)
	var strongest *NanoBot[T]
	for i, line := range lines {
		matches := reNanoBot.FindStringSubmatch(line)
		bots[i] = &NanoBot[T]{pos: Vector[T]{
			X: atoi[T](matches[1]),
			Y: atoi[T](matches[2]),
			Z: atoi[T](matches[3]),
		}, radius: atoi[T](matches[4])}

		if bots[i].radius > maxSignalRadius {
			maxSignalRadius = bots[i].radius
			strongest = bots[i]
		}
	}
	return bots, strongest
}

func atoi[T IntNumber](s string) T {
	val, _ := strconv.ParseInt(s, 10, 64)
	return T(val)
}

func abs[T IntNumber](x T) T {
	if x >= 0 {
		return x
	}
	return -x
}
