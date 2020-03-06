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

func contains(slice []string, search string) bool {
	for _, item := range slice {
		if item == search {
			return true
		}
	}
	return false
}

func containsAll(slice []string, searchs []string) bool {
	for _, search := range searchs {
		if !contains(slice, search) {
			return false
		}
	}
	return true
}

func split(term string) []string {
	separator := func(c rune) bool {
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(term, separator)
}

// RowAnalysis take a line from UCD and return it info
func RowAnalysis(ucdLine string) (rune, string, []string) {
	fields := strings.Split(ucdLine, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	words := split(name)

	if fields[10] != "" {
		name = fmt.Sprintf("%s (%s)", name, fields[10])
		for _, term := range split(fields[10]) {
			words = append(words, term)
		}
	}
	return rune(code), name, words
}

// ListUCD get a search and return findings list
func ListUCD(text io.Reader, search string) {
	searchs := strings.Fields(search)
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		runeUCD, name, words := RowAnalysis(line)
		if containsAll(words, searchs) {
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
