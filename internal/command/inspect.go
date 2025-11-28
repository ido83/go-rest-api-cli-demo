package command

import (
	"flag"
	"fmt"
	"sort"

	cfgstore "go-rest-api-cli-demo/internal/config"
)

// InspectCommand inspects stored profiles.
type InspectCommand struct{}

func NewInspectCommand() *InspectCommand {
	return &InspectCommand{}
}

func (i *InspectCommand) Name() string        { return "inspect" }
func (i *InspectCommand) Description() string { return "Inspect stored profiles" }

func (i *InspectCommand) Run(args []string) error {
	if len(args) == 0 || args[0] == "profiles" {
		return i.inspectProfiles()
	}

	switch args[0] {
	case "profile":
		return i.inspectProfile(args[1:])
	default:
		i.printUsage()
		return fmt.Errorf("unknown inspect target: %s", args[0])
	}
}

func (i *InspectCommand) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go-rest-api-cli-demo inspect profiles")
	fmt.Println("  go-rest-api-cli-demo inspect profile --name NAME")
}

func (i *InspectCommand) inspectProfiles() error {
	cfg, err := cfgstore.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles defined.")
		return nil
	}

	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Profiles:")
	for _, name := range names {
		pf := cfg.Profiles[name]
		authInfo := pf.AuthType
		if authInfo == "" {
			authInfo = "none"
		}
		fmt.Printf("- %s\n", name)
		fmt.Printf("  Base URL : %s\n", pf.BaseURL)
		fmt.Printf("  Auth     : %s\n", authInfo)
		if len(pf.Headers) > 0 {
			fmt.Println("  Headers  :")
			for k, v := range pf.Headers {
				fmt.Printf("    %s: %s\n", k, v)
			}
		}
	}
	return nil
}

func (i *InspectCommand) inspectProfile(args []string) error {
	fs := flag.NewFlagSet("inspect profile", flag.ContinueOnError)
	name := fs.String("name", "", "Profile name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("--name is required")
	}

	cfg, err := cfgstore.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pf, ok := cfg.Profiles[*name]
	if !ok {
		return fmt.Errorf("profile %q not found", *name)
	}

	authInfo := pf.AuthType
	if authInfo == "" {
		authInfo = "none"
	}

	fmt.Printf("Profile %q\n", *name)
	fmt.Printf("  Base URL : %s\n", pf.BaseURL)
	fmt.Printf("  Auth     : %s\n", authInfo)
	if pf.User != "" {
		fmt.Printf("  User     : %s\n", pf.User)
	}
	if pf.Token != "" {
		fmt.Printf("  Token    : (set)\n")
	}
	if len(pf.Headers) > 0 {
		fmt.Println("  Headers  :")
		for k, v := range pf.Headers {
			fmt.Printf("    %s: %s\n", k, v)
		}
	}
	return nil
}
