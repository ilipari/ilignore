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

func NewConsoleConflictConsumer(conflictsChannel <-chan Conflict, format string) <-chan string {
	return NewConflictConsumer(format, nil, os.Stdout, conflictsChannel, nil)
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
		fmt.Fprintf(p.outputWriter, "conflict -> %v\n", c)
	}
	close(p.errorChannel)
}
