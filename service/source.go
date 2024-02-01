package service

import (
	"log/slog"
	"os/exec"
	"strings"
)

type FixedFileSource struct {
	filesToCheck  []string
	outputChannel chan string
}

func NewFixedFileSource(filesToCheck []string) <-chan string {
	outputChannel := make(chan string)
	producer := FixedFileSource{filesToCheck, outputChannel}
	go producer.start()
	return outputChannel
}

func (p FixedFileSource) start() {
	for _, f := range p.filesToCheck {
		p.outputChannel <- f
	}
	close(p.outputChannel)
}

// list Files Command
type CommandFileSource struct {
	command       string
	outputChannel chan string
}

func NewFileSourceFromCommand(command string) <-chan string {
	outputChannel := make(chan string)
	// log.Printf("listCommand ->%v\n", command)
	producer := CommandFileSource{command, outputChannel}
	go producer.start()
	return outputChannel
}

func (p CommandFileSource) start() {
	files, error := getFilesToCommit(p.command)
	if error == nil {
		slog.Info("files command returned files", "count", len(files))
		for _, file := range files {
			p.outputChannel <- file
		}
	} else {
		logErrorWithMsg("Error executing command", error)
	}
	close(p.outputChannel)
}

func getFilesToCommit(getFilesCommand string) ([]string, error) {
	slog.Info("exec ", "command", getFilesCommand)
	cmd := exec.Command("bash", "-c", getFilesCommand)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// Converti l'output in una stringa e rimuovi l'ultimo carattere di new line
	outputStr := strings.TrimSuffix(string(output), "\n")
	var files []string
	if len(outputStr) > 0 {
		files = strings.Split(outputStr, "\n")
	}
	return files, nil
}
