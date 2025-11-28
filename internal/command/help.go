package command

import "fmt"

// HelpCommand prints available commands and an example.
type HelpCommand struct {
	reg     *Registry
	appName string
}

func NewHelpCommand(reg *Registry, appName string) *HelpCommand {
	return &HelpCommand{reg: reg, appName: appName}
}

func (h *HelpCommand) Name() string        { return "help" }
func (h *HelpCommand) Description() string { return "Show help" }

func (h *HelpCommand) Run(args []string) error {
	fmt.Printf("%s - simple REST API CLI\n\n", h.appName)
	fmt.Println("Usage:")
	fmt.Printf("  %s <command> [flags]\n\n", h.appName)

	fmt.Println("Commands:")
	for _, c := range h.reg.All() {
		fmt.Printf("  %-8s %s\n", c.Name(), c.Description())
	}
	fmt.Println()
	fmt.Println("Example:")
	fmt.Printf("  %s call --method GET --url \"https://api.agify.io/?name=meelad\"\n", h.appName)

	return nil
}
