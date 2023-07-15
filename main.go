package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/integrii/flaggy"
)

type cliArgs struct {
	IgnoreDirs []string
	BaseDir    string
	Term       string
}

type dirMatch struct {
	Directory string
	File      string
	Matches   []match
}

type match struct {
	Line   string
	Number int
}

var (
	DIRECTORIES_SEARCHED int
	FILES_SEARCHED       int
)

func main() {
	cliArgs := parseCommandLineArguments()

	if len(cliArgs.Term) <= 0 {
		log.Fatal("Please specify a search term (-term <some_term>)")
	}

	start := time.Now()
	dirMatches := searchDirectory(cliArgs.BaseDir, cliArgs.Term, cliArgs.IgnoreDirs)

	for _, dirMatch := range dirMatches {
		filePath := filepath.Join(dirMatch.Directory, dirMatch.File)
		for _, match := range dirMatch.Matches {
			fmt.Printf("%v:%v - %v\n", filePath, match.Number, match.Line)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s to traverse %d directories and %d files", elapsed, DIRECTORIES_SEARCHED, FILES_SEARCHED)
}

func searchDirectory(directory, term string, ignoreDirs []string) []dirMatch {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	var dirMatches []dirMatch
	for _, file := range files {
		filePath := filepath.Join(directory, file.Name())

		if file.IsDir() {
			if len(ignoreDirs) > 0 && contains(ignoreDirs, file.Name()) {
				continue
			}
			DIRECTORIES_SEARCHED++
			dirMatches = append(dirMatches, searchDirectory(filePath, term, ignoreDirs)...)
			continue
		}

		contents, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		FILES_SEARCHED++

		matches := scanFileForTerm(contents, term)

		if len(matches) <= 0 {
			continue
		}

		dirMatches = append(dirMatches, dirMatch{Directory: directory, File: file.Name(), Matches: matches})
	}

	return dirMatches
}

func contains(arr []string, target string) bool {
	for _, el := range arr {
		if el == target {
			return true
		}
	}
	return false
}

func scanFileForTerm(contents []byte, term string) []match {
	lines := strings.Split(string(contents), "\n")
	var matches []match

	for i, line := range lines {
		if !strings.Contains(line, term) {
			continue
		}

		idx := strings.Index(line, term)
		lowerBound := Max(idx-len(term)-10, 0)
		upperBound := Min(idx+len(term)+10, len(line))
		substr := line[lowerBound:upperBound]

		matches = append(matches, match{Line: substr, Number: i})
	}

	return matches
}

func parseCommandLineArguments() *cliArgs {
	var ignoreDir string
	flaggy.String(&ignoreDir, "ignore", "ignoreDir", "Name(s) of directories to ignore (comma-separated)")

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var baseDir string = currentDir
	flaggy.String(&baseDir, "base", "baseDir", "Base directory path to start from")

	var term string
	flaggy.String(&term, "term", "searchTerm", "Term to search")

	flaggy.Parse()

	dirs := strings.Split(ignoreDir, ",")

	for i, dir := range dirs {
		dirs[i] = strings.TrimSpace(dir)
	}

	return &cliArgs{
		IgnoreDirs: dirs,
		BaseDir:    baseDir,
		Term:       term,
	}
}
