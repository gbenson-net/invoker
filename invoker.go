// Package invoker abstracts running external commands.
package invoker

import "context"

// Invoker is the interface that wraps the Invoke method.
type Invoker interface {
	// Invoke runs the specified command and returns its combined
	// standard output and standard error.
	Invoke(ctx context.Context, name string, arg ...string) ([]byte, error)
}
