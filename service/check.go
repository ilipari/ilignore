package service

import (
	"fmt"
	"os"

	gitignore "github.com/denormal/go-gitignore"
)

type FileChecker interface {
	checkFile(string) (*Conflict, error)
}

const File = ".ilignore"

func NewFileChecker(ignoreFile string) FileChecker {
	ignore, err := gitignore.NewFromFile(ignoreFile)
	if err != nil {
		panic(err)
	}
	return &DenormalFileChecker{ignore}
}

type DenormalFileChecker struct {
	// ignore file with .gitignore syntax (default .ilignore)
	ignore gitignore.GitIgnore
}

func (s *DenormalFileChecker) checkFile(file string) (*Conflict, error) {
	fmt.Fprintf(os.Stderr, "DenormalFileChecker checkFile service called on file: %v\n", file)
	match := s.ignore.Match(file)
	var conflict *Conflict
	if match != nil {
		if match.Ignore() {
			conflict = NewConflict()
			conflict.File = file
			conflict.IgnoreFile = match.Position().File
			conflict.Line = match.Position().Line
			conflict.Pattern = match.String()
		} else if match.Include() {
			// in include per pattern negato
			fmt.Fprintf(os.Stderr,
				"Il file %q in include per il pattern %q alla riga %d\n",
				file, match, match.Position().Line,
			)
		}
	}
	return conflict, nil
}
