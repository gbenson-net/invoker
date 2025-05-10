package invoker

import (
	"context"
	"os/exec"
)

// ExecInvoker invokes commands using [exec.Command].
type ExecInvoker struct{}

// Exec invokes commands using [exec.Command].
var Exec = Invoker(&defaultExecInvoker)

var defaultExecInvoker ExecInvoker

// Invoke implements the Invoker interface.
func (e *ExecInvoker) Invoke(
	ctx context.Context,
	name string,
	arg ...string,
) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return exec.Command(name, arg...).CombinedOutput()
}
