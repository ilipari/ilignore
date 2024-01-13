package service

import (
	"fmt"
	"os"
)

const GIT_COMMIT_FILES_COMMAND = "git diff --cached --name-only --diff-filter=ACMD"
const IGNORE_FILE = ".ilignore"

func NewService(listFilesCommand, ignoreFile string) *IgnoreService {
	return &IgnoreService{
		listFilesCommand: listFilesCommand,
		ignoreFile:       ignoreFile,
	}
}

type IgnoreService struct {
	// command to obtain list of files to be checked against ignore file
	listFilesCommand string
	// ignore file with .gitignore syntax (default .ilignore)
	ignoreFile string
}

func (s IgnoreService) CheckFiles() []Conflict {
	fmt.Fprintf(os.Stderr, "CheckFiles service called\n")
	fmt.Fprintf(os.Stderr, "listCommand ->%v\n", s.listFilesCommand)
	fmt.Fprintf(os.Stderr, "ignore file ->%v\n", s.ignoreFile)
	return s.checkFiles(nil, nil)
}

func (s IgnoreService) CheckFilesFromStdin() []Conflict {
	fmt.Fprintf(os.Stderr, "CheckFilesFromStdin service called\n")
	return s.checkFiles(nil, nil)
}

func (s IgnoreService) CheckCommit() []Conflict {
	fmt.Fprintf(os.Stderr, "CheckCommit service called\n")
	return s.checkFiles(nil, nil)
}

func (s IgnoreService) checkFiles(filesToCheck, ignoreFiles []string) []Conflict {
	fmt.Fprintf(os.Stderr, "checkFiles service called\n")
	return nil
}

func (s IgnoreService) AddPatterns() {
	fmt.Println("AddPatterns service called")
}
