package invoker

import "fmt"

// A QuotedCommandError is returned by [Invoker.Invoke] if a command
// required quoting or escaping for invocaton and AllowQuotedCommands
// is not set on the invoker.
type QuotedCommandError struct {
	Command string
}

// Error implements the error interface.
func (e *QuotedCommandError) Error() string {
	return fmt.Sprintf(
		"%q: command required quoting but AllowQuotedCommands not set",
		e.Command,
	)
}
