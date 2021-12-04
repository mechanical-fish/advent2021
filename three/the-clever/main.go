// The clever solution constructs a trie (or trie-like?)
// data structure to solve this in linear time
package main

import (
	"bufio"
	"fmt"
	"os"
)

// Each binary number is represented by a leaf node of this tree. The
// tree is d nodes deep, where d is the number of bits in each
// number. Each node has up to two child nodes, corresponding to a 0
// or a 1.
//
// To add a number to the tree, we start at the root and take the bits
// of the number one at a time, traveling to the "zero" child for a 0 and to the "1" child for a one.
//
// Every time we pass through a node we add one to its "weight". We
// also keep count of how many nodes of each value (0 or 1) occur at
// each depth in the tree (which corresponds to the index of a digit
// in the original binary numbers.)
type Tree struct {
	Root                       *Node
	OccurrancesByDepthAndValue [][]uint
	NumNodes                   int
}

func NewTree(bitCount int) *Tree {
	t := &Tree{
		Root:                       NewNode(0),
		OccurrancesByDepthAndValue: make([][]uint, bitCount),
		NumNodes:                   1,
	}
	t.Root.Weight = 0
	for i := range t.OccurrancesByDepthAndValue {
		t.OccurrancesByDepthAndValue[i] = make([]uint, 2)
	}
	return t
}

func (t *Tree) Add(num string) {
	t.Root.Add(num, 0, t)
}

// Each leaf node has a bit value (0 or 1) plus it tracks the number
// of times this node was traversed while adding binary numbers to the
// tree. (the "weight")
type Node struct {
	Weight   int
	Bit      uint
	Children []*Node
}

func NewNode(bit uint) *Node {
	return &Node{
		Weight:   1,
		Bit:      bit,
		Children: make([]*Node, 2),
	}
}

func (n *Node) Add(num string, depth uint, t *Tree) {
	if len(num) == 0 {
		return
	}
	var nextChild uint = 0
	if num[0] == '1' {
		nextChild = 1
	}
	t.OccurrancesByDepthAndValue[depth][nextChild] += 1
	if n.Children[nextChild] == nil {
		n.Children[nextChild] = NewNode(nextChild)
		t.NumNodes += 1
	} else {
		n.Children[nextChild].Weight += 1
	}
	n.Children[nextChild].Add(num[1:], depth+1, t)
}

// Return gamma and epsilon rates
func (t Tree) Rates() (uint, uint) {
	var gamma uint = 0
	var epsilon uint = 0
	for i := range t.OccurrancesByDepthAndValue {
		gamma = 2 * gamma
		epsilon = 2 * epsilon
		if t.OccurrancesByDepthAndValue[i][1] > t.OccurrancesByDepthAndValue[i][0] {
			gamma += 1
		} else {
			epsilon += 1
		}
	}
	return gamma, epsilon
}

// Follow the tree from the root to a leaf, taking the path with the
// highest weight at each step (or the "1" path for a tie)
// and report the number at the leaf when you reach it.
func (t Tree) OxyRating() uint {
	return t.Root.Rating(0, func(children []*Node) *Node {
		if children[0].Weight > children[1].Weight {
			return children[0]
		}
		return children[1]
	})
}

// Follow the tree from the root to a leaf, taking the path with the
// lowest weight at each step (or the "0" path for a tie)
// and report the number at the leaf when you reach it.
func (t Tree) COTwoRating() uint {
	return t.Root.Rating(0, func(children []*Node) *Node {
		if children[0].Weight <= children[1].Weight {
			return children[0]
		}
		return children[1]
	})
}

func (n Node) Rating(acc uint, nextChild func([]*Node) *Node) uint {
	val := 2*acc + n.Bit
	if n.Children[0] == nil {
		if n.Children[1] == nil {
			return val
		}
		return n.Children[1].Rating(val, nextChild)
	}
	if n.Children[1] == nil {
		return n.Children[0].Rating(val, nextChild)
	}
	return nextChild(n.Children).Rating(val, nextChild)
}

func main() {
	var t *Tree
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if t == nil {
			t = NewTree(len(scanner.Text()))
		}
		t.Add(scanner.Text())
	}

	fmt.Printf("Allocated %d nodes\n\n", t.NumNodes)

	gam, eps := t.Rates()
	fmt.Printf("Gamma is %d\n", gam)
	fmt.Printf("Epsilon is %d\n", eps)
	fmt.Printf("  Product is %d\n", gam*eps)

	oxy := t.OxyRating()
	fmt.Printf("OxyRating is %d\n", oxy)
	coTwo := t.COTwoRating()
	fmt.Printf("COTwoRating is %d\n", coTwo)
	fmt.Printf("  Product is %d\n", oxy*coTwo)
}
