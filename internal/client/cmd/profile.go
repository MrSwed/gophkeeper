package cmd

import (
	cfg "gophKeeper/internal/client/config"
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func (a *app) addProfileCmd() *app {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profiles menu",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Current profile", cfg.Glob.Get("profile"))
			// todo
			_ = cmd.Usage()
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "list of profiles",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println(cmd.Short)
				prs := cfg.Glob.Get("profiles").(map[string]any)
				cmd.Println("Available profiles: ")
				for profile := range prs {
					if profile == cfg.Glob.Get("profile") {
						cmd.Println(" -", profile, "*")
					} else {
						cmd.Println(" -", profile)
					}
				}
				// cmd.Println(" -", strings.Join(maps.Keys(prs), "\n - "))
				cmd.Println()
				err := cmd.Usage()
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		&cobra.Command{
			Use:   "use",
			Short: "switch to another profile",
			Long:  "if it not exist, it will be created",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println("Current profile", cfg.Glob.Get("profile"))
				cmd.Println("Switching to profile.. ", args[0])
				cfg.Glob.Set("profile", args[0])
				err := cfg.UserLoad()
				if err != nil {
					cmd.PrintErrf("failed to load user profile: %v\n", err)
				}
			},
		},
	)

	a.root.AddCommand(cmd)
	return a
}
