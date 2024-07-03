package cmd

import (
	"fmt"
	cfg "gophKeeper/internal/client/config"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {

	// configCmd represents the cfg command
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Config action",
		// Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			isAction := false
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if flag.Changed {
					isAction = true
					cfg.Glob.Set(flag.Name, flag.Value)
					fmt.Println("Global configuration: set ", flag.Name, flag.Value)
				}
			})

			if !isAction {
				fmt.Println("Global configuration:")
				cfg.Glob.Print()
				if cfg.User.GetString("name") != "" {
					fmt.Printf("User %s configuration:\n", cfg.User.GetString("name"))
					cfg.User.Print()
				}

				fmt.Println()
				err := cmd.Usage()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				if cfg.Glob.Get("autosave") == nil || cfg.Glob.GetBool("autosave") {
					if err := cfg.Glob.Save(); err != nil {
						fmt.Println("Error autosave config", err)
					} else {
						fmt.Println("Success autosave config")
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
			if cfg.User.Get("user") == nil {
				fmt.Println("You should auth before your can edit your settings")
				return
			}
			fmt.Println("User params", cfg.User)
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
				fmt.Print("Saving global config.. ")
				if cfg.Glob.IsChanged() {
					err := cfg.Glob.Save()
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("success")
					}
				} else {
					fmt.Println("not changed")
				}

				fmt.Print("Saving user config.. ")
				if cfg.User.GetString("user") != "" {
					if cfg.User.IsChanged() {
						err := cfg.User.Save()
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("succeess")
						}
					} else {
						fmt.Println("not changed")
					}
				} else {
					fmt.Println("You should auth before your can edit your settings")
				}
			},
		},
	)

	rootCmd.AddCommand(configCmd)
}
