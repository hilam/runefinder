package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

const letterA = `0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;`

const line3Dto43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestRowAnalysis_cases(t *testing.T) {
	var cases = []struct {
		line    string
		runeUCD rune
		name    string
		words   []string
	}{
		{letterA,
			'A',
			"LATIN CAPITAL LETTER A",
			[]string{"LATIN", "CAPITAL", "LETTER", "A"},
		},
		{"0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;",
			'!',
			"EXCLAMATION MARK",
			[]string{"EXCLAMATION", "MARK"},
		},
		// {"002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;",
		// 	'-',
		// 	"HYPHEN-MINUS",
		// 	[]string{"HYPHEN", "MINUS"},
		// },
		// {"0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;",
		// 	'\'',
		// 	"APOSTROPHE (APOSTROPHE-QUOTE)",
		// 	[]string{"APOSTROPHE", "QUOTE"},
		// },
	}

	for _, caseX := range cases {
		runeUCD, name, words := RowAnalysis(caseX.line)
		if runeUCD != caseX.runeUCD ||
			name != caseX.name ||
			!reflect.DeepEqual(words, caseX.words) {
			t.Errorf("\nAnalisarLinha(%q)\n-> (%q, %q, %q)", // ➎
				caseX.line, runeUCD, name, words)
		}
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

func TestContains(t *testing.T) {
	var cases = []struct {
		slice  []string
		search string
		want   bool
	}{
		{[]string{"A", "B"}, "B", true},
		{[]string{}, "A", false},
		{[]string{"A", "B"}, "Z", false},
	}

	for _, caseX := range cases {
		got := contains(caseX.slice, caseX.search)
		if got != caseX.want {
			t.Errorf("contains(%#v, %#v) want: %v; got: %v",
				caseX.slice, caseX.search, caseX.want, got)
		}
	}
}

func TestContainsAll(t *testing.T) {
	var cases = []struct {
		slice   []string
		searchs []string
		want    bool
	}{
		{[]string{"A", "B"}, []string{"B"}, true},
		{[]string{}, []string{"A"}, false},
		{[]string{"A"}, []string{}, true},
		{[]string{}, []string{}, true},
		{[]string{"A", "B"}, []string{"Z"}, false},
		{[]string{"A", "B", "C"}, []string{"A", "C"}, true},
		{[]string{"A", "B", "C"}, []string{"A", "Z"}, false},
		{[]string{"A", "B"}, []string{"A", "B", "C"}, false},
	}

	for _, caseX := range cases {
		got := containsAll(caseX.slice, caseX.searchs)
		if got != caseX.want {
			t.Errorf("containsAll(%#v, %#v) want: %v; got: %v",
				caseX.slice, caseX.searchs, caseX.want, got)
		}
	}
}

func ExampleListUCD_two_words() {
	text := strings.NewReader(line3Dto43)
	ListUCD(text, "CAPITAL LATIN")
	// Output:
	// U+0041	A	LATIN CAPITAL LETTER A
	// U+0042	B	LATIN CAPITAL LETTER B
	// U+0043	C	LATIN CAPITAL LETTER C
}

func Example_searchTwoWords() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "cat ", "smiling"}
	main()
	// Output:
	// U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
	// U+1F63A	😺	SMILING CAT FACE WITH OPEN MOUTH
	// U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES

}
