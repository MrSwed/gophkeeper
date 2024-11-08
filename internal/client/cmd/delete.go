/*
Package cmd provides commands for performing delete operations on records.
It uses the Cobra library to define commands for deleting records by their keys.

Main functionalities include:

- Deleting one or more records by specifying their keys.
*/
package cmd

import (
	"database/sql"
	"errors"

	"github.com/spf13/cobra"
)

// addDeleteCmd adds a command for deleting records to the root command.
// The command allows users to delete records by providing their keys as arguments.
// If no keys are provided, an error message is displayed.
// For each key, the command attempts to delete the corresponding record,
// and it reports success or failure for each deletion attempt.
func (a *app) addDeleteCmd() *app {
	cmd := &cobra.Command{
		Use:   "delete [flags] record_key [...record_key]",
		Short: "delete records",
		Long:  `delete records by their keys`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.PrintErrln("You must specify a record key")
				return
			}
			for _, key := range args {
				// if !notConfirm {
				// todo: confirm.  may be github.com/manifoldco/promptui
				// }

				err := a.Srv().Delete(key)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						cmd.Printf("Record not exist: %s\n", key)
					} else {
						cmd.PrintErrf("Delete error: %s\n", err)
					}
				} else {
					cmd.Printf("%s successfully deleted \n", key)
				}
			}
		},
	}
	// cmd.Flags().BoolVarP(&notConfirm, "", "y", false, "do not confirm deleting")

	a.root.AddCommand(cmd)
	return a
}
