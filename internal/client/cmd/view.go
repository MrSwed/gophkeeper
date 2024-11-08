/*
This package provides commands for viewing data associated with specific keys.
It uses the Cobra library to define commands for retrieving and displaying
encrypted data.

Main functionalities include:

- Viewing data associated with a specific key.
- Decrypting the data and printing it to standard output.
- Handling errors related to data retrieval and formatting.
*/
package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/spf13/cobra"
)

// addViewCmd adds a command for viewing data associated with a specific key.
// The command takes a single argument, which is the key name.
// It retrieves the data from the server, decrypts it, and prints it to stdout.
// If no key is provided, it displays the command help.
// In case of an error during data retrieval, it prints an appropriate error message.
// If the record does not exist, it indicates that as well.
func (a *app) addViewCmd() *app {
	a.root.AddCommand(&cobra.Command{
		Use:   "view <key name>",
		Short: "View data",
		Long:  `Decrypt data and print it to stdout.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}
			data, err := a.Srv().Get(args[0])
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					cmd.Printf("Record not exist: %s\n", args[0])
				} else {
					cmd.Println("Data get error:", err)
				}
				return
			}
			out, err := json.MarshalIndent(data, "", " ")
			if err != nil {
				cmd.Printf("Data format output error %s %v", err, data)
				return
			}
			cmd.Println(string(out))
		},
	})
	return a
}
