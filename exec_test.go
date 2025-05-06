package invoker

import (
	"context"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestExecInvoker(t *testing.T) {
	exec := &ExecInvoker{}
	ctx := t.Context()
	b, err := exec.Invoke(ctx, "/bin/ls", "-1", "/")
	assert.NilError(t, err)
	s := string(b)
	assert.Check(t, strings.Contains(s, "\ndev\n"))
	assert.Check(t, strings.Contains(s, "\ntmp\n"))
	assert.Check(t, strings.Contains(s, "\nusr\n"))
}

func TestExecChecksContext(t *testing.T) {
	exec := &ExecInvoker{}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	b, err := exec.Invoke(ctx, "/bin/ls", "-1", "/")
	assert.ErrorIs(t, err, context.Canceled)
	assert.Check(t, b == nil)
}
