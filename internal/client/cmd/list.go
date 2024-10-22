/*
This package provides commands for listing data kept in the application.
It uses the Cobra library to define commands for retrieving and displaying
a list of stored data.

Main functionalities include:

- Displaying a list of kept data with details such as key, date, and description.
*/
package cmd

import (
	"time"

	"gophKeeper/internal/client/model"
	"gophKeeper/internal/helper"

	"github.com/spf13/cobra"
)

// addListCmd adds a command for listing kept data to the root command.
// The command retrieves and displays a list of stored data, including
// the total number of items and details for each item such as key, date,
// and description.
func (a *app) addListCmd() *app {
	query := model.ListQuery{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list kept data",
		Long:  `display list of kept data`,
		Run: func(cmd *cobra.Command, args []string) {
			dataList, err := a.Srv().List(query)
			if err != nil {
				cmd.Printf("Get list error: %s\n", err)
			}
			cmd.Printf("Total: %d\n", dataList.Total)
			for _, item := range dataList.Items {
				date := item.UpdatedAt
				if date == nil {
					date = &item.CreatedAt
				}
				cmd.Printf("%s\t%s\t%s\n", item.Key, date.Format(time.DateTime), item.Description)
			}
		},
	}
	err := helper.GenerateFlags(&query, cmd.Flags())
	if err != nil {
		cmd.Printf("GenerateFlags error: %s\n", err)
	}
	a.root.AddCommand(cmd)
	return a
}
