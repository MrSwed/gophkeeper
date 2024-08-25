package cmd

import (
	"github.com/spf13/cobra"
)

/*
	syncCmd

todo:

	0.1. sunc now
	0.2. sync cron(schedule)
	0.4. agent mode, runned as daemon (?)
*/
func (a *app) syncCmd() *cobra.Command {

	// syncCmd represents the sync command
	var cmd = &cobra.Command{
		Use:   "sync",
		Short: "Synchronize action",
		// Long: ``,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "status [command]",
		Short: "Sync status",
		Long:  `show info about last sync`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("sync status called")
			// todo
		},
	}, &cobra.Command{
		Use:   "now",
		Short: "Sync now",
		Long:  `synchronize now with server`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("sync now called")
			// todo
		},
	})

	a.root.AddCommand(cmd)

	return cmd
}
