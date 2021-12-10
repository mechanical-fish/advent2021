package main

import (
	"bufio"
	"fmt"
	"os"
)

func readFloor() ([][]uint8, error) {
	floor := make([][]uint8, 0, 5)
	row := make([]uint8, 0, 10)
	var rowLength int
	inBuf := bufio.NewReader(os.Stdin)
	for lin, err := inBuf.ReadBytes('\n'); err == nil; lin, err = inBuf.ReadBytes('\n') {
		for i, b := range lin {
			if b < '0' || b > '9' {
				break
			}
			if rowLength == 0 {
				row = append(row, uint8(b-'0'))
			} else {
				row[i] = uint8(b - '0')
			}
		}
		rowLength = len(row)
		floor = append(floor, row)
		row = make([]uint8, rowLength)
	}
	return floor, nil
}

type Point struct {
	Row    int
	Col    int
	Height uint8
}

func (p Point) RiskLevel() int {
	return int(p.Height + 1)
}

func main() {
	floor, err := readFloor()
	if err != nil {
		panic(err)
	}

	// process two rows at a time
	rowA := 0
	rowB := 0
	var Acandidates map[int]bool
	Bcandidates := make(map[int]bool)
	lows := make([]Point, 0, 5)
	for rowB < len(floor) {
		Acandidates = Bcandidates
		Bcandidates = make(map[int]bool)
		var candidate bool = false
		for i, d := range floor[rowB] {
			if candidate && floor[rowB][i-1] < d {
				// point to the left is a candidate, check the point
				// below it when we get to the next row
				// fmt.Printf("Consider row %d, col %d\n", rowB, i)
				Bcandidates[i-1] = true
			}
			// is current point a candidate? check points above and
			// to the left
			candidate = (i == 0 || floor[rowB][i-1] > d) &&
				(rowA == rowB || floor[rowA][i] > d)
			// if this is a candidate and the rightmost point,
			// it goes on the list
			if candidate && i == len(floor[rowB])-1 {
				// fmt.Printf("Consider row %d, col %d\n", rowB, i)
				Bcandidates[i] = true
			}
			if Acandidates[i] && floor[rowA][i] < d {
				lows = append(lows, Point{Row: rowA, Col: i, Height: floor[rowA][i]})
			}
		}
		rowA = rowB
		rowB++
	}
	// process the last row
	for i := range Bcandidates {
		lows = append(lows, Point{Row: rowA, Col: i, Height: floor[rowA][i]})
	}

	// print the lows
	lowArray := make([][]byte, len(floor))
	for i := range floor {
		lowArray[i] = make([]byte, len(floor[0]))
		for j := range lowArray[i] {
			lowArray[i][j] = '.'
		}
	}
	risk := 0
	for _, p := range lows {
		lowArray[p.Row][p.Col] = p.Height + '0'
		risk += p.RiskLevel()
	}

	for _, r := range lowArray {
		for _, d := range r {
			fmt.Printf("%c", d)
		}
		fmt.Println()
	}

	fmt.Printf("Total risk is %d\n", risk)
}
