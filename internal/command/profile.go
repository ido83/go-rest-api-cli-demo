package command

import (
	"flag"
	"fmt"

	cfgstore "go-rest-api-cli-demo/internal/config"
)

// ProfileCommand manages profiles (add/list/remove).
type ProfileCommand struct{}

func NewProfileCommand() *ProfileCommand {
	return &ProfileCommand{}
}

func (p *ProfileCommand) Name() string        { return "profile" }
func (p *ProfileCommand) Description() string { return "Manage profiles (add/list/remove)" }

func (p *ProfileCommand) Run(args []string) error {
	if len(args) == 0 {
		p.printUsage()
		return nil
	}

	switch args[0] {
	case "add":
		return p.runAdd(args[1:])
	case "list":
		return p.runList()
	case "remove":
		return p.runRemove(args[1:])
	default:
		p.printUsage()
		return fmt.Errorf("unknown profile action: %s", args[0])
	}
}

func (p *ProfileCommand) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("go-rest-api-cli-demo profile add --name NAME --base-url URL [--auth ...] [--header ...]")
	fmt.Println("  go-rest-api-cli-demo profile list")
	fmt.Println("  go-rest-api-cli-demo profile remove --name NAME")
}

func (p *ProfileCommand) runAdd(args []string) error {
	fs := flag.NewFlagSet("profile add", flag.ContinueOnError)

	name := fs.String("name", "", "Profile name (required)")
	baseURL := fs.String("base-url", "", "Base URL, e.g. https://api.example.com")
	authType := fs.String("auth", "none", "Auth type: none|basic|bearer")
	user := fs.String("user", "", "Username for basic auth")
	pass := fs.String("pass", "", "Password for basic auth")
	token := fs.String("token", "", "Bearer token for auth")

	headers := HeaderFlag{}
	fs.Var(&headers, "header", "Default header 'Key: Value' (can be repeated)")

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

	pf := cfgstore.Profile{
		Name:    *name,
		BaseURL: *baseURL,
		Headers: map[string]string(headers),

		AuthType: *authType,
		User:     *user,
		Pass:     *pass,
		Token:    *token,
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]cfgstore.Profile)
	}
	cfg.Profiles[pf.Name] = pf

	if err := cfgstore.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Printf("Profile %q saved\n", pf.Name)
	return nil
}

func (p *ProfileCommand) runList() error {
	cfg, err := cfgstore.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles defined.")
		return nil
	}

	fmt.Println("Profiles:")
	for name, pf := range cfg.Profiles {
		authInfo := pf.AuthType
		if authInfo == "" {
			authInfo = "none"
		}
		fmt.Printf("- %s (base-url: %s, auth: %s)\n", name, pf.BaseURL, authInfo)
	}
	return nil
}

func (p *ProfileCommand) runRemove(args []string) error {
	fs := flag.NewFlagSet("profile remove", flag.ContinueOnError)
	name := fs.String("name", "", "Profile name to remove")

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

	if _, ok := cfg.Profiles[*name]; !ok {
		return fmt.Errorf("profile %q not found", *name)
	}

	delete(cfg.Profiles, *name)

	if err := cfgstore.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Printf("Profile %q removed\n", *name)
	return nil
}
