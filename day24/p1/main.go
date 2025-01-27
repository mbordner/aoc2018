package main

import (
	"fmt"
	"github.com/mbordner/aoc2018/common/file"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	reGroup = regexp.MustCompile(`(\d+) units each with (\d+) hit points (?:\((.*)\)\s)*with an attack that does (\d+) (bludgeoning|fire|radiation|slashing) damage at initiative (\d+)`)
)

type DamageMap map[string]bool

func (dm DamageMap) Has(d string) bool {
	if b, e := dm[d]; e {
		return b
	}
	return false
}

type GroupMap map[*Group]bool

type GroupTargets map[*Group]*Group

func (gt GroupTargets) Has(g *Group) bool {
	if _, e := gt[g]; e {
		return true
	}
	return false
}

func (gm GroupMap) Clone() GroupMap {
	ogm := make(GroupMap)
	for g, alive := range gm {
		if alive {
			ogm[g] = true
		}
	}
	return ogm
}

type Group struct {
	immunity   DamageMap
	weakness   DamageMap
	units      int
	hp         int
	ap         int
	dt         string
	initiative int
}

func NewGroup(units, hp, ap, initiative int, dt string) *Group {
	return &Group{immunity: make(DamageMap), weakness: make(DamageMap), units: units, hp: hp, ap: ap, initiative: initiative, dt: dt}
}

func (g *Group) SetWeakness(dt string) {
	g.weakness[dt] = true
}

func (g *Group) SetImmunity(dt string) {
	g.immunity[dt] = true
}

func (g *Group) EffectivePower() int {
	return g.units * g.ap
}

func (g *Group) Initiative() int {
	return g.initiative
}

func (g *Group) DamageType() string {
	return g.dt
}

func (g *Group) WeakTo(d string) bool {
	return g.weakness.Has(d)
}

func (g *Group) ImmuneTo(d string) bool {
	return g.immunity.Has(d)
}

func (g *Group) DamageTo(og *Group) int {
	if og.ImmuneTo(g.DamageType()) {
		return 0
	}
	ep := g.EffectivePower()
	if og.WeakTo(g.DamageType()) {
		return 2 * ep
	}
	return ep
}

func (g *Group) DamageFrom(og *Group) {
	damage := og.DamageTo(g)
	unitsLost := damage / g.hp
	g.units -= unitsLost
	if g.units < 0 {
		g.units = 0
	}
}

func (g *Group) UnitsAlive() int {
	return g.units
}

type Army struct {
	name       string
	groups     GroupMap
	groupCount int
}

func NewArmy(name string) *Army {
	return &Army{name: name, groups: make(GroupMap)}
}

func (a *Army) AddGroup(group *Group) {
	a.groups[group] = true
	a.groupCount++
}

func (a *Army) KillGroup(group *Group) {
	a.groups[group] = false
	a.groupCount--
}

func sortGroupsByEffectivePower(groups []*Group) {
	sort.Slice(groups, func(i, j int) bool {
		iep, jep := groups[i].EffectivePower(), groups[j].EffectivePower()
		if iep > jep {
			return true
		} else if iep == jep {
			return groups[i].Initiative() > groups[j].Initiative()
		}
		return false
	})
}

func sortGroupsByInitiative(groups []*Group) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Initiative() > groups[j].Initiative()
	})
}

func (a *Army) GetGroupsByEffectivePower() []*Group {
	groups := make([]*Group, 0, a.groupCount)
	for group, alive := range a.groups {
		if alive {
			groups = append(groups, group)
		}
	}
	sortGroupsByEffectivePower(groups)
	return groups
}

