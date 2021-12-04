package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// A Caller calls a sequence of numbers at times t = 0..n where t is
// an integer. We track the sequence in order as a slice, as well as
// O(1) hashmaps of numbers to time and vice versa (there is a
// one-to-one relationship between numbers and times because bingo
// numbers don't repeat).
type Caller struct {
	Sequence   []int
	NumsByTime map[int]int
	TimesByNum map[int]int
}

func NewCaller() *Caller {
	return &Caller{
		Sequence:   make([]int, 0, 100),
		NumsByTime: make(map[int]int),
		TimesByNum: make(map[int]int),
	}
}

func (c *Caller) Add(n int) {
	t := len(c.Sequence)
	c.Sequence = append(c.Sequence, n)
	c.NumsByTime[t] = n
	c.TimesByNum[n] = t
}

func (c *Caller) ParseInput(s string) error {
	for _, n := range strings.Split(s, `,`) {
		i, err := strconv.Atoi(n)
		if err != nil {
			return err
		}
		c.Add(i)
	}
	return nil
}

func (c *Caller) String() string {
	result := make([]string, len(c.Sequence))
	for t, n := range c.Sequence {
		result[t] = strconv.Itoa(n)
	}
	return fmt.Sprintf("Call Sequence:\n%s\n", strings.Join(result, ` `))
}

// A Card is a single bingo card
type Card struct {
	Grid    [][]int
	WinTime int

	caller       *Caller
	nextEmptyRow int

	// As we construct the card, because we know the Caller in advance,
	// we can compute the time t at which each row and column wins the game --
	// we compute t for each number in the row/column and track the
	// largest value of t that has been seen in each row/column
	winTimeForRow []int
	winTimeForCol []int
}

func NewCard(caller *Caller, row string) (*Card, error) {
	nums, err := parseRow(row)
	if err != nil {
		return nil, err
	}
	c := &Card{
		Grid:          make([][]int, len(nums)),
		WinTime:       0,
		caller:        caller,
		nextEmptyRow:  0,
		winTimeForRow: make([]int, len(nums)),
		winTimeForCol: make([]int, len(nums)),
	}
	c.addRowNums(nums)
	return c, nil
}

// if t is a better winning time than the one we already know,
// record t as the best winning time
func (c *Card) setWinTime(t int) {
	if c.WinTime < len(c.Grid) || t < c.WinTime {
		c.WinTime = t
	}
}

func (c *Card) addRowNums(nums []int) {
	c.Grid[c.nextEmptyRow] = make([]int, len(c.Grid))
	for i, n := range nums {
		c.Grid[c.nextEmptyRow][i] = n
		t := c.caller.TimesByNum[n]
		if t > c.winTimeForRow[c.nextEmptyRow] {
			c.winTimeForRow[c.nextEmptyRow] = t
		}
		if t > c.winTimeForCol[i] {
			c.winTimeForCol[i] = t
		}
	}
	c.nextEmptyRow++
	if c.IsComplete() {
		// we know the whole card now, so compute the time when this
		// card wins
		for _, t := range c.winTimeForRow {
			c.setWinTime(t)
		}
		for _, t := range c.winTimeForCol {
			c.setWinTime(t)
		}
	}
}

func parseRow(row string) ([]int, error) {
	fields := strings.Fields(row)
	result := make([]int, len(fields))
	var err error
	for i, s := range fields {
		if err == nil {
			result[i], err = strconv.Atoi(s)
		}
	}
	return result, err
}

func (c *Card) AddRow(row string) error {
	nums, err := parseRow(row)
	if err == nil {
		c.addRowNums(nums)
	}
	return err
}

func (c Card) IsComplete() bool {
	return c.nextEmptyRow == len(c.Grid)
}

func (c Card) Score() int {
	// if this card has not won its score is 0
	if c.WinTime < len(c.Grid) {
		return 0
	}
	// add up unmarked numbers
	sum := 0
	for _, r := range c.Grid {
		for _, n := range r {
			if c.caller.TimesByNum[n] > c.WinTime {
				sum += n
			}
		}
	}
	// multiply by the final number called
	return sum * c.caller.NumsByTime[c.WinTime]
}

// Players are a collection of cards
type Players struct {
	Cards []*Card
}

func NewPlayers() *Players {
	return &Players{
		Cards: make([]*Card, 0),
	}
}

func (p *Players) Add(c *Card) {
	p.Cards = append(p.Cards, c)
}

func (p Players) Count() int {
	return len(p.Cards)
}

// the caller input has commas
var CallerRe = regexp.MustCompile(`,`)

// the bingo cards are space-separated lists of numbers with no
// leading whitespace
var CardRowRe = regexp.MustCompile(`\A\s?\d+\s+\d`)

// the bingo cards end with a blank line
var CardEndRe = regexp.MustCompile(`\A\z`)

func main() {
	caller := NewCaller()
	players := NewPlayers()
	var card *Card
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		row := scanner.Text()
		if CallerRe.MatchString(row) {
			if err = caller.ParseInput(row); err != nil {
				panic(err)
			}
		} else if CardRowRe.MatchString(row) {
			if card == nil {
				card, err = NewCard(caller, row)
			} else {
				err = card.AddRow(row)
			}
			if err != nil {
				panic(err)
			}
		} else if CardEndRe.MatchString(row) {
			if card != nil {
				players.Add(card)
			}
			card = nil
		}
	}
	// there might not be a final blank line to trigger the saving of the
	// final card
	if card != nil {
		players.Add(card)
	}
	fmt.Println(caller)

	fmt.Printf("Found %d cards\n", players.Count())
	var winningCard int = 0
	var winTime int
	var losingestCard int = 0
	var latestWinTime int
	for i, c := range players.Cards {
		fmt.Printf("  Card %d wins after %d numbers with score %d\n",
			i+1, c.WinTime+1, c.Score())
		if winTime == 0 || c.WinTime < winTime {
			winningCard = i
			winTime = c.WinTime
		}
		if latestWinTime == 0 || c.WinTime >= latestWinTime {
			losingestCard = i
			latestWinTime = c.WinTime
		}
	}
	fmt.Printf("Card %d is the overall winner after %d numbers\n with score %d\n",
		winningCard+1, winTime+1, players.Cards[winningCard].Score())
	fmt.Printf("Card %d will win last, after %d numbers\n with score %d\n",
		losingestCard+1, latestWinTime+1, players.Cards[losingestCard].Score())
}
