package cmd

import (
	cfg "gophKeeper/internal/client/config"

	"github.com/spf13/cobra"
)

// addProfileCmd
// Cobra command for profile operations
func (a *app) addProfileCmd() *app {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profiles menu",
		Run: func(cmd *cobra.Command, args []string) {
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
					cmd.PrintErrf("Error load global config %s from %s", err, cfg.Glob.GetString("config_path"))
					return
				}
				prs := cfg.Glob.GetStringMap("profiles")
				if len(prs) == 0 {
					cmd.Println(`No profiles yet. New default profile will be created, after first save data 
or set some profile data in config.

also you can create new profile by command
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
			Long:  "if it not exist, it will be created",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				err := a.Close()
				if err != nil {
					cmd.PrintErrf("Error close prev profile session  %s  %s", err, cfg.Glob.GetString("config_path"))
					return
				}
				err = cfg.GlobalLoad()
				if err != nil {
					cmd.PrintErrf("Error load global config %s from %s", err, cfg.Glob.GetString("config_path"))
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
				err := cfg.GlobalLoad()
				if err != nil {
					cmd.PrintErrf("Error load global config %s from %s", err, cfg.Glob.GetString("config_path"))
					return
				}
				err = cfg.UserLoad()
				if err != nil {
					cmd.PrintErrf("failed to load user profile: %v\n", err)
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
