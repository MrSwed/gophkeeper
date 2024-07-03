package cmd

import (
	"fmt"
	cfg "gophKeeper/internal/client/config"
	"os"
	"path/filepath"

	shell "github.com/brianstrauch/cobra-shell"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: func() string {
		_, file := filepath.Split(os.Args[0])

		return file
	}(),
	Short: "GophKeeper client",
	Long:  `Client for save encrypted data`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if cfg.Glob.GetBool("autosave") && cfg.Glob.Get("changed_at") != nil {
		fmt.Print("Saving global cfg files at exit..")
		err = cfg.Glob.Save()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(" ..Success")
	}
	if cfg.User.Get("name") != nil && cfg.User.GetBool("autosave") && cfg.User.Get("changed_at") != nil {
		fmt.Print("Saving user cfg files at exit..")
		err = cfg.User.Save()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(" ..Success")
	}
}

func init() {

	rootCmd.AddCommand(shell.New(rootCmd, nil))

}
