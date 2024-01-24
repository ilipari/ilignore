package service

type FixedFileSource struct {
	filesToCheck  []string
	outputChannel chan string
}

func NewFixedFileSource(filesToCheck []string) chan string {
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

func NewFileSourceFromCommand(command string) chan string {
	outputChannel := make(chan string)
	// fmt.Fprintf(os.Stderr, "listCommand ->%v\n", s.listFilesCommand)
	producer := CommandFileSource{command, outputChannel}
	go producer.start()
	return outputChannel
}

func (p CommandFileSource) start() {
	// TODO execute command
	close(p.outputChannel)
}
