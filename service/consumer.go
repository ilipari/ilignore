package service

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type DefaultConflictConsumer struct {
	format           string   // simple, table, csv, json
	fields           []string // Conflict fields to output, empty or nil means all
	outputWriter     io.Writer
	conflictsChannel <-chan Conflict
	errorChannel     chan string
}

func NewConsoleConflictConsumer(conflictsChannel <-chan Conflict, format string, fields []string) <-chan string {
	return NewConflictConsumer(format, fields, os.Stdout, conflictsChannel, nil)
}

func NewConflictConsumer(format string, fields []string, outputWriter io.Writer, conflictsChannel <-chan Conflict, errorChannel chan string) <-chan string {
	if conflictsChannel == nil {
		panic(errors.New("nil conflicts Channel"))
	}
	if errorChannel == nil {
		errorChannel = make(chan string)
	}
	consumer := DefaultConflictConsumer{format, fields, outputWriter, conflictsChannel, errorChannel}
	go consumer.start()
	return errorChannel
}

func (p DefaultConflictConsumer) start() {
	for c := range p.conflictsChannel {
		outStr := p.formatConflict(&c)
		fmt.Fprintf(p.outputWriter, "conflict -> %s\n", outStr)
	}
	close(p.errorChannel)
}

var ALL = []string{"File", "IgnoreFile", "Line", "Pattern"}

func (p DefaultConflictConsumer) formatConflict(c *Conflict) string {
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
	return outStr
}
