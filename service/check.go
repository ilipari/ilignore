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

func NewFileChecker(ignoreFile string, log bool) FileChecker {
	fileChecker, err :=
		NewDenormalFileChecker(
			ignoreFile,
			func(e gitignore.Error) bool {
				if log {
					logError(e.Underlying())
				}
				return true
			})
	if err != nil {
		panic(err)
	}
	return fileChecker
}

// NewDenormalFileChecker creates a FileChecker instance from the given file. An error
// will be returned if file cannot be opened or its absolute path determined.
func NewDenormalFileChecker(file string, errorHandler func(e gitignore.Error) bool) (FileChecker, error) {
	// define an error handler to catch any file access errors
	//		- record the first encountered error
	var err gitignore.Error
	var init = true
	wrapErrorHandler := func(e gitignore.Error) bool {
		if !init && errorHandler != nil {
			return errorHandler(e)
		}
		if err == nil {
			err = e
		}
		return true
	}

	// attempt to retrieve the GitIgnore represented by this file
	ignore := gitignore.NewWithErrors(file, wrapErrorHandler)

	// did we encounter an error?
	//		- if the error has a zero Position then it was encountered
	//		  before parsing was attempted, so we return that error
	if err != nil {
		if err.Position().Zero() {
			return nil, err.Underlying()
		}
	}

	// otherwise, we ignore the parser errors on .gitignore file
	init = false
	return &DenormalFileChecker{ignore}, nil
} // NewDenormalFileChecker()

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
