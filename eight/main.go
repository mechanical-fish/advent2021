package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strings"
)

// Return the binary representation of a given set of signals
func signalsToBinary(s string) uint {
	var result uint = 0
	for _, c := range s {
		result += uint(1 << (c - 'a'))
	}
	return result
}

type Example struct {
	// ciphertext for example digits, by length
	L map[int][]uint

	// ciphertext for output digits
	Outputs []uint

	// ciphertext representation of each cleartext digit 0-9
	Key []uint
}

func NewExample(s string) *Example {
	e := &Example{
		L:       make(map[int][]uint),
		Outputs: make([]uint, 4),
		Key:     make([]uint, 10),
	}
	parts := strings.Split(s, "|")
	for _, p := range strings.Fields(parts[0]) {
		b := signalsToBinary(p)
		if e.L[len(p)] == nil {
			e.L[len(p)] = []uint{b}
		} else {
			e.L[len(p)] = append(e.L[len(p)], b)
		}
	}
	for i, p := range strings.Fields(parts[1]) {
		e.Outputs[i] = signalsToBinary(p)
	}
	return e
}

// Return the ciphertext for the only example digit with N segments.
// If more than one example digit has N segments, panic.
func (e *Example) CipherWithSegments(n int) uint {
	if len(e.L[n]) != 1 {
		panic("example has wrong segments")
	}
	return e.L[n][0]
}

// figure out the key for digits which have a unique number of segments
func (e *Example) ComputeSimpleKey() {
	e.Key[1] = e.CipherWithSegments(2)
	e.Key[4] = e.CipherWithSegments(4)
	e.Key[7] = e.CipherWithSegments(3)
	e.Key[8] = e.CipherWithSegments(7)
}

// do the math to figure out the ciphertext representation of each
// cleartext digit
func (e *Example) ComputeKey() {
	e.ComputeSimpleKey()
	// here we do a whole bunch of boolean algebra, as worked out
	// on paper, to get ciphertext for segments a to g
	Sa := e.Key[7] ^ e.Key[1]
	var fivesPairwise []uint = []uint{
		e.L[5][0] ^ e.L[5][1],
		e.L[5][0] ^ e.L[5][2],
		e.L[5][1] ^ e.L[5][2],
	}
	var allFives uint = 0
	for _, n := range fivesPairwise {
		allFives = allFives | n
	}
	var sixesPairwise []uint = []uint{
		e.L[6][0] ^ e.L[6][1],
		e.L[6][0] ^ e.L[6][2],
		e.L[6][1] ^ e.L[6][2],
	}
	var allSixes uint = 0
	for _, n := range sixesPairwise {
		allSixes = allSixes | n
	}
	Sc := allFives ^ allSixes ^ e.Key[4]
	Sf := Sc ^ e.Key[1]
	var Sb uint
	var Se uint
	// fmt.Printf("%v\n", e.L)
	for _, n := range fivesPairwise {
		// wow i love golang's stock "bits" package
		// fmt.Printf("%b ^ %b is %b with %d ones\n", n, Sc, n^Sc, bits.OnesCount(n^Sc))
		if bits.OnesCount(n^Sc) == 1 {
			Sb = n ^ Sc
		}
		if bits.OnesCount(n^Sf) == 1 {
			Se = n ^ Sf
		}
	}
	Sd := Sb ^ Sf ^ allFives ^ allSixes
	Sg := e.Key[8] ^ Sa ^ Sb ^ Sc ^ Sd ^ Se ^ Sf

	// fmt.Printf("Segment a is %b\n", Sa)
	// fmt.Printf("Segment b is %b\n", Sb)
	// fmt.Printf("Segment c is %b\n", Sc)
	// fmt.Printf("Segment d is %b\n", Sd)
	// fmt.Printf("Segment e is %b\n", Se)
	// fmt.Printf("Segment f is %b\n", Sf)
	// fmt.Printf("Segment g is %b\n", Sb)

	// finally we compute the remaining ciphertexts in the key
	// in terms of the ciphertext representation of individual segments
	e.Key[0] = Sa | Sb | Sc | Se | Sf | Sg
	e.Key[2] = Sa | Sc | Sd | Se | Sg
	e.Key[3] = Sa | Sc | Sd | Sf | Sg
	e.Key[5] = Sa | Sb | Sd | Sf | Sg
	e.Key[6] = Sa | Sb | Sd | Se | Sf | Sg
	e.Key[9] = Sa | Sb | Sc | Sd | Sf | Sg
}

// count the number of decodable outputs, for part one
func (e *Example) CountKnownOutputs() int {
	c := 0
	for _, o := range e.Outputs {
		for _, k := range e.Key {
			if k == o {
				c++
			}
		}
	}
	return c
}

// return the number displayed by the decoded outputs
func (e *Example) Output() float64 {
	var result float64
	for i, o := range e.Outputs {
		for d, k := range e.Key {
			if k == o {
				result += math.Pow10(3-i) * float64(d)
			}
		}
	}
	return result
}

func main() {
	// The given data is groups of signals.
	// Each group represents a 7-segment digit.
	// Represent a group of signals as a 7-digit binary number,
	// with a = 1 and g = 2^7 = 64.
	examples := []*Example{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		examples = append(examples, NewExample(scanner.Text()))
	}

	// part one
	total := 0
	for i, e := range examples {
		fmt.Printf("Example %d:\n", i)
		e.ComputeSimpleKey()
		for i, k := range e.Key {
			if k != 0 {
				fmt.Printf(" digit %d is %b\n", i, k)
			}
		}
		c := e.CountKnownOutputs()
		fmt.Printf(" we know %d outputs already\n", c)
		total += c
	}
	fmt.Printf("\nwe know a total of %d outputs in part one\n\n", total)

	// part two
	var sum float64
	for i, e := range examples {
		e.ComputeKey()
		// for d, k := range e.Key {
		// 	fmt.Printf(" Digit %d represented by %b\n", d, k)
		// }
		o := e.Output()
		fmt.Printf("Output %d is %d\n", i, int(o))
		sum += o
	}
	fmt.Printf("\nTotal is %d\n", int(sum))
}
