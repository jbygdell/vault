package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*RotateCommand)(nil)
var _ cli.CommandAutocomplete = (*RotateCommand)(nil)

// RotateCommand is a Command that rotates the encryption key being used
type RotateCommand struct {
	*BaseCommand
}

func (c *RotateCommand) Synopsis() string {
	return "Rotates the underlying encryption key"
}

func (c *RotateCommand) Help() string {
	helpText := `
Usage: vault rotate [options]

  Rotates the underlying encryption key which is used to secure data written
  to the storage backend. This installs a new key in the key ring. This new
  key is used to encrypted new data, while older keys in the ring are used to
  decrypt older data.

  This is an online operation and does not cause downtime. This command is run
  per-cluser (not per-server), since Vault servers in HA mode share the same
  storeage backend.

  Rotate Vault's encryption key:

      $ vault rotate

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *RotateCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *RotateCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *RotateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RotateCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Rotate the key
	err = client.Sys().Rotate()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error rotating key: %s", err))
		return 2
	}

	// Print the key status
	status, err := client.Sys().KeyStatus()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading key status: %s", err))
		return 2
	}

	c.UI.Output("Success! Rotated key")
	c.UI.Output("")
	c.UI.Output(printKeyStatus(status))
	return 0
}
