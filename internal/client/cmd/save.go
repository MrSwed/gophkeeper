/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/card"

	"github.com/spf13/cobra"
)

func commonFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("file", "f", "", "read from file")
	cmd.Flags().StringP("key", "k", "", "set your entry key-identifier")
	cmd.Flags().StringP("description", "d", "", "description, will be displayed in the list of entries list")
}

func (a *app) saveAuthCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:       "auth",
		Short:     "Save auth data",
		Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		ValidArgs: []string{"l", "s", "k", "f", "d"},
		Long:      `Encrypts login/password pairs.`,
		Example: `  save auth -l login -s password
  save auth -l login -s password -d site.com`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("save auth called")
			// todo here
		},
	}
	commonFlags(cmd)
	cmd.Flags().StringP("login", "l", "", "login")
	cmd.Flags().StringP("password", "p", "", "password")
	return
}

func (a *app) saveTextCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:       "text [flags]",
		Short:     "Save text data",
		Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		ValidArgs: []string{"s", "k", "f", "d"},
		Long:      `Encrypts text data`,
		Example: `  save text -f filename
  save text -k switch -d description -s
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("save text called")
			// todo here
		},
	}
	commonFlags(cmd)
	cmd.Flags().StringP("test", "t", "", "text data")
	return
}

func (a *app) saveBinCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:       "bin [flags]",
		Short:     "Save binary data",
		Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		ValidArgs: []string{"k", "f", "d"},
		Long:      `Encrypts bin data`,
		Example:   `   save bin -f filename`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("save bin called")
			// todo here
		},
	}
	commonFlags(cmd)
	return
}

func (a *app) saveCardCmd() (cmd *cobra.Command) {
	data := &card.Model{
		Common: model.Common{},
		Data:   &card.Data{},
	}

	cmd = &cobra.Command{
		Use:       "card [flags]",
		Short:     "Save card data",
		Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		ValidArgs: []string{"-num", "-exp-mon", "-exp-year", "-cvv", "-name", "f", "k", "d"},
		Long:      `Encrypts bank cards data`,
		Example:   `  save card --num 2222-4444-5555-1111 --exp 10/29 --cvv 123 --owner "Max Space"`,
		Run: func(cmd *cobra.Command, args []string) {

			data.Data.Sanitize()
			err := data.Validate()

			if err != nil {
				fmt.Println("Validate: Error: ", err)
			}

			// todo is draft yet
			fmt.Println(data.Data)

			err = a.srv.Save(data)
			if err != nil {
				fmt.Println(err.Error())

			}
		},
	}

	err := GenerateFlags(&data.Common, cmd.Flags())
	if err != nil {
		cmd.Printf("GenerateFlags error: %s\n", err)
	}
	err = GenerateFlags(data.Data, cmd.Flags())
	if err != nil {
		cmd.Printf("GenerateFlags error: %s\n", err)
	}

	return
}

func (a *app) addSaveCmd() *app {
	var saveCmd = &cobra.Command{
		Use:       "save [command]",
		Short:     "Save data",
		ValidArgs: []string{},
		Args:      cobra.NoArgs,
		Long:      `Encrypts and save data`,
	}

	saveCmd.AddCommand(a.saveAuthCmd(), a.saveTextCmd(), a.saveBinCmd(), a.saveCardCmd())

	a.root.AddCommand(saveCmd)
	return a
}
