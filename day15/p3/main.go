package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"github.com/pkg/errors"
	"io"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	Gray    = Color("\033[1;37m%s\033[0m")
	White   = Color("\033[1;97m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var (
	N    = Pos{-1, 0}
	E    = Pos{0, 1}
	S    = Pos{1, 0}
	W    = Pos{0, -1}
	dirs = Positions{N, W, E, S}
)

const (
	Wall   = '#'
	Space  = '.'
	Goblin = 'G'
	Elf    = 'E'
)

type PosContainer map[Pos]bool

func (pc PosContainer) Has(p Pos) bool {
	if b, e := pc[p]; e {
		return b
	}
	return false
}

type PosLinker map[Pos]Pos

type Pos struct {
	Y int
	X int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func (p Pos) Add(o Pos) Pos {
	return Pos{p.Y + o.Y, p.X + o.X}
}

func (p Pos) GridAdjacent(grid Grid) Positions {
	ps := make(Positions, 0, 4)
	for _, d := range dirs {
		o := p.Add(d)
		if grid[o.Y][o.X] != Wall {
			ps = append(ps, o)
		}
	}
	return ps
}

func (p Pos) Adjacent() Positions {
	ps := make(Positions, 0, 4)
	for _, d := range dirs {
		ps = append(ps, p.Add(d))
	}
	return ps
}

type Positions []Pos

type Grid [][]byte

func (g Grid) Clone() Grid {
	og := make(Grid, len(g))
	for y := range g {
		og[y] = make([]byte, len(g[y]))
		copy(og[y], g[y])
	}
	return og
}

type Player struct {
	pos       Pos
	pt        byte
	hp        int
	ap        int
	attacking bool
}

func (p *Player) String() string {
	return fmt.Sprintf("{%s%shp:%d,ap:%d}", string(p.pt), p.pos, p.hp, p.ap)
}

// CanAttack checks if a grid open adjacent position has an enemy team member, returns team member if so
func (p *Player) CanAttack(grid Grid, enemy Team) *Player {
	gridOpenAdj := p.pos.GridAdjacent(grid)
	for _, op := range gridOpenAdj {
		if enemy.Has(op) {
			return enemy[op]
		}
	}
	return nil
}

func (p *Player) Attack(grid Grid, enemy Team) *Player {
	aps := make(Players, 0, 4)
	gridOpenAdj := p.pos.GridAdjacent(grid)
	for _, op := range gridOpenAdj {
		if enemy.Has(op) {
			aps = append(aps, enemy[op])
		}
	}
	sort.Slice(aps, func(i, j int) bool {
		if aps[i].hp < aps[j].hp {
			return true
		} else if aps[i].hp == aps[j].hp {
			if aps[i].pos.Y < aps[j].pos.Y {
				return true
			} else if aps[i].pos.Y == aps[j].pos.Y {
				return aps[i].pos.X < aps[j].pos.X
			}
		}
		return false
	})
	if len(aps) > 0 {
		ap := aps[0]
		ap.hp -= p.ap
		return ap
	}
	return nil
}

func (p *Player) Alive() bool {
	return p.hp > 0
}

func (p *Player) Type() byte {
	return p.pt
}

type Team map[Pos]*Player
type Players []*Player

func (t *Team) Has(p Pos) bool {
	if _, e := (*t)[p]; e {
		return true
	}
	return false
}

func (t *Team) Remove(p Pos) {
	if !t.Empty() {
		delete(*t, p)
	}
}

func (t *Team) HasAdjacent(p Pos) *Player {
	for _, d := range dirs {
		if m, e := (*t)[p.Add(d)]; e {
			return m
		}
	}
	return nil
}

func (t *Team) Type() byte {
	if !t.Empty() {
		for _, m := range *t {
			return m.pt
		}
	}
	return 0
}

func (t *Team) Empty() bool {
	return len(*t) == 0
}

func (t *Team) MovePlayer(p *Player, np Pos) {
	t.Remove(p.pos)
	(*t)[np] = p
	p.pos = np
}

type Teams []Team

func (ts *Teams) AllPlayers() Players {
	count := 0
	for _, t := range *(ts) {
		count += len(t)
	}
	ps := make(Players, 0, count)
	for _, t := range *(ts) {
		for _, p := range t {
			ps = append(ps, p)
		}
	}
	return ps
}

func (ts *Teams) Has(p Pos) bool {
	for _, t := range *(ts) {
		if t.Has(p) {
			return true
		}
	}
	return false
}

func (ts *Teams) EnemyFor(p *Player) *Team {
	for _, t := range *(ts) {
		if t.Type() != p.pt {
			return &t
		}
	}
	return nil
}

func (ts *Teams) TeamFor(p *Player) *Team {
	for _, t := range *(ts) {
		if t.Type() == p.pt {
			return &t
		}
	}
	return nil
}

// Referee func can evaluate a player death and return false if battle should end
type Referee func(p *Player) bool

type Battle struct {
	grid         Grid
	goblins      Team
	elves        Team
	round        int
	dead         Players
	ref          Referee
	disqualified bool
	uiUpdate     chan bool
}

func (b *Battle) UIUpdate() {
	b.uiUpdate <- true
}

func (b *Battle) GridString() string {
	g := b.grid.Clone()
	for p := range b.elves {
		g[p.Y][p.X] = Elf
	}
	for p := range b.goblins {
		g[p.Y][p.X] = Goblin
	}
	ss := make([]string, 0, len(g))
	for y := range g {
		ss = append(ss, string(g[y]))
	}
	return strings.Join(ss, "\n")
}

func (b *Battle) PrintStats(out io.Writer) {
	lines := make([][]string, len(b.grid))
	for y := range b.grid {
		lines[y] = make([]string, 0, 10)
	}
	players := b.GetPlayerTurnOrder()
	for _, player := range players {
		y := player.pos.Y
		if player.attacking {
			lines[y] = append(lines[y], Red(player.String()))
		} else {
			if player.pt == Elf {
				lines[y] = append(lines[y], Teal(player.String()))
			} else {
				lines[y] = append(lines[y], Green(player.String()))
			}
		}
	}
	for y := range lines {
		fmt.Fprintln(out, strings.Join(lines[y], ", "))
	}
}

func (b *Battle) PrintGrid(out io.ReadWriter) {

	grid := make([][]string, len(b.grid))
	for y := range grid {
		grid[y] = make([]string, len(b.grid[y]))
		for x := range grid[y] {
			if b.grid[y][x] == Wall {
				grid[y][x] = fmt.Sprintf("%s", Black(string(Wall)))
			} else {
				grid[y][x] = fmt.Sprintf("%s", Gray(string(Space)))
			}
		}
	}

	for pos, p := range b.elves {
		if p.attacking {
			grid[pos.Y][pos.X] = fmt.Sprintf("%s", Red(string(Elf)))
		} else {
			grid[pos.Y][pos.X] = fmt.Sprintf("%s", Teal(string(Elf)))
		}
	}
	for pos, p := range b.goblins {
		if p.attacking {
			grid[pos.Y][pos.X] = fmt.Sprintf("%s", Red(string(Goblin)))
		} else {
			grid[pos.Y][pos.X] = fmt.Sprintf("%s", Green(string(Goblin)))
		}
	}

	for y := range grid {
		fmt.Fprintln(out, "\u001b[1m"+strings.Join(grid[y], ""))
	}

}

func (b *Battle) GetPlayerTurnOrder() Players {
	ps := make(Players, 0, len(b.goblins)+len(b.elves))
	for _, p := range b.goblins {
		ps = append(ps, p)
	}
	for _, p := range b.elves {
		ps = append(ps, p)
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].pos.Y < ps[j].pos.Y {
			return true
		} else if ps[i].pos.Y == ps[j].pos.Y {
			return ps[i].pos.X < ps[j].pos.X
		}
		return false
	})
	return ps
}

