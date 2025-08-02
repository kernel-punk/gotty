package gottylib

type RunParameters struct {
	Cmd  string
	Args []string
	Ssl  bool
	Port int
}
