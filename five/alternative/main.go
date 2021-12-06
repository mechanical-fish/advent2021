package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Line struct {
	Xa, Ya, Xb, Yb int
}

func NewLine(xa, ya, xb, yb int) *Line {
	return &Line{Xa: xa, Ya: ya, Xb: xb, Yb: yb}
}

func ParseLine(s string) *Line {
	f := strings.Fields(s)
	if f[1] != "->" {
		panic("A line of input did not contain an arrow")
	}
	ptStrings := append(strings.Split(f[0], `,`), strings.Split(f[2], `,`)...)
	pts := make([]int, len(ptStrings))
	var err error
	for i, ps := range ptStrings {
		if pts[i], err = strconv.Atoi(ps); err != nil {
			panic(err)
		}
	}
	return NewLine(pts[0], pts[1], pts[2], pts[3])
}

type Point struct {
	X int
	Y int
}

func (n Line) PointsOn(useDiagonals bool) []Point {
	result := make([]Point, 0, 2)
	var dX, dY int = 0, 0
	if n.Xa < n.Xb {
		dX = 1
	} else if n.Xa > n.Xb {
		dX = -1
	}
	if n.Ya < n.Yb {
		dY = 1
	} else if n.Ya > n.Yb {
		dY = -1
	}
	if !useDiagonals && dX != 0 && dY != 0 {
		return result
	}
	var x int = n.Xa
	var y int = n.Ya
	for x != n.Xb || y != n.Yb {
		result = append(result, Point{x, y})
		x += dX
		y += dY
	}
	return append(result, Point{n.Xb, n.Yb})
}

type Plane struct {
	Grid [][]uint8
}

func NewPlane(xSize, ySize int) *Plane {
	p := &Plane{
		Grid: make([][]uint8, xSize),
	}
	for i := range p.Grid {
		p.Grid[i] = make([]uint8, ySize)
	}
	return p
}

func (p *Plane) AddLine(n *Line, useDiagonals bool) {
	for _, pt := range n.PointsOn(useDiagonals) {
		p.Grid[pt.X][pt.Y] = p.Grid[pt.X][pt.Y] + 1
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	plane := NewPlane(1000, 1000)
	for scanner.Scan() {
		plane.AddLine(ParseLine(scanner.Text()), true)
	}
	overlaps := 0
	for _, row := range plane.Grid {
		for _, o := range row {
			if o >= 2 {
				overlaps++
			}
		}
	}
	fmt.Printf("Found %d overlaps\n", overlaps)
}
