package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func printBuckets(buckets []int, t int) {
	fmt.Printf("t=%d:\t", t)
	total := 0
	for _, n := range buckets {
		fmt.Printf("%d\t", n)
		total += n
	}
	fmt.Printf(" -- total %d\n", total)
}

func main() {
	inBuckets := make([]int, 9)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		for _, n := range strings.Split(scanner.Text(), `,`) {
			i, err := strconv.Atoi(n)
			if err != nil {
				panic(err)
			}
			inBuckets[i]++
		}
	}
	printBuckets(inBuckets, 0)

	buckets := make([]int, 9)
	for i, n := range inBuckets {
		buckets[i] = n
	}

	// first, a fast but not ultimately-fast algorithm:

	endTime := 256
	start := time.Now()
	for t := 1; t <= endTime; t++ {
		births := buckets[0]
		for i := range buckets {
			if i == 8 {
				buckets[i] = births
				continue
			}
			buckets[i] = buckets[i+1]
		}
		buckets[6] += births
	}
	end := time.Since(start)
	printBuckets(buckets, endTime)
	fmt.Printf("Executed the 'slow' way in %v\n", end)

	// Now the even faster way, with fewer copies --
	// this has successfully made me miss assembly language programming
	for i, n := range inBuckets {
		buckets[i] = n
	}
	start = time.Now()
	// why copy all those buckets? they mostly just cascade, so
	// just move a cursor which always points to the bucket that is 0 at time t
	cursor := 0
	for t := 1; t <= endTime; t++ {
		buckets[(cursor+7)%9] += buckets[cursor]
		cursor = (cursor + 1) % 9
	}
	end = time.Since(start)
	printBuckets(buckets, endTime)
	fmt.Printf("Executed the 'improved' way in %v\n", end)
}
