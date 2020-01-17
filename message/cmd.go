package message

// Cmd is a message that represents a command to run
// this is a convenience message type that implements
// the x.Commander interface
type Cmd struct {
	Name string
	Args []string
}

// Command returns the command name
func (c Cmd) Command() string {
	return c.Name
}

// Arguments returns the arguments for the command
func (c Cmd) Arguments() []string {
	return c.Args
}
