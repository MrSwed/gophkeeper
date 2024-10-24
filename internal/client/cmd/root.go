/*
This package provides the root command for the GophKeeper client application.
It uses the Cobra library to define the main command for the application,
which serves as the entry point for various subcommands related to managing
encrypted data.

Main functionalities include:

- Displaying application information, including version and build date.
- Setting up the root command and adding shell command support.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	shell "github.com/brianstrauch/cobra-shell"
	"github.com/spf13/cobra"
)

// shortVersionInfo returns a short version information string.
func (a *app) shortVersionInfo() string {
	return fmt.Sprintf("Version: %s", a.v.Version)
}

// fullVersionInfo returns a full version information string.
func (a *app) fullVersionInfo() string {
	return fmt.Sprintf("Version: %s\nBuild date: %s\nCommit: %s", a.v.Version, a.v.Date, a.v.Commit)
}

// addRootCmd adds the root command to the application.
// The root command serves as the main entry point for the GophKeeper client,
// providing information about the application and allowing access to subcommands.
func (a *app) addRootCmd() *app {
	c := &cobra.Command{
		Use: func() string {
			_, file := filepath.Split(os.Args[0])
			return file
		}(),
		Short: "GophKeeper client",
		Long:  "Client for saving encrypted data.",
	}

	// Add a flag for displaying the full version information
	c.Flags().BoolP("version", "v", false, "Display full version information")
	c.Run = func(cmd *cobra.Command, args []string) {
		if showVersion, _ := cmd.Flags().GetBool("version"); showVersion {
			fmt.Println(a.fullVersionInfo())
		} else {
			fmt.Println(a.shortVersionInfo())
			fmt.Println()
			_ = cmd.Usage()
		}
	}

	c.AddCommand(shell.New(c, nil))
	a.root = c
	return a
}
