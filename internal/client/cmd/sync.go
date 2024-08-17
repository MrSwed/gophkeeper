/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	// syncCmd represents the sync command
	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Synchronize action",
		// Long: ``,
	}

	syncCmd.AddCommand(&cobra.Command{
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

	// rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.generateSaveFlags.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.generateSaveFlags.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
