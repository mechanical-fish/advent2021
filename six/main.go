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
	buckets := make([]int, 9)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		for _, n := range strings.Split(scanner.Text(), `,`) {
			i, err := strconv.Atoi(n)
			if err != nil {
				panic(err)
			}
			buckets[i]++
		}
	}
	printBuckets(buckets, 0)
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
	printBuckets(buckets, endTime)
	fmt.Printf("Executed in %v", time.Since(start))
}
