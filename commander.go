package cmd

import "fmt"

type Helper func(string, []string) string
type Handler func(string, []string) error
type M map[string]Handler

type Error interface {
	error
	Help() string
}

type cmdError struct {
	x  string
	xs []string
	fn Helper
}

func (c *cmdError) Error() string {
	return "command not found"
}

func (c *cmdError) Help() string {
	return c.fn(c.x, c.xs)
}

func newError(commands M, command string, helper Helper) *cmdError {
	err := cmdError{
		command,
		[]string{},
		helper,
	}
	for k, _ := range commands {
		err.xs = append(err.xs, k)
	}
	return &err
}

// Looks up the first command in args and applies the rest to the found function
// in M. An error is returned if the function fails or the command is not found.
func Exec(args []string, helper Helper, commands M) error {
	var cmd string
	if len(args) >= 1 {
		cmd = args[0]
	} else {
		return newError(commands, cmd, helper)
	}
	if fn, ok := commands[cmd]; !ok {
		return newError(commands, cmd, helper)
	} else {
		if err := fn(cmd, args[1:]); err != nil {
			return err
		}
	}
	return nil
}

func isInvalid(cmd string) bool {
	return cmd != "" && cmd != "help" && cmd != "?"
}

// The default help information providing a list of commands you can run.
func DefaultHelper(cmd string, cmds []string) string {
	prefix := ""
	if isInvalid(cmd) {
		prefix = fmt.Sprintf("Invalid command '%s'. ", cmd)
	}
	options := fmt.Sprintf("%sPossible commands\n", prefix)
	for _, cmd := range cmds {
		options += fmt.Sprintf("	%s\n", cmd)
	}
	return options
}

// This creates a simple manual by introducing the command name, listing it's options
// and then further examples and usage information.
func Manual(usage string, man string) Helper {
	return func(cmd string, cmds []string) string {
		if isInvalid(cmd) {
			return DefaultHelper(cmd, cmds)
		}
		return fmt.Sprintf("%s\n\n%s\n%s", usage, DefaultHelper(cmd, cmds), man)
	}
}

