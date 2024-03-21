package service

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type FileSource interface {
	start()
}

type ReaderFileSource struct {
	reader          io.Reader
	outputChannel   chan string
	continueOnEmpty bool
}

// start implements FileSource.
func (p *ReaderFileSource) start() {
	defer close(p.outputChannel)
	if err := p.readLinesFromReader(); err != nil {
		slog.Error("Error reading lines", "err", err.Error(), "reader", p.reader)
	}
}

func (p *ReaderFileSource) readLinesFromReader() error {
	scanner := bufio.NewScanner(p.reader)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		slog.Debug("Read " + line)
		line = strings.TrimSpace(line)
		if line == "" {
			if p.continueOnEmpty {
				continue
			} else {
				break
			}
		}
		p.outputChannel <- line
		count++
	}
	slog.Debug("read lines", "count", count)
	return scanner.Err()
}

// list Files Command
type CommandFileSource struct {
	command string
	ReaderFileSource
}

func (p *CommandFileSource) start() {
	slog.Info("exec ", "command", p.command)
	defer close(p.outputChannel)
	cmd, err := p.startCommand()
	if err != nil {
		slog.Error("Error starting command", "err", err.Error(), "command", p.command)
		return
	}
	if err = p.readLinesFromReader(); err != nil {
		slog.Error("Error reading from command Stdout", "err", err.Error(), "command", p.command)
	}
	err = cmd.Wait()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			slog.Error("Error getting files", "err", err.Error(), "command", p.command, "exitCode", ee.ProcessState.ExitCode())
		} else {
			slog.Error("I/O Problems getting files", "err", err.Error(), "command", p.command)
		}
	}
}

func (p *CommandFileSource) startCommand() (*exec.Cmd, error) {
	fields := strings.Fields(p.command)
	cmd := exec.Command(fields[0], fields[1:]...)
	stdout, err := cmd.StdoutPipe()
	p.reader = stdout // set reader to read command output
	if err == nil {
		err = cmd.Start()
	}
	return cmd, err
}

func NewFileSource(filesToCheck []string, command string, continueOnEmpty bool) <-chan string {
	outputChannel := make(chan string)

	var producer FileSource
	if len(filesToCheck) > 0 {
		slog.Info("Check Fixed files")
		reader := strings.NewReader(strings.Join(filesToCheck, "\n"))
		producer = &ReaderFileSource{reader, outputChannel, continueOnEmpty}
	} else if command != "" {
		slog.Info("Check files from command")
		producer = &CommandFileSource{command, ReaderFileSource{nil, outputChannel, continueOnEmpty}}
	} else {
		slog.Info("Check files from Stdin")
		producer = &ReaderFileSource{os.Stdin, outputChannel, false}
	}
	go producer.start()
	return outputChannel
}

const GIT_DIFF_COMMAND_TEMPLATE = "git diff%s --name-only --diff-filter=%s"

func NewGitDiffFileSource(cached bool, diffFilter string) <-chan string {
	cachedFlag := ""
	if cached {
		cachedFlag = " --cached"
	}
	if diffFilter == "" {
		diffFilter = "ACMD"
	}
	command := fmt.Sprintf(GIT_DIFF_COMMAND_TEMPLATE, cachedFlag, diffFilter)
	return NewFileSource(nil, command, false)
}
