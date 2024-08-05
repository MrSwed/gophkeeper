package cmd

import (
	"fmt"
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
			fmt.Println("list called")
			// todo
			query := model.ListQuery{}
			dataList, err := a.srv.List(query)
			if err != nil {
				fmt.Printf("Get list error: %s\n", err)
			}
			fmt.Printf("List view: %v\n", dataList)
		},
	}
	cmd.Flags().StringP("filter", "f", "", "filter list of kept data")
	a.root.AddCommand(cmd)
	return a
}