func (a *Army) SelectTargets(oa *Army) GroupTargets {
	possibleTargets := oa.GetGroupsByEffectivePower()

	selectedTargets := make(GroupTargets)

	for _, group := range a.GetGroupsByEffectivePower() {
		maxDamage := 0
		var maxDamageTarget *Group
		for _, target := range possibleTargets {
			if !selectedTargets.Has(target) {
				damage := group.DamageTo(target)
				if damage > 0 {
					if damage > maxDamage {
						maxDamageTarget = target
						maxDamage = damage
					} else if damage == maxDamage {
						step, tep := maxDamageTarget.EffectivePower(), target.EffectivePower()
						if tep > step {
							maxDamageTarget = target
						} else if tep == step {
							if target.Initiative() > maxDamageTarget.Initiative() {
								maxDamageTarget = target
							}
						}
					}
				}
			}
		}
		if maxDamageTarget != nil {
			selectedTargets[maxDamageTarget] = group
		}
	}
	return selectedTargets
}

func (a *Army) GetAliveGroups() GroupMap {
	return a.groups.Clone()
}

func (a *Army) GetGroupsAliveCount() int {
	return a.groupCount
}

type Armies []*Army

func (as Armies) AllGroupsHaveUnits() bool {
	for _, a := range as {
		if a.GetGroupsAliveCount() == 0 {
			return false
		}
	}
	return true
}

func (as Armies) GetGroupsAliveCount() int {
	count := 0
	for _, a := range as {
		count += a.GetGroupsAliveCount()
	}
	return count
}

func (as Armies) GetUnitsAliveCount() int {
	count := 0
	for _, a := range as {
		for g := range a.GetAliveGroups() {
			count += g.UnitsAlive()
		}
	}
	return count
}

func (as Armies) Fight() {
	selected := make(GroupTargets)
	attackers := make([]*Group, 0, as.GetGroupsAliveCount())
	gId := make(map[*Group]int)
	for i := range as {
		oi := i + 1
		if oi == len(as) {
			oi = 0
		}
		selectedTargets := as[i].SelectTargets(as[oi])
		for target, attacker := range selectedTargets {
			attackers = append(attackers, attacker)
			selected[attacker] = target
			gId[attacker] = i
			gId[target] = oi
		}
	}
	sortGroupsByInitiative(attackers)

	for _, attacker := range attackers {
		if attacker.UnitsAlive() > 0 {
			target := selected[attacker]
			target.DamageFrom(attacker)
			if target.UnitsAlive() == 0 {
				as[gId[target]].KillGroup(target)
			}
		}
	}
}

func main() {
	armies := getData("../data.txt")

	for armies.AllGroupsHaveUnits() {
		armies.Fight()
	}

	fmt.Println(armies.GetUnitsAliveCount())
}

func getData(filename string) Armies {
	content, _ := file.GetContent(filename)
	teams := strings.Split(strings.TrimSpace(string(content)), "\n\n")
	armies := make(Armies, len(teams))
	for i, team := range teams {
		lines := strings.Split(team, "\n")
		name := lines[0][0 : len(lines[0])-1]
		armies[i] = NewArmy(name)
		for _, line := range lines[1:] {
			matches := reGroup.FindStringSubmatch(line)
			g := NewGroup(atoi(matches[1]), atoi(matches[2]), atoi(matches[4]), atoi(matches[6]), matches[5])

			tokens := strings.Split(strings.TrimSpace(matches[3]), "; ")
			for _, token := range tokens {
				immuneTo := "immune to "
				weakTo := "weak to "
				if strings.HasPrefix(token, immuneTo) {
					immunities := strings.Split(token[len(immuneTo):], ", ")
					for _, immunity := range immunities {
						g.SetImmunity(immunity)
					}
				} else if strings.HasPrefix(token, weakTo) {
					weaknesses := strings.Split(token[len(weakTo):], ", ")
					for _, weakness := range weaknesses {
						g.SetWeakness(weakness)
					}
				}
			}

			armies[i].AddGroup(g)
		}
	}
	return armies
}

func atoi(s string) int {
	val, _ := strconv.ParseInt(s, 10, 64)
	return int(val)
}
