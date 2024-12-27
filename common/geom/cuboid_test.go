package geom

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Volume(t *testing.T) {

	cases := []struct {
		input  string
		volume uint64
	}{
		{input: `0,0,0,1,1,1`, volume: uint64(1)},
		{input: `0,0,0,0,0,0`, volume: uint64(0)},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			v := NewCuboid(tc.input).Volume()
			assert.Equal(t, tc.volume, v)
		})
	}
}

func Test_CubuoidsVolume(t *testing.T) {

	cases := []struct {
		cube     string
		splitpt  string
		expected uint64
	}{
		{cube: `0,0,0,5,5,5`, splitpt: `2,2,2`, expected: 64},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.cube)
			splits := c.SplitAt(NewPoint(tc.splitpt))
			assert.Equal(t, tc.expected, splits.Volume())
		})
	}

}

func Test_PointsCount(t *testing.T) {

	cases := []struct {
		input  string
		volume uint64
	}{
		{input: `0,0,0,1,1,1`, volume: uint64(8)},
		{input: `0,0,0,0,0,0`, volume: uint64(1)},
		{input: `0,0,0,2,2,2`, volume: uint64(27)},
		{input: `0,0,0,2,2,1`, volume: uint64(18)},
		{input: `0,0,0,2,2,0`, volume: uint64(9)},
		{input: `0,0,0,3,3,0`, volume: uint64(16)},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			v := NewCuboid(tc.input).PointsCount()
			assert.Equal(t, tc.volume, v)
		})
	}
}

func Test_PointsContainsViaCorners(t *testing.T) {
	cases := []struct {
		cuboid   string
		points   []string
		expected bool
	}{
		{cuboid: `0,0,0,1,1,1`, points: []string{`0,0,0`, `0,1,0`, `1,1,0`, `1,0,0`, `0,0,1`, `0,1,1`, `1,1,1`, `1,0,1`}, expected: true},
		{cuboid: `0,0,0,1,1,1`, points: []string{`2,2,2`}, expected: false},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.cuboid)
			corners := c.Corners()
			for _, s := range tc.points {
				p := NewPoint(s)
				val := corners.Contains(p)

				assert.Equal(t, tc.expected, val)
			}

		})
	}

}

func Test_TransformPoint(t *testing.T) {
	cases := []struct {
		point    string
		vector   string
		expected string
	}{
		{point: `0,0,0`, vector: `1,-1,1`, expected: `1,-1,1`},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewPoint(tc.point).Transform(NewVector(tc.vector)).String())
		})
	}
}

func Test_CuboidContains(t *testing.T) {
	cases := []struct {
		cuboid   string
		point    string
		expected bool
	}{
		{`0,0,0,2,2,2`, `1,1,1`, true},
		{`0,0,0,2,2,2`, `-1,-1,-1`, false},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewCuboid(tc.cuboid).Contains(NewPoint(tc.point)))
		})
	}
}

func Test_CuboidsContains(t *testing.T) {
	cases := []struct {
		cuboids  []string
		check    string
		expected bool
	}{
		{
			cuboids:  []string{`0,0,0,1,1,1`, `2,2,2,3,3,3`, `4,4,4,5,5,5`},
			check:    `2,2,2,3,3,3`,
			expected: true,
		},
		{
			cuboids:  []string{`0,0,0,1,1,1`, `2,2,2,3,3,3`, `4,4,4,5,5,5`},
			check:    `6,6,6,7,7,7`,
			expected: false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			cuboids := make(Cuboids, 0, len(tc.cuboids))
			for _, c := range tc.cuboids {
				cuboids = append(cuboids, NewCuboid(c))
			}
			assert.Equal(t, tc.expected, cuboids.Contains(NewCuboid(tc.check)))
		})
	}
}

func Test_IntersectingCorners(t *testing.T) {
	cases := []struct {
		c1      string
		c2      string
		corners []string
	}{
		{`0,0,0,2,2,2`, `-1,-1,-1,1,1,1`, []string{`1,1,1`}},
		{`0,0,0,2,2,2`, `-2,-2,-2,-1,-1,-1`, []string{}},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			corners := NewCuboid(tc.c1).IntersectingCorners(NewCuboid(tc.c2))

			assert.Equal(t, len(tc.corners), len(corners))
		})
	}
}

func Test_IsCorner(t *testing.T) {
	cases := []struct {
		c1       string
		p        string
		expected bool
	}{
		{`0,0,0,1,1,1`, `1,1,1`, true},
		{`0,0,0,1,1,1`, `-1,-1,-1`, false},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewCuboid(tc.c1).IsCorner(NewPoint(tc.p)))
		})
	}
}

