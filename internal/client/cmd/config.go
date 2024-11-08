/*
This package provides commands for configuration operations in the application.
It uses the Cobra library to define commands and subcommands for managing global
and user configurations.

Main functionalities include:

- Viewing and modifying global configuration.
- Viewing and modifying user configuration.
- Saving configurations manually.
*/
package cmd

import (
	"encoding/json"
	cfg "gophKeeper/internal/client/config"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// addConfigCmd adds commands for configuration operations to the root command.
// The main command provides actions for managing global and user configurations.
func (a *app) addConfigCmd() *app {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config actions",
	}

	globalCmd := &cobra.Command{
		Use:   "global",
		Short: "Global config actions",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.GlobalLoad()
			if err != nil {
				cmd.PrintErrf("Error loading global config %s from %s", err, cfg.Glob.GetString("config_path"))
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
					cmd.Printf("Data format output error: %s %v", err, cfg.Glob.AllSettings())
					return
				}
				cmd.Println(string(out))
				cmd.Println()
			} else if cfg.Glob.Get("autosave") == nil || cfg.Glob.GetBool("autosave") {
				if err := cfg.Glob.Save(); err != nil {
					cmd.Println("Error autosaving config", err)
				} else {
					cmd.Println("Success autosaving config")
				}
			}
		},
	}

	globalCmd.Flags().BoolP("autosave", "a", true, "Global: auto save config")
	globalCmd.Flags().BoolP("debug", "d", false, "Global: debug mode")

	updUserCmd := &cobra.Command{
		Use:   "user",
		Short: "change user config parameters",
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
					cmd.Printf("Data format output error: %s %v", err, cfg.User.AllSettings())
					return
				}
				cmd.Println(string(out))
				cmd.Println()
			} else {
				if cfg.User.Get("autosave") == nil || cfg.User.GetBool("autosave") {
					if err := cfg.User.Save(); err != nil {
						cmd.Println("Error autosaving config", err)
					} else {
						cmd.Println("Success autosaving config")
					}
				}
			}
		},
	}

	updUserCmd.Flags().StringP("server", "s", "", "server address")
	updUserCmd.Flags().DurationP("sync.timeout.sync", "t", 0, "synchronization timeout")
	updUserCmd.Flags().DurationP("sync.timeout.register", "", 0, "register at server timeout")
	updUserCmd.Flags().StringP("email", "e", "", "User email")
	updUserCmd.Flags().BoolP("autosave", "a", true, "Auto save user config")

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Save now, for shell mode, if autosave is disabled",
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

	configCmd.AddCommand(globalCmd, updUserCmd, saveCmd)
	a.root.AddCommand(configCmd)
	return a
}
