package cmd

import (
	cfg "gophKeeper/internal/client/config"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (a *app) addConfigCmd() *app {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config action",
		// Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			isAction := false
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if flag.Changed {
					isAction = true
					cfg.Glob.Set(flag.Name, flag.Value)
					cmd.Println("Global configuration: set ", flag.Name, flag.Value)
				}
			})

			if !isAction {
				cmd.Println("Global configuration:")
				cfg.Glob.Print()
				cmd.Println("")

				err := cfg.UserLoad()
				if err != nil {
					cmd.Printf("error load user configuration: %s\n", err)
					return
				}
				cmd.Printf("User \"%s\" configuration:\n", cfg.User.GetString("name"))
				cfg.User.Print()

				cmd.Println()
				err = cmd.Usage()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				if cfg.Glob.Get("autosave") == nil || cfg.Glob.GetBool("autosave") {
					if err := cfg.Glob.Save(); err != nil {
						cmd.Println("Error autosave config", err)
					} else {
						cmd.Println("Success autosave config")
					}
				}
			}
		},
	}
	configCmd.Flags().BoolP("autosave", "a", true, "Global: auto save config")

	updUserCmd := &cobra.Command{
		Use:   "user",
		Short: "change user config params",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.User.Get("name") == nil {
				cmd.Println("You should auth before your can edit your settings")
				return
			}
			cmd.Println("User params")
			cfg.User.Print()

			// todo handle flags
		},
	}
	// updUserCmd.Flags().StringP("mode", "m", "", "remote mode")
	updUserCmd.Flags().StringP("server", "s", "", "server address")
	updUserCmd.Flags().StringP("server_type", "t", "", "server type (default grpc)")
	updUserCmd.Flags().StringP("sync_interval", "i", "", "synchronisation interval ")
	updUserCmd.Flags().BoolP("autosave", "a", true, "User: auto save config")

	configCmd.AddCommand(updUserCmd,
		&cobra.Command{
			Use:   "save",
			Short: "Save now",
			Long:  ``,
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Print("Saving global config.. ")
				if cfg.Glob.IsChanged() {
					err := cfg.Glob.Save()
					if err != nil {
						cmd.Println(err)
					} else {
						cmd.Println("success")
					}
				} else {
					cmd.Println("not changed")
				}

				cmd.Print("Saving user config.. ")
				if cfg.User.GetString("user") != "" {
					if cfg.User.IsChanged() {
						err := cfg.User.Save()
						if err != nil {
							cmd.Println(err)
						} else {
							cmd.Println("succeess")
						}
					} else {
						cmd.Println("not changed")
					}
				} else {
					cmd.Println("You should auth before your can edit your settings")
				}
			},
		},
	)

	a.root.AddCommand(configCmd)
	return a
}
