package cmd

import (
	"flag"
	"gophKeeper/internal/client/model"

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

	var fs = new(flag.FlagSet)
	err := GenerateFlags(&query, fs)
	if err != nil {
		cmd.Printf("GenerateFlags error: %s\n", err)
	}
	cmd.Flags().AddGoFlagSet(fs)

	a.root.AddCommand(cmd)
	return a
}
