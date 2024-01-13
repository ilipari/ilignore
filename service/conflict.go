package service

type Conflict struct {
	// offending file command to obtain list of files to be checked against ignore file
	File string
	// matching pattern in .gitignore file
	Pattern string
	// ignore file with .gitignore syntax
	IgnoreFile string
	// line in ignore file containing matching pattern, if available
	Line int
}

func NewConflict() *Conflict {
	return &Conflict{"", "", "", -1}
}
