package service

import (
	"log/slog"
	"sync"
)

const IGNORE_FILE = ".ilignore"
const DEFAULT_CONFLICTS_BUFFER_SIZE = 5

func NewService(ignoreFiles []string, concurrency int) *IgnoreService {
	if len(ignoreFiles) > 0 {
		fileCheckers := []FileChecker{}
		slog.Info("creating ignore service", "files", ignoreFiles, "concurrency", concurrency)
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
	concurrency int
}

type checkJob struct {
	checker FileChecker
	file    string
}

// Channel to obtain list of files to be checked against ignore file
func (s *IgnoreService) CheckFiles(filesCh <-chan string) <-chan Conflict {
	slog.Debug("CheckFiles service called")
	jobsCh := s.startJobsBuilder(filesCh)
	conflictsCh := s.startJobExecutors(jobsCh)
	return conflictsCh
}

func (s *IgnoreService) startJobsBuilder(filesCh <-chan string) <-chan checkJob {
	jobs := make(chan checkJob)
	go func() {
		defer close(jobs)
		for file := range filesCh {
			for _, checker := range s.checkers {
				jobs <- checkJob{
					checker: checker,
					file:    file,
				}
			}
		}
	}()
	return jobs
}

func (s *IgnoreService) startJobExecutors(jobs <-chan checkJob) <-chan Conflict {
	conflictsCh := make(chan Conflict)
	var wg sync.WaitGroup

	// Start worker pool
	wg.Add(s.concurrency)
	for i := 0; i < s.concurrency; i++ {
		go func(idx int) {
			defer wg.Done()
			// Each worker processes jobs until channel closes
			for job := range jobs {
				slog.Debug("Checking file", "worker", idx, "file", job.file)
				conflict, err := job.checker.checkFile(job.file)
				logError(err)
				if conflict != nil {
					conflictsCh <- *conflict
				}
			}
		}(i)
	}

	// Close conflicts channel when all workers done
	go func() {
		wg.Wait()
		close(conflictsCh)
	}()

	return conflictsCh
}

func logError(err error) {
	if err != nil {
		slog.Error("Error: %v\n", err)
	}
}

func (s IgnoreService) AddPatterns() {
	slog.Info("AddPatterns service called")
}
