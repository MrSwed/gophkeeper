package cmd

import (
	"fmt"
	"gophKeeper/internal/client/config"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func init() {
	rootCmd.AddCommand(profileCmd())
}

// listCmd represents the list command
func profileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profiles menu",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Current profile", config.Glob.Get("profile"))
			// todo
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "list of profiles",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(cmd.Short)
				prs := config.Glob.Get("profiles").(map[string]any)
				fmt.Println("Available profiles: ")
				fmt.Println(" -", strings.Join(maps.Keys(prs), "\n -"))
				fmt.Println()
				err := cmd.Usage()
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		&cobra.Command{
			Use:   "use",
			Short: "switch to another profile",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Current profile", config.Glob.Get("profile.name"))
				// todo
			},
		},
	)
	return cmd
}
