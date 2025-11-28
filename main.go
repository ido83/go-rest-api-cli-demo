package main

import (
	"fmt"
	"go-rest-api-cli-demo/internal/command"
	"go-rest-api-cli-demo/internal/httpclient"
	"os"
)

func main() {
	reg := command.NewRegistry()

	factory := httpclient.Factory{}
	reg.Register(command.NewCallCommand(factory))
	reg.Register(command.NewProfileCommand())
	reg.Register(command.NewInspectCommand())
	reg.Register(command.NewHelpCommand(reg, "go-rest-api-cli-demo"))

	if len(os.Args) < 2 {
		if helpCmd, ok := reg.Get("help"); ok {
			_ = helpCmd.Run(nil)
		}
		os.Exit(1)
	}

	name := os.Args[1]
	cmd, ok := reg.Get(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", name)
		if helpCmd, ok := reg.Get("help"); ok {
			_ = helpCmd.Run(nil)
		}
		os.Exit(1)
	}

	if err := cmd.Run(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
