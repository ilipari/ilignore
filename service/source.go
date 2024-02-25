package service

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type StdinFileSource struct {
	outputChannel chan string
}

func NewStdinFileSource() <-chan string {
	outputChannel := make(chan string)
	producer := StdinFileSource{outputChannel}
	go producer.start()
	return outputChannel
}

func (p StdinFileSource) start() {
	if err := readLines(os.Stdin, p.outputChannel, false); err != nil {
		slog.Error("Error in reading from STDIN", "err", err.Error())
	}
	close(p.outputChannel)
}

func readLines(in io.Reader, out chan string, continueOnEmpty bool) error {
	scanner := bufio.NewScanner(in)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		slog.Debug("Read " + line)
		line = strings.TrimSpace(line)
		if line == "" {
			if continueOnEmpty {
				continue
			} else {
				break
			}
		}
		out <- line
		count++
	}
	slog.Debug("read lines", "count", count)
	return scanner.Err()
}

// list Files Command
type CommandFileSource struct {
	command string
	StdinFileSource
}

func NewCommandFileSource(command string) <-chan string {
	outputChannel := make(chan string)
	// log.Printf("listCommand ->%v\n", command)
	producer := CommandFileSource{command, StdinFileSource{outputChannel}}
	go producer.start()
	return outputChannel
}

func (p CommandFileSource) start() {
	slog.Info("exec ", "command", p.command)
	cmd := exec.Command("bash", "-c", p.command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error starting command", "err", err.Error(), "command", p.command)
	} else {
		if err = cmd.Start(); err != nil {
			slog.Error(err.Error())
		} else {
			if err = readLines(stdout, p.outputChannel, true); err != nil {
				slog.Error("Error reading from command Stdout", "err", err.Error(), "command", p.command)
			}
			err = cmd.Wait()
			if err != nil {
				slog.Error("Error Waiting for command to exit", "err", err.Error(), "command", p.command)
			}
		}

	}
	close(p.outputChannel)
}

// for debug/test purpose
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
