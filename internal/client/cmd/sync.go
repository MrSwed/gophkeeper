/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *app) syncCmd() *cobra.Command {

	// syncCmd represents the sync command
	var cmd = &cobra.Command{
		Use:   "sync",
		Short: "Synchronize action",
		// Long: ``,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "status [command]",
		Short: "Sync status",
		Long:  `show info about last sync`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sync status called")
			// todo
		},
	}, &cobra.Command{
		Use:   "now",
		Short: "Sync now",
		Long:  `synchronize now with server`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sync now called")
			// todo
		},
	})

	a.root.AddCommand(cmd)

	return cmd
}
