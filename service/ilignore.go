package service

import (
	"fmt"
	"os"
)

const GIT_COMMIT_FILES_COMMAND = "git diff --cached --name-only --diff-filter=ACMD"
const IGNORE_FILE = ".ilignore"

func NewService(listFilesCommand, ignoreFile string) *IgnoreService {
	fileChecker := NewFileChecker(ignoreFile, false)
	return &IgnoreService{
		listFilesCommand: listFilesCommand,
		fileChecker:      fileChecker,
		concurrency:      false,
	}
}

type IgnoreService struct {
	// command to obtain list of files to be checked against ignore file
	listFilesCommand string
	fileChecker      FileChecker
	concurrency      bool
}

func (s IgnoreService) CheckFiles() []Conflict {
	fmt.Fprintf(os.Stderr, "CheckFiles service called\n")
	fmt.Fprintf(os.Stderr, "listCommand ->%v\n", s.listFilesCommand)
	// fmt.Fprintf(os.Stderr, "ignore file ->%v\n", s.ignoreFile)
	files := []string{"ciao.txt", "mondo.csv", ".vscode"}
	return s.checkFiles(files)
}

func (s IgnoreService) CheckFilesFromStdin() []Conflict {
	fmt.Fprintf(os.Stderr, "CheckFilesFromStdin service called\n")
	return s.checkFiles(nil)
}

func (s IgnoreService) CheckCommit() []Conflict {
	fmt.Fprintf(os.Stderr, "CheckCommit service called\n")
	return s.checkFiles(nil)
}

func (s IgnoreService) checkFiles(filesToCheck []string) []Conflict {
	fmt.Fprintf(os.Stderr, "checkFiles service called\n")
	var conflicts []Conflict
	if !s.concurrency {
		for _, file := range filesToCheck {
			conflict := s.checkFile(file)
			if conflict != nil {
				conflicts = append(conflicts, *conflict)
			}
		}
	} else {
		conflictsOutputChannel := make(chan Conflict)
		for _, file := range filesToCheck {
			// TODO
			go s.checkFileToChannel(file, conflictsOutputChannel)
		}
	}
	return conflicts
}

func (s IgnoreService) checkFile(file string) *Conflict {
	fmt.Fprintf(os.Stderr, "checkFile service called\n")
	conflict, err := s.fileChecker.checkFile(file)
	logError(err)
	return conflict
}

func (s IgnoreService) checkFileToChannel(file string, ch chan Conflict) {
	conflict := s.checkFile(file)
	if conflict != nil {
		ch <- *conflict
	}
}

func logError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func (s IgnoreService) AddPatterns() {
	fmt.Println("AddPatterns service called")
}
