package service

import (
	"fmt"
	"os"
	"sync"
)

const GIT_COMMIT_FILES_COMMAND = "git diff --cached --name-only --diff-filter=ACMD"
const IGNORE_FILE = ".ilignore"

func NewService(ignoreFile string) *IgnoreService {
	fileChecker := NewFileChecker(ignoreFile, false)
	return &IgnoreService{
		fileChecker: fileChecker,
		concurrency: false,
	}
}

type IgnoreService struct {
	fileChecker FileChecker
	concurrency bool
}

// Channel to obtain list of files to be checked against ignore file
func (s IgnoreService) CheckFiles(filesChannel chan string, conflictConsumer ConflictConsumer) {
	fmt.Fprintf(os.Stderr, "CheckFiles service called\n")
	s.checkFiles(filesChannel, conflictConsumer)
}

func (s IgnoreService) CheckFilesFromStdin() {
	fmt.Fprintf(os.Stderr, "CheckFilesFromStdin service called\n")
	s.checkFiles(nil, nil)
}

func (s IgnoreService) CheckCommit() {
	fmt.Fprintf(os.Stderr, "CheckCommit service called\n")
	s.checkFiles(nil, nil)
}

func (s IgnoreService) checkFiles(filesToCheck chan string, conflictConsumer ConflictConsumer) {
	fmt.Fprintf(os.Stderr, "checkFiles service called\n")
	if !s.concurrency {
		for file := range filesToCheck {
			conflict := s.checkFile(file)
			if conflict != nil {
				conflictConsumer.ConflictsChannel() <- *conflict
			}
		}
		close(conflictConsumer.ConflictsChannel())
		for err := range conflictConsumer.ErrorChannel() {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	} else {
		// conflictsChannel := make(chan Conflict)
		var wg sync.WaitGroup
		for file := range filesToCheck {
			wg.Add(1)
			go s.checkFileToChannel(file, conflictConsumer.ConflictsChannel(), &wg)
		}
		// TODO start Conflict reader
		wg.Wait()
		fmt.Println("All go routines finished executing")
		close(conflictConsumer.ConflictsChannel())
	}
}

func (s IgnoreService) checkFile(file string) *Conflict {
	fmt.Fprintf(os.Stderr, "checkFile service called\n")
	conflict, err := s.fileChecker.checkFile(file)
	logError(err)
	return conflict
}

func (s IgnoreService) checkFileToChannel(file string, ch chan Conflict, wg *sync.WaitGroup) {
	conflict := s.checkFile(file)
	if conflict != nil {
		ch <- *conflict
	}
	wg.Done()
}

func logError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func (s IgnoreService) AddPatterns() {
	fmt.Println("AddPatterns service called")
}
