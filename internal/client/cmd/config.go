package cmd

import (
	"encoding/json"
	cfg "gophKeeper/internal/client/config"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// addConfigCmd
// Cobra commands for config operations
func (a *app) addConfigCmd() *app {

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config action",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.GlobalLoad()
			if err != nil {
				cmd.PrintErrf("Error load global config %s from %s", err, cfg.Glob.GetString("config_path"))
				return
			}
			isAction := false
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if flag.Changed {
					isAction = true
					cfg.Glob.Set(flag.Name, flag.Value)
					cmd.Printf("Global configuration: set `%s` = `%s`\n", flag.Name, flag.Value)
				}
			})

			if !isAction {
				cmd.Println("Global configuration:")
				out, err := json.MarshalIndent(cfg.Glob.AllSettings(), "", " ")
				if err != nil {
					cmd.Printf("Data format out err %s %v", err, cfg.Glob.AllSettings())
					return
				}
				cmd.Println(string(out))
				cmd.Println()
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
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
			}

			isAction := false
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if flag.Changed {
					isAction = true
					cfg.User.Set(flag.Name, flag.Value)
					cmd.Printf("User configuration: set `%s` = `%s`\n", flag.Name, flag.Value)
				}
			})

			if !isAction {
				cmd.Println("User configuration:")
				out, err := json.MarshalIndent(cfg.User.AllSettings(), "", " ")
				if err != nil {
					cmd.Printf("Data format out err %s %v", err, cfg.User.AllSettings())
					return
				}
				cmd.Println(string(out))
				cmd.Println()
			} else {
				if cfg.User.Get("autosave") == nil || cfg.User.GetBool("autosave") {
					if err := cfg.User.Save(); err != nil {
						cmd.Println("Error autosave config", err)
					} else {
						cmd.Println("Success autosave config")
					}
				}
			}
		},
	}
	// updUserCmd.Flags().StringP("mode", "m", "", "remote mode")
	// updUserCmd.Flags().StringP("server_type", "t", "", "server type (default grpc)")
	updUserCmd.Flags().StringP("server", "s", "", "server address")
	updUserCmd.Flags().DurationP("sync.timeout.sync", "t", 0, "synchronisation timeout")
	updUserCmd.Flags().DurationP("sync.timeout.register", "", 0, "register at server timeout")
	updUserCmd.Flags().StringP("email", "e", "", "User email")
	updUserCmd.Flags().BoolP("autosave", "a", true, "Auto save user config")

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Save now, for shell mode, if autosave disabled",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
			}

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
			if cfg.User.IsChanged() {
				err := cfg.User.Save()
				if err != nil {
					cmd.Println(err)
				} else {
					cmd.Println("success")
				}
			} else {
				cmd.Println("not changed")
			}
		},
	}

	configCmd.AddCommand(updUserCmd, saveCmd)

	a.root.AddCommand(configCmd)
	return a
}
