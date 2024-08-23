package cmd

import (
	"database/sql"
	"errors"

	"github.com/spf13/cobra"
)

func (a *app) addDeleteCmd() *app {
	// var (
	// notConfirm bool
	// )
	cmd := &cobra.Command{
		Use:   "delete [flags] record_key [...record_key]",
		Short: "delete records",
		Long:  `delete records by it keys`,
		// Args:  cobra.ExactArgs(1),
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
					cmd.Printf("%s success deleted \n", key)
				}
			}
		},
	}
	// cmd.Flags().BoolVarP(&notConfirm, "", "y", false, "do not confirm deleting")

	a.root.AddCommand(cmd)
	return a
}
