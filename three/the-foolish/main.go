package main

import (
	"bufio"
	"fmt"
	"os"
)

type Diagnostic struct {
	RawData []string
	BitSums []int
	Count   int
}

func NewDiagnostic(numBits int) *Diagnostic {
	d := Diagnostic{}
	d.BitSums = make([]int, numBits)
	d.RawData = []string{}
	return &d
}

func (d *Diagnostic) Add(num string) {
	d.Count += 1
	d.RawData = append(d.RawData, num)
	// it would probably be cleaner to do all this in binary
	// arithmetic instead of this ugly, repeated char-by-char parsing
	// but it just doesn't matter
	for i := 0; i < len(num); i++ {
		if num[len(num)-1-i] == '1' {
			d.BitSums[i] += 1
		}
	}
}

func (d *Diagnostic) Gamma() uint {
	var result uint
	for i, s := range d.BitSums {
		if 2*s > d.Count {
			result += 1 << i
		}
	}
	return result
}

func (d *Diagnostic) Epsilon() uint {
	var result uint
	for i, s := range d.BitSums {
		if 2*s < d.Count {
			result += 1 << i
		}
	}
	return result
}

func (d *Diagnostic) OxyDigitAt(idx int) byte {
	if 2*d.BitSums[idx] >= d.Count {
		return '1'
	}
	return '0'
}

func (d *Diagnostic) COTwoDigitAt(idx int) byte {
	if 2*d.BitSums[idx] < d.Count {
		return '1'
	}
	return '0'
}

func (d *Diagnostic) OxyRating(currBit int) uint {
	if len(d.RawData) == 1 {
		return binStrToNum(d.RawData[0])
	}
	keepDigit := d.OxyDigitAt(currBit)
	currIdx := len(d.BitSums) - currBit - 1
	newD := NewDiagnostic(len(d.BitSums))
	for _, num := range d.RawData {
		if num[currIdx] == keepDigit {
			newD.Add(num)
			// fmt.Printf("  digit %d, keeping %s\n", currBit, num)
		}
	}
	return newD.OxyRating(currBit - 1)
}

// Look, I'm trying to write the code fast,
// not avoid copypasta at all costs :)

func (d *Diagnostic) COTwoRating(currBit int) uint {
	if len(d.RawData) == 1 {
		return binStrToNum(d.RawData[0])
	}
	keepDigit := d.COTwoDigitAt(currBit)
	currIdx := len(d.BitSums) - currBit - 1
	newD := NewDiagnostic(len(d.BitSums))
	for _, num := range d.RawData {
		if num[currIdx] == keepDigit {
			newD.Add(num)
			// fmt.Printf("  digit %d, keeping %s\n", currBit, num)
		}
	}
	return newD.COTwoRating(currBit - 1)
}

func binStrToNum(s string) uint {
	var result uint = 0
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			result += 1 << (len(s) - i - 1)
		}
	}
	return result
}

func main() {
	var d *Diagnostic
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if d == nil {
			d = NewDiagnostic(len(scanner.Text()))
		}
		d.Add(scanner.Text())
		// fmt.Printf("Value of %s is %d\n", scanner.Text(), binStrToNum(scanner.Text()))
		// for i, s := range d.BitSums {
		// 	fmt.Printf("Sum of bit %d: %d\n", i, s)
		// }
	}

	fmt.Printf("Gamma is %d\n", d.Gamma())
	fmt.Printf("Epsilon is %d\n", d.Epsilon())
	fmt.Printf("  Product is %d\n", d.Gamma()*d.Epsilon())
	oxy := d.OxyRating(len(d.BitSums) - 1)
	fmt.Printf("OxyRating is %d\n", oxy)
	coTwo := d.COTwoRating(len(d.BitSums) - 1)
	fmt.Printf("COTwoRating is %d\n", coTwo)
	fmt.Printf("  Product is %d\n", oxy*coTwo)
}
