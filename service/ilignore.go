package service

import (
	"log/slog"
	"sync"
)

const GIT_COMMIT_FILES_COMMAND = "git diff --cached --name-only --diff-filter=ACMD"
const IGNORE_FILE = ".ilignore"
const DEFAULT_CONFLICTS_BUFFER_SIZE = 5

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
func (s *IgnoreService) CheckFiles(filesChannel <-chan string) <-chan Conflict {
	slog.Debug("CheckFiles service called")
	conflictsChannel := make(chan Conflict, DEFAULT_CONFLICTS_BUFFER_SIZE)
	go s.checkFiles(filesChannel, conflictsChannel)
	return conflictsChannel
}

func (s *IgnoreService) checkFiles(filesToCheck <-chan string, conflictChannel chan Conflict) {
	slog.Debug("checkFiles service called")
	if !s.concurrency {
		for file := range filesToCheck {
			conflict := s.checkFile(file)
			if conflict != nil {
				conflictChannel <- *conflict
			}
		}
		close(conflictChannel)
	} else {
		// conflictsChannel := make(chan Conflict)
		var wg sync.WaitGroup
		for file := range filesToCheck {
			wg.Add(1)
			go s.checkFileToChannel(file, conflictChannel, &wg)
		}
		// TODO start Conflict reader
		wg.Wait()
		slog.Debug("All go routines finished executing")
		close(conflictChannel)
	}
}

func (s IgnoreService) checkFile(file string) *Conflict {
	slog.Debug("checkFile service called")
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
		slog.Error("Error: %v\n", err)
	}
}

func logErrorWithMsg(msg string, err error) {
	if err != nil {
		slog.Error("%v: %v\n", msg, err)
	}
}

func (s IgnoreService) AddPatterns() {
	slog.Info("AddPatterns service called")
}
