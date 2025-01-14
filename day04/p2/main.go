package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	reEvent = regexp.MustCompile(`\[(.*)\] (falls asleep|wakes up|Guard \#(\d+) begins shift)`)
)

type EventAction int

const (
	UNKNOWN EventAction = iota
	START
	SLEEP
	WAKE
)

func (a EventAction) String() string {
	switch a {
	case START:
		return "start"
	case SLEEP:
		return "sleep"
	case WAKE:
		return "wake"
	}
	return "unknown"
}

type Event struct {
	guard    int
	action   EventAction
	ts       string
	t        time.Time
	startMin int
	dur      int
}

func (e *Event) String() string {
	return fmt.Sprintf("{ts: %s, guard: %d, action: %s}", e.ts, e.guard, e.action)
}

func NewEvent(matches []string) *Event {
	layout := "2006-01-02 15:04" // Go's reference time layout
	t, _ := time.Parse(layout, matches[1])
	e := &Event{t: t, ts: matches[1]}
	e.startMin = t.Minute()
	if strings.Index(matches[2], "sleep") != -1 {
		e.action = SLEEP
	} else if strings.Index(matches[2], "wakes") != -1 {
		e.action = WAKE
	} else if strings.Index(matches[2], "shift") != -1 {
		e.action = START
		e.guard = atoi(matches[3])
	}
	return e
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}

func main() {
	events := getData("../data.txt")

	first := events[0]
	last := events[len(events)-1]

	mins := last.t.Sub(first.t).Minutes()
	fmt.Println(mins)

	guards := make(map[int]map[int]int)
	for _, event := range events {
		g := event.guard
		if _, e := guards[g]; !e {
			guards[g] = make(map[int]int)
		}
		if event.action == SLEEP {
			for m := event.startMin; m < event.startMin+event.dur; m++ {
				if _, minRecorded := guards[g][m]; minRecorded {
					guards[g][m]++
				} else {
					guards[g][m] = 1
				}
			}
		}
	}

	guardID := 0
	guardMinSleptCount := 0
	guardMin := 0

	for gId, minCounts := range guards {
		highestMinCount := 0
		highestMin := 0

		for m, count := range minCounts {
			if count > highestMinCount {
				highestMin = m
				highestMinCount = count
			}
		}

		if highestMinCount > guardMinSleptCount {
			guardMinSleptCount = highestMinCount
			guardID = gId
			guardMin = highestMin
		}
	}

	fmt.Println(guardID * guardMin)

}

func getData(filename string) []*Event {
	lines, _ := file.GetLines(filename)
	events := make([]*Event, len(lines))
	for i, line := range lines {
		matches := reEvent.FindStringSubmatch(line)
		events[i] = NewEvent(matches)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].t.Before(events[j].t)
	})

	var lastGuard int
	var lastAction EventAction
	var lastEventTime time.Time

	for i, event := range events {
		if i > 0 {
			events[i-1].dur = int(event.t.Sub(lastEventTime).Minutes())
		}

		if event.action == START {
			lastGuard = event.guard
		} else if event.guard == 0 {
			event.guard = lastGuard
		}

		if event.action != START {
			if event.action == WAKE {
				if lastAction != SLEEP {
					panic("how can you wake without being asleep?")
				}
			} else if event.action == SLEEP {
				if lastAction == SLEEP {
					panic("going to sleep when already sleeping?")
				}
			}
		}

		lastAction = event.action
		lastEventTime = event.t
	}

	return events
}
