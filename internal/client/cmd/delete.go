package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd())
}

// listCmd represents the list command
func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [flags] record_key [...record_key]",
		Short: "delete record",
		Long:  `delete record by it key (ID)`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
			// todo
		},
	}
	cmd.Flags().StringP("quite", "q", "", "do not show log")

	return cmd
}
