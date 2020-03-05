package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// RowAnalysis take a line from UCD and return it info
func RowAnalysis(ucdLine string) (rune, string) {
	fields := strings.Split(ucdLine, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	return rune(code), fields[1]
}

// ListUCD get a search and return findings list
func ListUCD(text io.Reader, search string) {
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		runeUCD, name := RowAnalysis(line)
		if strings.Contains(name, search) {
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runeUCD, name)
		}
	}
}

func main() {
	ucd, err := os.Open("UnicodeData.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ucd.Close()
	search := strings.Join(os.Args[1:], " ")
	ListUCD(ucd, strings.ToUpper(search))
}
