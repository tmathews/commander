package commander

import "fmt"

type Handler func(string, []string) error
type M map[string]Handler

type Error interface {
	error
	Help() string
}

type cmdError struct {
	x  string
	xs []string
}

func (c *cmdError) Error() string {
	return "command not found"
}

func (c *cmdError) Help() string {
	prefix := ""
	if c.x != "" && c.x != "help" {
		prefix = fmt.Sprintf("Invalid command '%s'. ", c.x)
	}
	options := fmt.Sprintf("%sPossible commands\n", prefix)
	for _, cmd := range c.xs {
		options += fmt.Sprintf("	%s\n", cmd)
	}
	options += "\nPraise the sun!\n"
	return options
}

func newError(commands M, command string) *cmdError {
	err := cmdError{
		command,
		[]string{},
	}
	for k, _ := range commands {
		err.xs = append(err.xs, k)
	}
	return &err
}

func Exec(commands M, args []string) error {
	var cmd string
	if len(args) >= 1 {
		cmd = args[0]
	} else {
		return newError(commands, cmd)
	}
	if fn, ok := commands[cmd]; !ok {
		return newError(commands, cmd)
	} else {
		if err := fn(cmd, args[1:]); err != nil {
			return err
		}
	}
	return nil
}

