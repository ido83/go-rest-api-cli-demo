package command

import (
	"flag"
	"fmt"

	"go-rest-api-cli/internal/version"
)

// VersionCommand prints the CLI version information.
type VersionCommand struct{}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

func (c *VersionCommand) Name() string        { return "version" }
func (c *VersionCommand) Description() string { return "Show go-rest-api-cli version information" }

func (c *VersionCommand) Run(args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	short := fs.Bool("short", false, "Print only the version string")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *short {
		fmt.Println(version.Version)
		return nil
	}

	fmt.Printf("go-rest-api-cli %s\n", version.Full())
	return nil
}