func Test_CuboidTransform(t *testing.T) {
	cases := []struct {
		c1       string
		v        string
		expected string
	}{
		{`0,0,0,1,1,1`, `1,1,1`, `1,1,1,2,2,2`},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewCuboid(tc.c1).Transform(NewVector(tc.v)).String())
		})
	}
}

func Test_SplitAt(t *testing.T) {
	type pointCheck struct {
		p     string
		count int
	}

	cases := []struct {
		c      string
		p      string
		l      int
		checks []pointCheck // loops on points across the split cuboids and assets the points exit in *count* cuboids
	}{
		{
			c: `-2,-2,-2,2,2,2`,
			p: `0,0,0`,
			l: 8,
		},
		{
			c: `0,0,0,2,2,2`,
			p: `1,1,1`,
			l: 8,
		},
		{
			c: `0,0,0,3,3,3`,
			p: `1,1,1`,
			l: 8,
		},
		{
			c: `0,0,0,4,4,4`,
			p: `1,1,1`,
			l: 8,
		},
		{
			c: `0,0,0,4,4,4`,
			p: `2,2,2`,
			l: 8,
		},
		{
			c: `0,0,0,5,5,5`,
			p: `2,2,2`,
			l: 8,
		},
		{
			c: `0,0,0,10,10,10`,
			p: `2,2,2`,
			l: 8,
			checks: []pointCheck{
				{
					p:     `2,2,2`,
					count: 1,
				},
			},
		},
		{
			c: `0,0,0,10,10,10`,
			p: `-1,-1,-1`,
			l: 0,
		},
		{
			c: `0,0,0,10,10,10`,
			p: `0,0,0`,
			l: 8,
			checks: []pointCheck{
				{
					p:     `0,0,0`,
					count: 1,
				},
				{
					p:     `1,1,1`,
					count: 1,
				},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.c)
			p := NewPoint(tc.p)

			splits := c.SplitAt(p)

			assert.Equal(t, tc.l, len(splits))

			if len(splits) > 0 {
				cuboidPointsCount := c.PointsCount()
				splitsPointsCount := splits.PointsCount()
				assert.Equal(t, cuboidPointsCount, splitsPointsCount)
			}

			if tc.checks != nil && len(tc.checks) > 0 {
				for _, check := range tc.checks {
					tp := NewPoint(check.p)

					count := 0
					for _, c := range splits {
						if c.Contains(tp) {
							count++
						}
					}

					assert.Equal(t, check.count, count)
				}
			}
		})
	}
}

func Test_Points(t *testing.T) {

	cases := []struct {
		c  string
		ps string
	}{
		{
			c:  `0,0,0,1,1,1`,
			ps: `0,0,0,0,0,1,0,1,0,0,1,1,1,0,0,1,0,1,1,1,0,1,1,1`,
		},
		{
			c:  `0,0,0,0,0,0`,
			ps: `0,0,0`,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.c)
			ps := c.Points()
			assert.Equal(t, tc.ps, fmt.Sprintf("%s", ps))
		})
	}
}

func Test_Snap(t *testing.T) {

	cases := []struct {
		cuboid string
	}{
		{`0,0,0,1,1,1`},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			c := NewCuboid(tc.cuboid)
			assert.True(t, c.IsOnEdge(c.Min))
			assert.True(t, c.IsOnEdge(c.Max))

			assert.True(t, c.Min.Snap(c.Max, X).Snap(c.Max, Y).Snap(c.Max, Z) == c.Max)
			assert.True(t, c.Max.Snap(c.Min, X).Snap(c.Min, Y).Snap(c.Min, Z) == c.Min)
		})
	}
}

func Test_Encloses(t *testing.T) {

	cases := []struct {
		c string
		o string
	}{
		{
			c: `0,0,0,2,2,1`,
			o: `0,0,0,1,1,1`,
		},
		{
			c: `0,0,0,1,1,1`,
			o: `0,0,0,2,2,1`,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.c)
			o := NewCuboid(tc.o)
			assert.True(t, (c.Encloses(o) || o.Encloses(c)) && c.Overlaps(o))
		})
	}
}

func Test_SplitCombineBack(t *testing.T) {
	cases := []struct {
		cuboid string
		point  string
	}{
		{
			cuboid: `0,0,0,3,3,3`,
			point:  `1,1,1`,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.cuboid)
			p := NewPoint(tc.point)
			combined := c.SplitAt(p).Combine()
			assert.Equal(t, 1, len(combined))
			assert.Equal(t, c, combined[0])
		})
	}
}

