package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"sort"
	"strings"
)

var (
	reStepDep = regexp.MustCompile(`Step (\w) must be finished before step (\w) can begin.`)
)

const (
	MaxWorkers = 5
	ExtraTime  = 60
)

func main() {
	steps := getSteps("../data.txt")

	executed := make([]string, 0, len(steps))

	workers := make(Working)

	secs := 0
	for !steps.AllFinished() {

		can := steps.CanStart()

		for i := 0; i < len(can) && workers.Len() < MaxWorkers; i++ {
			workers.Add(can[i], steps[can[i]])
			steps[can[i]].Start()
		}

		secs++
		steps.Tick(1)

		finished := make([]string, 0, len(workers))
		for id := range workers {
			if steps[id].IsFinished() {
				workers.Remove(id)
				finished = append(finished, id)
			}
		}
		if len(finished) > 0 {
			sort.Strings(finished)
			executed = append(executed, finished...)
		}
	}

	fmt.Println(strings.Join(executed, ""))
	fmt.Println(secs)
}

type Working map[string]*Step

func (w Working) Add(id string, step *Step) {
	w[id] = step
}

func (w Working) Remove(id string) {
	delete(w, id)
}

func (w Working) Len() int {
	return len(w)
}

type Steps []*Step

type Step struct {
	name         string
	finished     bool
	dependencies []*Step
	extraTime    int
	time         int
}

func (s *Step) AddDependency(dep *Step) {
	s.dependencies = append(s.dependencies, dep)
}

func (s *Step) Start() {
	s.time = int((s.name[0])-byte('A')) + 1 + s.extraTime
}

func (s *Step) Tick(sec int) {
	if !s.IsFinished() && s.IsRunning() {
		s.time -= sec
		if s.time < 0 {
			s.time = 0
		}
		if s.time == 0 {
			s.finished = true
		}
	}
}

func (s *Step) IsRunning() bool {
	return s.time > 0
}

func (s *Step) IsFinished() bool {
	return s.finished
}

func (s *Step) CanStart() bool {
	if s.IsFinished() || s.IsRunning() {
		return false
	}
	for _, dep := range s.dependencies {
		if !dep.IsFinished() {
			return false
		}
	}
	return true
}

type StepsMap map[string]*Step

func (sm StepsMap) CanStart() []string {
	var can []string
	for _, step := range sm {
		if step.CanStart() {
			can = append(can, step.name)
		}
	}
	sort.Strings(can)
	return can
}

func (sm StepsMap) Tick(sec int) {
	for _, step := range sm {
		step.Tick(sec)
	}
}

func (sm StepsMap) AllFinished() bool {
	for _, step := range sm {
		if !step.IsFinished() {
			return false
		}
	}
	return true
}

func (sm StepsMap) Get(s string) *Step {
	if step, e := sm[s]; e {
		return step
	}
	step := &Step{name: s, finished: false, dependencies: []*Step{}, extraTime: ExtraTime}
	sm[s] = step
	return sm[s]
}

func getSteps(filename string) StepsMap {
	stepsMap := make(StepsMap)
	lines, _ := file.GetLines(filename)

	for _, line := range lines {
		matches := reStepDep.FindStringSubmatch(line)
		step := stepsMap.Get(matches[2])
		dep := stepsMap.Get(matches[1])
		step.AddDependency(dep)
	}

	return stepsMap
}
