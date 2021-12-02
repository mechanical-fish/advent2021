package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type SonarData struct {
	Data []int
}

type StatsResult struct {
	Increases int
	Decreases int
	Same      int
}

func (ss *SonarData) Add(num int) {
	ss.Data = append(ss.Data, num)
}

func (ss *SonarData) Diffs() StatsResult {
	r := StatsResult{}
	for i := range ss.Data {
		if i == 0 {
			continue
		}
		if ss.Data[i] > ss.Data[i-1] {
			r.Increases += 1
		} else if ss.Data[i] < ss.Data[i-1] {
			r.Decreases += 1
		} else {
			r.Same += 1
		}
	}
	return r
}

func (ss *SonarData) WindowDiffs() StatsResult {
	windowedData := &SonarData{}
	for i := range ss.Data {
		if i == 0 || i == len(ss.Data)-1 {
			continue
		}
		windowedData.Add(ss.Data[i-1] + ss.Data[i] + ss.Data[i+1])
	}
	return windowedData.Diffs()
}

func main() {
	sonar := &SonarData{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		sonar.Add(i)
	}
	diffs := sonar.Diffs()
	fmt.Printf("%d were larger\n%d were smaller\n%d were the same\n", diffs.Increases, diffs.Decreases, diffs.Same)
	winDiffs := sonar.WindowDiffs()
	fmt.Println("\nIn 3-measurement windows:")
	fmt.Printf("%d were larger\n%d were smaller\n%d were the same\n", winDiffs.Increases, winDiffs.Decreases, winDiffs.Same)
}
