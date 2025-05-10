package invoker

import (
	"context"
	"encoding/json"
	"fmt"

	"gotest.tools/v3/assert"
)

// TestingT is the subset of testing.T used by the invoker package.
type TestingT interface {
	assert.TestingT
	Helper()
}

// MockInvoker mocks Invoker.
type MockInvoker struct {
	t            TestingT
	expectations []*mockInvocation
}

// NewMock creates a MockInvoker for testing.
func NewMock(t TestingT) *MockInvoker {
	return &MockInvoker{t: t}
}

// A MockInvocation represents an expected invocation.
type MockInvocation interface {
	// Returns specifies the result this mocked invocation will return.
	Returns([]byte, error)
}

// ExpectInvoke adds a new expected invocation.
func (mi *MockInvoker) ExpectInvoke(name string, arg ...string) MockInvocation {
	e := &mockInvocation{name: name, args: arg}
	mi.expectations = append(mi.expectations, e)
	return e
}

// mockInvocation implements MockInvocation.
type mockInvocation struct {
	name   string
	args   []string
	result string // immutable
	err    error
}

// Returns implements MockInvocation.
func (e *mockInvocation) Returns(result []byte, err error) {
	if e.result != "" || e.err != nil {
		panic("expectation already set")
	}

	e.result = string(result)
	e.err = err
}

// MarshalJSON() implements json.Marshaler.
func (e *mockInvocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(append([]string{e.name}, e.args...))
}

// Invoke implements Invoker.
func (mi *MockInvoker) Invoke(
	ctx context.Context,
	name string,
	arg ...string,
) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if mi.t == nil {
		panic("nil testing.T")
	}
	mi.t.Helper()

	assert.Assert(mi.t, len(mi.expectations) > 0)
	expect := mi.expectations[0]
	mi.expectations = mi.expectations[1:]

	assert.Equal(mi.t, name, expect.name)
	assert.DeepEqual(mi.t, arg, expect.args)

	return []byte(expect.result), expect.err
}

// ExpectationsWereMet returns an error unless all expectations were met.
func (mi *MockInvoker) ExpectationsWereMet() error {
	if len(mi.expectations) == 0 {
		return nil
	}

	b, err := json.Marshal(mi.expectations)
	assert.NilError(mi.t, err)
	return fmt.Errorf("unmet expectations: %v", string(b))
}