func Test_CombineToOne(t *testing.T) {

	cases := []struct {
		cuboids []string
	}{
		{
			cuboids: []string{
				`0,0,0,1,1,1`,
				`2,0,0,3,1,1`,
				`4,0,0,5,1,1`,

				`0,2,0,1,3,1`,
				`2,2,0,3,3,1`,
				`4,2,0,5,3,1`,

				`0,4,0,1,5,1`,
				`2,4,0,3,5,1`,
				`4,4,0,5,5,1`,

				`0,0,2,1,1,3`,
				`2,0,2,3,1,3`,
				`4,0,2,5,1,3`,

				`0,2,2,1,3,3`,
				`2,2,2,3,3,3`,
				`4,2,2,5,3,3`,

				`0,4,2,1,5,3`,
				`2,4,2,3,5,3`,
				`4,4,2,5,5,3`,

				`0,0,4,1,1,5`,
				`2,0,4,3,1,5`,
				`4,0,4,5,1,5`,

				`0,2,4,1,3,5`,
				`2,2,4,3,3,5`,
				`4,2,4,5,3,5`,

				`0,4,4,1,5,5`,
				`2,4,4,3,5,5`,
				`4,4,4,5,5,5`,
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			cuboids := make(Cuboids, 0, 8)
			for _, c := range tc.cuboids {
				cuboids = append(cuboids, NewCuboid(c))
			}

			cuboids = cuboids.Combine()

			assert.Equal(t, 1, len(cuboids))

		})
	}

}

func TestCuboid_IntersectCuboidsPoints(t *testing.T) {

	cases := []struct {
		cuboids []string
		points  string
	}{
		{
			cuboids: []string{`0,0,0,1,1,1`, `2,2,2,2,2,2`, `3,3,3,3,3,3`},
			points: Points{
				NewPoint(`0,0,0`),
				NewPoint(`0,1,0`),
				NewPoint(`1,0,0`),
				NewPoint(`1,1,0`),

				NewPoint(`0,0,1`),
				NewPoint(`0,1,1`),
				NewPoint(`1,0,1`),
				NewPoint(`1,1,1`),

				NewPoint(`2,2,2`),

				NewPoint(`3,3,3`),
			}.DeDup().Sort().String(),
		},
		{
			cuboids: []string{`0,0,0,1,2,0`, `2,1,0,3,2,0`},
			points:  `0,0,0,0,1,0,0,2,0,1,0,0,1,1,0,1,2,0,2,1,0,2,2,0,3,1,0,3,2,0`,
		},
		{
			cuboids: []string{`1,1,0,5,3,0`, `2,0,0,6,2,0`, `0,0,0,3,2,0`},
			points:  `0,0,0,0,1,0,0,2,0,1,0,0,1,1,0,1,2,0,1,3,0,2,0,0,2,1,0,2,2,0,2,3,0,3,0,0,3,1,0,3,2,0,3,3,0,4,0,0,4,1,0,4,2,0,4,3,0,5,0,0,5,1,0,5,2,0,5,3,0,6,0,0,6,1,0,6,2,0`,
		},
		{
			cuboids: []string{`3,1,1,3,3,3`, `1,1,3,2,3,3`}, //  `0,0,0,2,2,2`
			points:  `1,1,3,1,2,3,1,3,3,2,1,3,2,2,3,2,3,3,3,1,1,3,1,2,3,1,3,3,2,1,3,2,2,3,2,3,3,3,1,3,3,2,3,3,3`,
		},
		{
			cuboids: []string{`13,11,11,13,13,13`, `11,11,13,12,12,13`},
			points:  NewPoints(`11,11,13,11,12,13,12,11,13,12,12,13,13,11,11,13,11,12,13,11,13,13,12,11,13,12,12,13,12,13,13,13,11,13,13,12,13,13,13`).DeDup().Sort().String(),
		},
		{
			cuboids: []string{`10,10,10,12,12,12`, `11,11,11,13,13,13`},
			points:  NewPoints(`10,10,10,10,10,11,10,10,12,10,11,10,10,11,11,10,11,12,10,12,10,10,12,11,10,12,12,11,10,10,11,10,11,11,10,12,11,11,10,11,11,11,11,11,12,11,11,13,11,12,10,11,12,11,11,12,12,11,12,13,11,13,11,11,13,12,11,13,13,12,10,10,12,10,11,12,10,12,12,11,10,12,11,11,12,11,12,12,11,13,12,12,10,12,12,11,12,12,12,12,12,13,12,13,11,12,13,12,12,13,13,13,11,11,13,11,12,13,11,13,13,12,11,13,12,12,13,12,13,13,13,11,13,13,12,13,13,13`).DeDup().Sort().String(),
		},
		{
			cuboids: []string{`13,11,11,13,13,13`, `11,11,13,12,12,13`, `11,13,11,12,13,13`},
			points:  NewPoints(`11,11,13,11,12,13,11,13,11,11,13,12,11,13,13,12,11,13,12,12,13,12,13,11,12,13,12,12,13,13,13,11,11,13,11,12,13,11,13,13,12,11,13,12,12,13,12,13,13,13,11,13,13,12,13,13,13`).DeDup().Sort().String(),
		},
		{
			cuboids: []string{`10,10,10,12,12,12`, `13,11,11,13,13,13`, `11,11,13,12,13,13`},
			points:  NewPoints(`10,10,10,10,10,11,10,10,12,10,11,10,10,11,11,10,11,12,10,12,10,10,12,11,10,12,12,11,10,10,11,10,11,11,10,12,11,11,10,11,11,11,11,11,12,11,11,13,11,12,10,11,12,11,11,12,12,11,12,13,11,13,13,12,10,10,12,10,11,12,10,12,12,11,10,12,11,11,12,11,12,12,11,13,12,12,10,12,12,11,12,12,12,12,12,13,12,13,13,13,11,11,13,11,12,13,11,13,13,12,11,13,12,12,13,12,13,13,13,11,13,13,12,13,13,13`).DeDup().Sort().String(),
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			var cuboids Cuboids
			var points Points

			for n := 0; n < len(tc.cuboids); n++ {
				cuboidsPointsCount := cuboids.PointsCount()
				newCuboid := NewCuboid(tc.cuboids[n])
				newCuboidPointsCount := newCuboid.PointsCount()
				cuboids = cuboids.Merge(newCuboid)
				newCuboidsPointCount := cuboids.PointsCount()
				assert.True(t, cuboidsPointsCount+newCuboidPointsCount >= newCuboidsPointCount)
			}

			points = cuboids.Points()

			checkPoints := NewPoints(tc.points)

			assert.Equal(t, len(checkPoints), len(points))

			assert.Equal(t, tc.points, points.String())

			noneOverlap := true
		outer:
			for a := 0; a < len(cuboids); a++ {
				for b := 0; b < len(cuboids); b++ {
					if a != b {
						c := cuboids[a]
						o := cuboids[b]
						if c.Overlaps(o) {
							noneOverlap = false
							break outer
						}
					}
				}
			}

			assert.True(t, noneOverlap)

		})
	}

}

