package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	var inCrabs []int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		numStrings := strings.Split(scanner.Text(), `,`)
		if inCrabs == nil {
			inCrabs = make([]int, len(numStrings))
		}
		for i, n := range strings.Split(scanner.Text(), `,`) {
			pos, err := strconv.Atoi(n)
			if err != nil {
				panic(err)
			}
			inCrabs[i] = pos
		}
	}

	// for _, n := range inCrabs {
	// 	fmt.Printf("%d ", n)
	// }
	// fmt.Println()

	start := time.Now()

	sort.Ints(inCrabs)
	end := time.Since(start)
	fmt.Printf("Sort time is %v\n", end)

	// Consider the leftmost and rightmost crabs. To bring those two
	// crabs together at any position between them will cost the same
	// amount of fuel, because a position 1 unit closer to the right
	// crab will be 1 unit further from the left. Any position not
	// bracketed by the two crabs will cost more fuel than that.
	//
	// Repeat this logic to see that the min-fuel position must also
	// lie between the second-leftmost and second-rightmost pair of
	// crabs.
	//
	// Repeat this for each pair of crabs until we have only one pair
	// left. Either of their positions, or any position in between, is
	// the optimum. Or, if there is one crab left, its position is the
	// optimum. In short, the optimum place for crab meetups is the
	// integer closest to the median position.
	start = time.Now()
	answer := inCrabs[len(inCrabs)/2]
	fuelCost := 0
	for _, pos := range inCrabs {
		if answer > pos {
			fuelCost += answer - pos
		} else {
			fuelCost += pos - answer
		}
	}
	end = time.Since(start)
	fmt.Printf("The optimum location is %d\n", answer)
	fmt.Printf("The fuel cost is %d\n", fuelCost)
	fmt.Printf("Executed in %v\n", end)

	// If crab engines burn fuel at a nonlinear rate, the Gauss
	// formula for the sum of ints from 1 to N applies.
	//
	// Some calculus on a piece of paper suggests that the minimum
	// fuel cost is obtained by putting the rendezvous point close
	// to the mean location of all the crabs. if that location is x,
	// then the fuel cost to get a crab at a to point x is (x-a)(x-a+1)/2

	start = time.Now()
	fuelCosts := make(map[int]int, 5)
	sum := 0
	for _, pos := range inCrabs {
		sum += pos
	}
	var mean int = int(math.Round(float64(sum) / float64(len(inCrabs))))
	for x := mean - 2; x <= mean+2; x++ {
		fuelCosts[x] = fuelCostNonlinear(x, inCrabs)
	}
	end = time.Since(start)
	fmt.Printf("The optimum location is around %d\n", mean)
	for x, c := range fuelCosts {
		fmt.Printf("The fuel cost at %d is %d\n", x, c)
	}

	fmt.Printf("Executed in %v\n", end)
}

func fuelCostNonlinear(pos int, crabs []int) int {
	var fuelCost int = 0
	for _, x := range crabs {
		if x < pos {
			fuelCost += (pos - x) * (pos - x + 1) / 2
		} else {
			fuelCost += (x - pos) * (x - pos + 1) / 2
		}
	}
	return fuelCost
}
