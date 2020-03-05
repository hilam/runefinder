package main

import (
	"os"
	"strings"
	"testing"
)

const letterA = `0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;`
const nameA = "LATIN CAPITAL LETTER A"
const line3Dto43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestRowAnalysis(t *testing.T) {
	runeUCD, name := RowAnalysis(letterA)
	if runeUCD != 'A' {
		t.Errorf("Want 'A', got %q", runeUCD)
	}
	if name != nameA {
		t.Errorf("Want %q, got %q", nameA, name)
	}
}

func ExampleListUCD() {
	text := strings.NewReader(line3Dto43)
	ListUCD(text, "MARK")
	// Output: U+003F	?	QUESTION MARK
}

func ExampleListUCD_two_results() {
	text := strings.NewReader(line3Dto43)
	ListUCD(text, "SIGN")
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}

func Example() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "cruzeiro"}
	main()
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
}
