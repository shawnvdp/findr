package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/integrii/flaggy"
)

type cliArgs struct {
	IgnoreDir string
	BaseDir   string
	Term      string
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

func main() {
	cliArgs := parseCommandLineArguments()

	if len(cliArgs.Term) <= 0 {
		log.Fatal("Please specify a search term (-term <some_term>)")
	}

	dirMatches := searchDirectory(cliArgs.BaseDir, cliArgs.Term, cliArgs.IgnoreDir)

	for _, dirMatch := range dirMatches {
		filePath := filepath.Join(dirMatch.Directory, dirMatch.File)
		for _, match := range dirMatch.Matches {
			fmt.Printf("%v:%v - %v\n", filePath, match.Number, match.Line)
		}
	}
}

func searchDirectory(directory, term, ignoreDir string) []dirMatch {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	var dirMatches []dirMatch
	for _, file := range files {
		if file.IsDir() {
			if strings.Contains(file.Name(), ignoreDir) {
				continue
			}
			// todo: implement recursion
			continue
		}
		fmt.Println(file.Name(), file.IsDir())

		filePath := filepath.Join(directory, file.Name())

		contents, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		matches := scanFileForTerm(contents, term)

		if len(matches) <= 0 {
			continue
		}

		dirMatches = append(dirMatches, dirMatch{Directory: directory, File: file.Name(), Matches: matches})
	}

	return dirMatches

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
	flaggy.String(&ignoreDir, "ignoreDir", "ignore", "Name of directories to ignore")

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var baseDir string = currentDir
	flaggy.String(&baseDir, "baseDir", "base", "Base directory path to start from (defaults to current directory)")

	var term string
	flaggy.String(&term, "term", "search term", "Term to search")

	flaggy.Parse()

	return &cliArgs{
		IgnoreDir: ignoreDir,
		BaseDir:   baseDir,
		Term:      term,
	}
}
