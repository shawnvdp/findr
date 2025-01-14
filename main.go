package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/integrii/flaggy"
)

type cliArgs struct {
	IgnoreDirs []string
	IgnoreExts []string
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

	dirMatches := make(chan dirMatch)
	taskChan := make(chan string, 1024)
	maxWorkers := runtime.NumCPU()
	wg := &sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		go (func() {
			for dir := range taskChan {
				searchDirectory(wg, dir, cliArgs.Term, cliArgs.IgnoreDirs, cliArgs.IgnoreExts, dirMatches, taskChan)
				wg.Done()
			}
		})()
	}

	go func() {
		for match := range dirMatches {
			filePath := filepath.Join(match.Directory, match.File)
			for _, match := range match.Matches {
				fmt.Printf("%v:%v - %v\n", filePath, match.Number, match.Line)
			}
		}
	}()

	wg.Add(1)
	taskChan <- cliArgs.BaseDir

	wg.Wait()
	close(dirMatches)
	close(taskChan)

	elapsed := time.Since(start)
	fmt.Printf("Took %s to traverse %d directories and %d files", elapsed, DIRECTORIES_SEARCHED, FILES_SEARCHED)
}

func searchDirectory(wg *sync.WaitGroup, directory, term string, ignoreDirs []string, ignoreExts []string, matchChan chan dirMatch, taskChan chan string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := filepath.Join(directory, file.Name())

		if file.IsDir() {
			if len(ignoreDirs) > 0 && contains(ignoreDirs, file.Name()) {
				continue
			}
			DIRECTORIES_SEARCHED++
			wg.Add(1)
			taskChan <- filePath
			continue
		}

		if contains(ignoreExts, path.Ext(file.Name())) {
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

		matchChan <- dirMatch{Directory: directory, File: file.Name(), Matches: matches}
	}
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
	flaggy.String(&ignoreDir, "igd", "ignoreDir", "Name(s) of directories to ignore (comma-separated)")

	var ignoreExt string
	flaggy.String(&ignoreExt, "ige", "ignoreExt", "File extensions to ignore (comma-separated)")

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

	exts := strings.Split(ignoreExt, ",")
	for i, ext := range exts {
		exts[i] = strings.TrimSpace(ext)
	}

	return &cliArgs{
		IgnoreDirs: dirs,
		IgnoreExts: exts,
		BaseDir:    baseDir,
		Term:       term,
	}
}
