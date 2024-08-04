package cmd

import (
	"os"
	"path/filepath"

	shell "github.com/brianstrauch/cobra-shell"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
func (a *app) addRootCmd() *app {
	c := &cobra.Command{
		Use: func() string {
			_, file := filepath.Split(os.Args[0])

			return file
		}(),
		Short: "GophKeeper client",
		Long:  `Client for save encrypted data`,
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	c.AddCommand(shell.New(c, nil))
	a.root = c
	return a
}