func NewBattle(filename string, eHP, eAP, gHP, gAP int) *Battle {
	b := &Battle{goblins: make(Team), elves: make(Team), round: 0}

	lines, _ := file.GetLines(filename)
	b.grid = make(Grid, len(lines))
	for y, line := range lines {
		b.grid[y] = []byte(line)
		for x, char := range b.grid[y] {
			if char == Elf || char == Goblin {
				p := &Player{pos: Pos{Y: y, X: x}, pt: char}
				if char == Elf {
					p.hp = eHP
					p.ap = eAP
					b.elves[p.pos] = p
				} else {
					p.hp = gHP
					p.ap = gAP
					b.goblins[p.pos] = p
				}
				b.grid[y][x] = Space
			}
		}
	}

	b.dead = make(Players, 0, len(b.elves)+len(b.goblins))

	return b
}

func (b *Battle) Over() bool {
	if b.disqualified {
		return true
	}
	if b.goblins.Empty() || b.elves.Empty() {
		return true
	}
	return false
}

func (b *Battle) GetMovePath(p *Player) Positions {
	var playerTeam, enemyTeam Team
	if p.pt == Elf {
		playerTeam, enemyTeam = b.elves, b.goblins
	} else {
		playerTeam, enemyTeam = b.goblins, b.elves
	}

	queue := make(common.Queue[Pos], 0, 20)
	visited := make(PosContainer)
	prev := make(PosLinker)

	queue.Enqueue(p.pos)
	visited[p.pos] = true

	for !queue.Empty() {
		cur := *(queue.Dequeue())
		if ep := enemyTeam.HasAdjacent(cur); ep != nil {
			path := Positions{cur}
			for pp := prev[cur]; pp != p.pos; pp = prev[pp] {
				path = append(Positions{pp}, path...)
			}
			return path
		} else {
			gridOpenAdjacent := cur.GridAdjacent(b.grid)
			for _, np := range gridOpenAdjacent {
				if !visited.Has(np) {
					visited[np] = true
					if !((&Teams{playerTeam, enemyTeam}).Has(np)) {
						prev[np] = cur
						queue.Enqueue(np)
					}
				}
			}
		}
	}

	return Positions{}
}

