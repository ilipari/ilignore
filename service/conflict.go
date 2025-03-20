package service

type Conflict struct {
	// offending file
	File string
	// matching pattern in .gitignore file
	Pattern string
	// ignore file with .gitignore syntax
	IgnoreFile string
	// line in ignore file containing matching pattern, if available
	Line int
}

func newConflict() *Conflict {
	return &Conflict{"", "", "", -1}
}
