package invoker

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMockInvoker(t *testing.T) {
	mock := NewMock(t)
	mock.ExpectInvoke("hello?").Returns([]byte("world!"), nil)
	ctx := t.Context()
	b, err := mock.Invoke(ctx, "hello?")
	assert.NilError(t, err)
	assert.Equal(t, string(b), "world!")
}

func TestMockChecksContext(t *testing.T) {
	mock := NewMock(t)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	b, err := mock.Invoke(ctx, "hello", "world")
	assert.ErrorIs(t, err, context.Canceled)
	assert.Check(t, b == nil)
}
