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
	"os"
	"path/filepath"

	shell "github.com/brianstrauch/cobra-shell"
	"github.com/spf13/cobra"
)

// addRootCmd adds the root command to the application.
// The root command serves as the main entry point for the GophKeeper client,
// providing information about the application and allowing access to subcommands.
func (a *app) addRootCmd() *app {
	appInfo := `Client for saving encrypted data.`
	if a.v.Version != "" {
		appInfo += ` Version = ` + a.v.Version
	}
	if a.v.Date != "" {
		appInfo += `, build date: ` + a.v.Date
	}
	c := &cobra.Command{
		Use: func() string {
			_, file := filepath.Split(os.Args[0])
			return file
		}(),
		Short: "GophKeeper client",
		Long:  appInfo,
	}
	c.AddCommand(shell.New(c, nil))
	a.root = c
	return a
}
