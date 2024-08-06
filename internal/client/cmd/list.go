package cmd

import (
	"gophKeeper/internal/client/model"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func (a *app) addListCmd() *app {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list kept data",
		Long:  `display list of kept data`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("list called")
			// todo
			query := model.ListQuery{}
			dataList, err := a.srv.List(query)
			if err != nil {
				cmd.Printf("Get list error: %s\n", err)
			}
			cmd.Printf("Total: %d\n", dataList.Total)
			for _, item := range dataList.Items {
				cmd.Printf("%s\t%s\t%s\n", item.Key, item.UpdatedAt, item.Description)
			}

		},
	}
	cmd.Flags().StringP("filter", "f", "", "filter list of kept data")
	a.root.AddCommand(cmd)
	return a
}