func Test_IntersectCuboids(t *testing.T) {

	type results struct {
		c1   []string
		both []string
		c2   []string
	}

	cases := []struct {
		c1     string
		c2     string
		checks results
	}{
		{
			c1: `5,5,5,10,10,10`,
			c2: `0,0,0,4,4,4`,
			checks: results{
				c1:   []string{`5,5,5,10,10,10`},
				both: []string{},
				c2:   []string{`0,0,0,4,4,4`},
			},
		},
		{
			c1: `0,0,0,5,5,5`,
			c2: `-5,-5,-5,2,2,2`,
			checks: results{
				c1:   []string{`3,0,0,5,5,5`, `0,0,3,2,2,5`, `0,3,0,2,5,5`},
				both: []string{`0,0,0,2,2,2`},
				c2:   []string{`-5,-5,-5,-1,2,2`, `0,0,-5,2,2,-1`, `0,-5,-5,2,-1,2`},
			},
		},
		{
			c1: `0,0,0,10,10,10`,
			c2: `0,0,0,10,10,5`,
			checks: results{
				c1:   []string{`0,0,6,10,10,10`},
				both: []string{`0,0,0,10,10,5`},
				c2:   []string{},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c1 := NewCuboid(tc.c1)
			c2 := NewCuboid(tc.c2)
			fromC1, both, fromC2 := c1.Intersect(c2)

			c1pc := c1.PointsCount()
			c2pc := c2.PointsCount()
			interpc := both.PointsCount()

			fromc1pc := fromC1.PointsCount()
			fromc2pc := fromC2.PointsCount()

			val1 := c1pc + c2pc - interpc
			val2 := fromc1pc + interpc + fromc2pc

			assert.Equal(t, val1, val2)

			checkFromC1 := make(Cuboids, 0, len(tc.checks.c1))
			for _, c := range tc.checks.c1 {
				checkFromC1 = append(checkFromC1, NewCuboid(c))
			}

			checkBoth := make(Cuboids, 0, len(tc.checks.both))
			for _, c := range tc.checks.both {
				checkBoth = append(checkBoth, NewCuboid(c))
			}

			checkFromC2 := make(Cuboids, 0, len(tc.checks.c2))
			for _, c := range tc.checks.c2 {
				checkFromC2 = append(checkFromC2, NewCuboid(c))
			}

			assert.Equal(t, len(checkFromC1), len(fromC1))
			assert.Equal(t, len(checkBoth), len(both))
			assert.Equal(t, len(checkFromC2), len(fromC2))

			for _, check := range [][]Cuboids{{checkFromC1, fromC1}, {checkBoth, both}, {checkFromC2, fromC2}} {
				for _, c := range check[0] {
					if check[1].Contains(c) == false {
						fmt.Println("wtf")
					}
					assert.Equal(t, true, check[1].Contains(c))
				}
			}
		})
	}
}
