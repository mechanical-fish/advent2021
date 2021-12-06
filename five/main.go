package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Line struct {
	Xa int
	Ya int
	Xb int
	Yb int
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

func (n Line) IsHorizontal() bool {
	return n.Ya == n.Yb
}

func (n Line) IsVertical() bool {
	return n.Xa == n.Xb
}

type Point struct {
	X int
	Y int
}

func (n Line) PointsOn() []Point {
	result := make([]Point, 0, 2)
	if n.IsHorizontal() {
		xS, xL := n.Xa, n.Xb
		if xS > xL {
			xS, xL = xL, xS
		}
		for x := xS; x <= xL; x++ {
			result = append(result, Point{x, n.Ya})
		}
	} else if n.IsVertical() {
		yS, yL := n.Ya, n.Yb
		if yS > yL {
			yS, yL = yL, yS
		}
		for y := yS; y <= yL; y++ {
			result = append(result, Point{n.Xa, y})
		}
	} else {
		// rearrange the endpoints so the first point is on the left
		xL, yL, xR, yR := n.Xa, n.Ya, n.Xb, n.Yb
		if xL > xR {
			xL, xR, yL, yR = xR, xL, yR, yL
		}
		var slope int = 1
		if yR < yL {
			slope = -1
		}
		y := yL
		for x := xL; x <= xR; x++ {
			result = append(result, Point{x, y})
			y += slope
		}
	}
	return result
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

func (p *Plane) AddLine(n *Line) {
	for _, pt := range n.PointsOn() {
		p.Grid[pt.X][pt.Y] = p.Grid[pt.X][pt.Y] + 1
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	plane := NewPlane(1000, 1000)
	for scanner.Scan() {
		plane.AddLine(ParseLine(scanner.Text()))
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
