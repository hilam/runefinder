package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// UCDURL is the url to download de UnicodeDat.txt
const UCDURL = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"

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
	count := 0
	searchs := strings.Fields(search)
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		runeUCD, name, words := RowAnalysis(line)
		if containsAll(words, searchs) {
			count++
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runeUCD, name)
		}
	}
	fmt.Println("# of results: " + strconv.FormatInt(int64(count), 10))
}

func restore(envVar, value string, existed bool) {
	if existed {
		os.Setenv(envVar, value)
	} else {
		os.Unsetenv(envVar)
	}
}

func getPathUCD() string {
	pathUCD := os.Getenv("UCD_PATH")
	if pathUCD == "" {
		user, err := user.Current()
		exitIf(err)
		pathUCD = user.HomeDir + "/UnicodeData.txt"
	}
	return pathUCD
}

func exitIf(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func openUCD(path string) (*os.File, error) {
	ucd, err := os.Open(path)
	if os.IsNotExist(err) {
		fmt.Printf("%s not found\nDownloading %s\n", path, UCDURL)
		done := make(chan bool)
		go downloadUCD(UCDURL, path, done)
		progress(done)
		ucd, err = os.Open(path)
	}
	return ucd, err
}

func progress(done <-chan bool) {
	for {
		select {
		case <-done:
			fmt.Println()
			return
		default:
			fmt.Print(".")
			time.Sleep(150 * time.Millisecond)
		}
	}
}

func downloadUCD(url, path string, done chan<- bool) {
	response, err := http.Get(url)
	exitIf(err)
	defer response.Body.Close()
	file, err := os.Create(path)
	exitIf(err)
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	exitIf(err)
	done <- true
}

func main() {
	ucd, err := openUCD(getPathUCD())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ucd.Close()
	search := strings.Join(os.Args[1:], " ")
	ListUCD(ucd, strings.ToUpper(search))
}
