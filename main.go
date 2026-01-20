package main

import (
	"fmt"
	"os"

	"go-rest-api-cli/internal/command"
	"go-rest-api-cli/internal/httpclient"
	"go-rest-api-cli/internal/version"
)

const appName = "go-rest-api-cli"

const asciiLogo = `
#     _______  _______     _______  _______  _______ _________ _______  _______ _________    _______  _       _________
#    (  ____ \(  ___  )   (  ____ )(  ____ \(  ____ \\__   __/(  ___  )(  ____ )\__   __/   (  ____ \( \      \__   __/
#    | (    \/| (   ) |   | (    )|| (    \/| (    \/   ) (   | (   ) || (    )|   ) (      | (    \/| (         ) (   
#    | |      | |   | |   | (____)|| (__    | (_____    | |   | (___) || (____)|   | |      | |      | |         | |   
#    | | ____ | |   | |   |     __)|  __)   (_____  )   | |   |  ___  ||  _____)   | |      | |      | |         | |   
#    | | \_  )| |   | |   | (\ (   | (            ) |   | |   | (   ) || (         | |      | |      | |         | |   
#    | (___) || (___) | _ | ) \ \__| (____/\/\____) |   | | _ | )   ( || )      ___) (___ _ | (____/\| (____/\___) (___
#    (_______)(_______)(_)|/   \__/(_______/\_______)   )_((_)|/     \||/       \_______/(_)(_______/(_______/\_______/
#                                                                                                                                                               
				  200000083            380000005                                   
                                200000000000003      180000000000002                                
                               600002     180008    800001     300009                               
                              40000         700008800007         80004                              
                              00007           20000005            0000                              
                              0000             200002             0000                              
                              00001           90000009           70000                              
                              20000         500002200004         00002                              
                               5000097   7400004    4000047   7600005                               
                                 9000000000006        6000000000008                                 
                                    4000005              5000006                             
																		
                                                                                                    
`

// printBanner prints the ASCII logo and version info.
// Note: use Println without embedded "\n" to keep 'go vet' happy.
func printBanner() {
	fmt.Print(asciiLogo) // use Print to avoid adding an extra newline
	fmt.Println(appName)
	fmt.Printf("Version: %s | Commit: %s | Built: %s\n\n", version.Version, version.Commit, version.Date)
}

func main() {
	factory := httpclient.Factory{}

	registry := command.NewRegistry()
	registry.Register(command.NewCallCommand(factory))
	registry.Register(command.NewProfileCommand())
	registry.Register(command.NewInspectCommand())

	// Create help command with app name and register it
	helpCmd := command.NewHelpCommand(registry, appName)
	registry.Register(helpCmd)

	// Version command
	registry.Register(command.NewVersionCommand())

	// No command specified → show banner + global help
	if len(os.Args) < 2 {
		printBanner()
		_ = helpCmd.Run([]string{})
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd, ok := registry.Get(cmdName)
	if !ok || cmd == nil {
		// Unknown command → show banner + global help
		printBanner()
		fmt.Fprintf(os.Stderr, "unknown command %q (try: %s help)\n\n", cmdName, appName)
		_ = helpCmd.Run([]string{})
		os.Exit(1)
	}

	if err := cmd.Run(cmdArgs); err != nil {
		// Any error → show error + banner + relevant help
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		printBanner()

		// Try help for this specific command, otherwise global help
		if cmdName != "help" {
			_ = helpCmd.Run([]string{cmdName})
		} else {
			_ = helpCmd.Run([]string{})
		}
		os.Exit(1)
	}
}
