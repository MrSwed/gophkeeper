package cmd

import (
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/helper"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func (a *app) addListCmd() *app {
	query := model.ListQuery{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list kept data",
		Long:  `display list of kept data`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("list called")
			// todo

			dataList, err := a.Srv().List(query)
			if err != nil {
				cmd.Printf("Get list error: %s\n", err)
			}
			cmd.Printf("Total: %d\n", dataList.Total)
			for _, item := range dataList.Items {
				cmd.Printf("%s\t%s\t%s\n", item.Key, item.UpdatedAt, item.Description)
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
