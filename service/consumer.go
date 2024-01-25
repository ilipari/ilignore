package service

import (
	"fmt"
	"io"
	"os"
)

type ConflictConsumer interface {
	ConflictsChannel() chan Conflict
	ErrorChannel() chan string
}

type DefaultConflictConsumer struct {
	format           string   // simple, table, csv, json
	fields           []string // Conflict fields to output, empty or nil means all
	outputWriter     io.Writer
	conflictsChannel chan Conflict
	errorChannel     chan string
}

const DEFAULT_CONFLICTS_BUFFER_SIZE = 5

func NewConsoleConflictConsumer(format string) ConflictConsumer {
	return NewConflictConsumer(format, nil, os.Stdout, nil, nil)
}

func NewConflictConsumer(format string, fields []string, outputWriter io.Writer, conflictsChannel chan Conflict, errorChannel chan string) ConflictConsumer {
	if conflictsChannel == nil {
		conflictsChannel = make(chan Conflict, DEFAULT_CONFLICTS_BUFFER_SIZE)
	}
	if errorChannel == nil {
		errorChannel = make(chan string)
	}
	consumer := DefaultConflictConsumer{format, fields, outputWriter, conflictsChannel, errorChannel}
	go consumer.start()
	return &consumer
}

func (p DefaultConflictConsumer) start() {
	for c := range p.conflictsChannel {
		fmt.Fprintf(p.outputWriter, "conflict -> %v\n", c)
	}
	close(p.errorChannel)
}

func (p *DefaultConflictConsumer) ConflictsChannel() chan Conflict {
	return p.conflictsChannel
}

func (p *DefaultConflictConsumer) ErrorChannel() chan string {
	return p.errorChannel
}
