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

func main() {
	steps := getSteps("../data.txt")

	executed := make([]string, 0, len(steps))

	for can := steps.CanExecute(); len(can) > 0; can = steps.CanExecute() {
		steps[can[0]].Execute()
		executed = append(executed, can[0])
	}

	fmt.Println(strings.Join(executed, ""))
}

type Steps []*Step

type Step struct {
	name         string
	finished     bool
	dependencies []*Step
}

func (s *Step) AddDependency(dep *Step) {
	s.dependencies = append(s.dependencies, dep)
}

func (s *Step) Execute() {
	s.finished = true
}

func (s *Step) IsFinished() bool {
	return s.finished
}

func (s *Step) CanExecute() bool {
	if s.IsFinished() {
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

func (sm StepsMap) CanExecute() []string {
	var can []string
	for _, step := range sm {
		if step.CanExecute() {
			can = append(can, step.name)
		}
	}
	sort.Strings(can)
	return can
}

func (sm StepsMap) Get(s string) *Step {
	if step, e := sm[s]; e {
		return step
	}
	step := &Step{name: s, finished: false, dependencies: []*Step{}}
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