func (b *Battle) RunRound() (int, int, int) {
	numMoved, numAttacked, numDied := 0, 0, 0

	if !b.Over() {
		players := b.GetPlayerTurnOrder()
		for _, p := range players {

			if b.Over() {
				return numMoved, numAttacked, numDied // break out of round run, full round can't complete
			}

			if p.Alive() {
				var playerTeam, enemyTeam Team
				if p.pt == Elf {
					playerTeam, enemyTeam = b.elves, b.goblins
				} else {
					playerTeam, enemyTeam = b.goblins, b.elves
				}

				p.attacking = false

				var enemy *Player
				if enemy = p.CanAttack(b.grid, enemyTeam); enemy == nil {

					path := b.GetMovePath(p)
					if len(path) > 0 {
						np := path[0]
						playerTeam.MovePlayer(p, np)
						numMoved++
					}

					// after move, recheck if we can attack
					enemy = p.CanAttack(b.grid, enemyTeam)
				}

				if enemy != nil { // could attack an enemy

					p.attacking = true
					attacked := p.Attack(b.grid, enemyTeam)
					if attacked != nil {
						numAttacked++
						if !attacked.Alive() {
							enemyTeam.Remove(attacked.pos)
							b.dead = append(b.dead, attacked)
							numDied++
							if b.ref != nil {
								if !b.ref(attacked) {
									b.disqualified = true
									// referee says no, battle should end
									return numMoved, numAttacked, numDied
								}
							}
						}
					}

				}
			}
		}

		b.round++
	}

	return numMoved, numAttacked, numDied
}

func (b *Battle) HasAnyAlive(pt byte) bool {
	if pt == Elf {
		return len(b.elves) > 0
	}
	if pt == Goblin {
		return len(b.goblins) > 0
	}
	return false
}

func (b *Battle) Run(wg *sync.WaitGroup) {
	b.UIUpdate()
	for !b.Over() {
		time.Sleep(time.Millisecond * 1100)
		b.RunRound()
		b.UIUpdate()
	}
	wg.Done()
}

func (b *Battle) Outcome() (int, int) {
	hpSum := 0
	for _, p := range b.goblins {
		hpSum += p.hp
	}
	for _, p := range b.elves {
		hpSum += p.hp
	}
	return hpSum * b.round, hpSum
}

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

}

func loop(g *gocui.Gui) {

	elfAP := 20
	dataFile := "../data.txt"

	b := NewBattle(dataFile, 200, elfAP, 200, 3)

	referee := func(p *Player) bool {
		if p.Type() == Elf {
			return false
		}
		return true
	}

	b.ref = referee
	b.uiUpdate = make(chan bool)

	var wg sync.WaitGroup
	wg.Add(1)

	uiUpdateKill := make(chan bool)

	go func() {
		for {
			select {
			case <-b.uiUpdate:

				g.Update(func(g *gocui.Gui) error {
					v, _ := g.View("grid")
					v.Clear()
					b.PrintGrid(v)
					v, _ = g.View("stats")
					v.Clear()
					b.PrintStats(v)
					return nil
				})

			case <-uiUpdateKill:
				return
			}
		}
	}()

	go b.Run(&wg)

	wg.Wait()

	uiUpdateKill <- true
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX > 0 && maxY > 0 {
		if v, err := g.SetView("stats", 34, 0, maxX-4, 33); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
		}
		if v, err := g.SetView("grid", 0, 0, 33, 33); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.BgColor = gocui.ColorWhite
			v.FgColor = gocui.ColorBlack
			go loop(g)
		}
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
