package service

import (
	"log/slog"
)

const IGNORE_FILE = ".ilignore"
const DEFAULT_CONFLICTS_BUFFER_SIZE = 5

func NewService(ignoreFiles []string, concurrency bool) *IgnoreService {
	if len(ignoreFiles) > 0 {
		fileCheckers := []FileChecker{}
		for _, ignoreFile := range ignoreFiles {
			fileChecker := NewFileChecker(ignoreFile, false)
			fileCheckers = append(fileCheckers, fileChecker)
		}
		return &IgnoreService{
			checkers:    fileCheckers,
			concurrency: concurrency,
		}
	}
	panic("At least one ignore file is required")
}

type IgnoreService struct {
	checkers    []FileChecker
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
			for _, fileChecker := range s.checkers {
				conflict, err := fileChecker.checkFile(file)
				logError(err)
				if conflict != nil {
					conflictChannel <- *conflict
				}
			}
		}
		close(conflictChannel)
	} else {
		// TODO
		slog.Error("Unimplemented")
	}
}

func logError(err error) {
	if err != nil {
		slog.Error("Error: %v\n", err)
	}
}

func (s IgnoreService) AddPatterns() {
	slog.Info("AddPatterns service called")
}
