/*
Package cmd provides commands for managing user profiles in the application.
It uses the Cobra library to define commands and subcommands for profile operations.

Main functionalities include:

- Viewing the current profile.
- Listing available profiles.
- Switching to another profile.
- Changing the password for the current profile.
*/
package cmd

import (
	cfg "gophKeeper/internal/client/config"

	"github.com/spf13/cobra"
)

// addProfileCmd adds commands for profile operations to the root command.
// The main command provides a menu for profile management.
func (a *app) addProfileCmd() *app {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profiles menu",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.GlobalLoad()
			if err != nil {
				cmd.PrintErrf("Error loading global config %s from %s", err, cfg.Glob.GetString("config_path"))
				return
			}
			cmd.Println("Current profile", cfg.GetUserName())
			// todo
			_ = cmd.Usage()
		},
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "list of profiles",
			Run: func(cmd *cobra.Command, args []string) {
				err := cfg.GlobalLoad()
				if err != nil {
					cmd.PrintErrf("Error loading global config %s from %s", err, cfg.Glob.GetString("config_path"))
					return
				}
				prs := cfg.Glob.GetStringMap("profiles")
				if len(prs) == 0 {
					cmd.Println(`No profiles yet. A new default profile will be created after the first save of data 
or setting some profile data in the config.
You can also create a new profile by the command
    profile use <new_name>`)
				}
				cmd.Println("Available profiles: ")
				for name, profile := range prs {
					p, ok := profile.(map[string]any)
					if !ok {
						continue
					}
					if n, ok := p["name"]; ok {
						name = n.(string)
					}
					if name == cfg.GetUserName() {
						cmd.Println(" -", name, "*")
					} else {
						cmd.Println(" -", name)
					}
				}
				cmd.Println()
			},
		},
		&cobra.Command{
			Use:   "use",
			Short: "switch to another profile",
			Long:  "If it does not exist, it will be created.",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				err := a.Close()
				if err != nil {
					cmd.PrintErrf("Error closing previous profile session  %s  %s", err, cfg.Glob.GetString("config_path"))
					return
				}
				err = cfg.GlobalLoad()
				if err != nil {
					cmd.PrintErrf("Error loading global config %s from %s", err, cfg.Glob.GetString("config_path"))
					return
				}
				cmd.Println("Current profile", cfg.GetUserName())
				cmd.Println("Switching to profile.. ", args[0])
				cfg.Glob.Set("profile", args[0])
				err = cfg.UserLoad(true)
				if err != nil {
					cmd.PrintErrf("failed to load user profile: %v\n", err)
				}
			},
		},
		&cobra.Command{
			Use:   "password",
			Short: "set new password",
			Run: func(cmd *cobra.Command, args []string) {
				err := cfg.UserLoad()
				if err != nil {
					cmd.PrintErrf("failed to load config: %v\n", err)
				}
				cmd.Println("Current profile", cfg.GetUserName())
				cmd.Println("Change password:.. ")
				err = a.Srv().ChangePasswd()
				if err != nil {
					cmd.PrintErrf("failed to change password: %v\n", err)
				}
			},
		},
	)
	a.root.AddCommand(cmd)
	return a
}
