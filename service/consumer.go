package service

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type DefaultConflictConsumer struct {
	fields       []string // Conflict fields to output, empty or nil means all
	outputWriter io.Writer
	conflictsCh  <-chan Conflict
	errorCh      chan string
}

func NewConsoleConflictConsumer(conflictsChannel <-chan Conflict, format string, fields []string) <-chan string {
	return NewConflictConsumer(format, fields, os.Stdout, conflictsChannel, nil)
}

func NewConflictConsumer(format string, fields []string, outputWriter io.Writer, conflictsChannel <-chan Conflict, errorsCh chan string) <-chan string {
	if conflictsChannel == nil {
		panic(errors.New("nil conflicts Channel"))
	}
	if errorsCh == nil {
		errorsCh = make(chan string)
	}
	consumer := DefaultConflictConsumer{fields, outputWriter, conflictsChannel, errorsCh}
	go consumer.start()
	return errorsCh
}

func (p DefaultConflictConsumer) start() {
	for c := range p.conflictsCh {
		outStr, err := p.formatConflict(&c)
		if err != nil {
			p.errorCh <- fmt.Sprintf("error formatting conflict: %s", err)
			continue
		}
		fmt.Fprintf(p.outputWriter, "%s\n", outStr)
	}
	close(p.errorCh)
}

var ALL = []string{"File", "IgnoreFile", "Line", "Pattern"}

func (p DefaultConflictConsumer) formatConflict(c *Conflict) (string, error) {
	args := []any{}
	pattern := ""
	if p.fields == nil || len(p.fields) == 0 {
		p.fields = ALL
	}

	for i, f := range p.fields {
		if i > 0 {
			pattern += ", "
		}
		pattern += "%s"
		if f == "File" {
			args = append(args, c.File)
		}
		if f == "IgnoreFile" {
			args = append(args, c.IgnoreFile)
		}
		if f == "Line" {
			args = append(args, fmt.Sprint(c.Line))
		}
		if f == "Pattern" {
			args = append(args, c.Pattern)
		}
	}
	outStr := fmt.Sprintf(pattern, args...)
	return outStr, nil
}
