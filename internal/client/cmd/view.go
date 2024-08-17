package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

func (a *app) addViewCmd() *app {
	a.root.AddCommand(&cobra.Command{
		Use:   "view <key name>",
		Short: "View data",
		Long:  `Decrypt data and print it to stdout.`,
		// Example: `  `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}

			data, err := a.srv.Get(args[0])

			if err != nil {
				cmd.Println("Data get error", err)
				return
			}
			out, err := json.MarshalIndent(data, "", " ")
			if err != nil {
				cmd.Printf("Data format out err %s %v", err, data)
				return
			}
			cmd.Println(string(out))
		},
	})
	return a
}
