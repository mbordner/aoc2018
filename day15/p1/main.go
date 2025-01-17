package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common"
	"github.com/mbordner/aoc2018/common/file"
	"log"
	"sort"
	"strings"
)

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
	pos Pos
	pt  byte
	hp  int
	ap  int
}

func (p *Player) String() string {
	return fmt.Sprintf("{t:%s, p:%s, hp:%d, ap:%d}", string(p.pt), p.pos, p.hp, p.ap)
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
		log.Printf("player %s attacks %s\n", p, ap)
		ap.hp -= p.ap
		log.Printf("attacked player %s\n", ap)
		return ap
	}
	return nil
}

func (p *Player) Alive() bool {
	return p.hp > 0
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

type Battle struct {
	grid    Grid
	goblins Team
	elves   Team
	round   int
	dead    Players
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
		log.Printf("starting round %d\n", b.round)

		players := b.GetPlayerTurnOrder()
		for _, p := range players {

			if b.Over() {
				log.Printf("round ended early\n")
				return numMoved, numAttacked, numDied // break out of round run, full round can't complete
			}

			if p.Alive() {
				var playerTeam, enemyTeam Team
				if p.pt == Elf {
					playerTeam, enemyTeam = b.elves, b.goblins
				} else {
					playerTeam, enemyTeam = b.goblins, b.elves
				}

				log.Printf("player %s taking turn\n", p)

				var enemy *Player
				if enemy = p.CanAttack(b.grid, enemyTeam); enemy == nil {

					// check if player can move since we can't attack yet
					path := b.GetMovePath(p)
					if len(path) > 0 {
						np := path[0]
						gp := path[len(path)-1]
						log.Printf("player %s wanting to move to %s starting with %s\n", p, gp, np)
						playerTeam.MovePlayer(p, np)
						numMoved++
					} else {
						log.Printf("player %s tried to move but can't\n", p)
					}

					// after move, recheck if we can attack
					enemy = p.CanAttack(b.grid, enemyTeam)
				}

				if enemy != nil { // could attack an enemy

					attacked := p.Attack(b.grid, enemyTeam)
					if attacked != nil {
						numAttacked++
						if !attacked.Alive() {
							enemyTeam.Remove(attacked.pos)
							b.dead = append(b.dead, attacked)
							numDied++
							log.Printf("player %s died\n", attacked)
						}
					}

				}
			}
		}

		log.Printf("round %d complete. moved: %d, attacked: %d, died: %d\n", b.round, numMoved, numAttacked, numDied)
		b.round++
	}

	return numMoved, numAttacked, numDied
}

func (b *Battle) Run() (int, int) {
	fmt.Println(b.GridString())
	for !b.Over() {
		numMoved, numAttacked, numDied := b.RunRound()
		fmt.Printf("After round %d.. moved: %d, attacked: %d, died %d\n", b.round, numMoved, numAttacked, numDied)
		if numMoved > 0 || numDied > 0 {
			fmt.Println(b.GridString())
		}
	}
	return b.Outcome()
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
	// test1.txt ✓  27730
	// test2.txt ✓  36334
	// test3.txt ✓  39514
	// test4.txt ✓  27755
	// test5.txt ✓  28944
	// test6.txt ✓  18740
	b := NewBattle("../data.txt", 200, 3, 200, 3)

	outcome, hp := b.Run()

	log.Printf("outcome: %d after %d rounds with hp %d\n", outcome, b.round, hp)
	log.Println(b.GetPlayerTurnOrder())
	fmt.Printf("outcome: %d after %d rounds with hp %d\n", outcome, b.round, hp)
	fmt.Println(b.GetPlayerTurnOrder())
}
