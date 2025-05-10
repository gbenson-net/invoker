package invoker

import (
	"context"
	"slices"
	"strings"

	"github.com/kballard/go-shellquote"
	"golang.org/x/crypto/ssh"
)

// SSHInvoker invokes commands using an [ssh.Client].  Note that the
// SSH protocol requires that the command and arguments passed to
// [Invoker.Invoke] be joined into a quoted, space-delimited string
// for transmission.  SSHInvoker does this using [shellquote.Join],
// which assumes /bin/sh-style quoting, but how the receiver will
// interpret this quoted string cannot be known, so SSHInvoker will
// reject any command requiring quoting unless [AllowQuotedCommands]
// has been explicitly set to true.
type SSHInvoker struct {
	// Client is the SSH client connection to use.
	Client *ssh.Client

	// AllowQuotedCommands permits using commands requiring quoting or
	// escaping for transmission.  If false, Invoke will return an
	// error rather than transmit a command the receiver could
	// misinterpret.
	AllowQuotedCommands bool
}

// NewSSH returns an invoker that invokes commands using the supplied
// [ssh.Client].
func NewSSH(c *ssh.Client) *SSHInvoker {
	return &SSHInvoker{Client: c}
}

// Invoke implements the Invoker interface.
func (si *SSHInvoker) Invoke(
	ctx context.Context,
	name string,
	arg ...string,
) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	argv := append([]string{name}, arg...)
	cmd := shellquote.Join(argv...)
	if !si.AllowQuotedCommands {
		check := strings.Fields(cmd)
		if !slices.Equal(check, argv) {
			return nil, &QuotedCommandError{cmd}
		}
	}

	session, err := si.Client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return session.CombinedOutput(cmd)
}
