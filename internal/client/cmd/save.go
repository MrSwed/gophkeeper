/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
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

func saveAuthCmd() (cmd *cobra.Command) {
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

func saveTextCmd() (cmd *cobra.Command) {
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

func saveBinCmd() (cmd *cobra.Command) {
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

func saveCardCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:       "card [flags]",
		Short:     "Save card data",
		Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		ValidArgs: []string{"-num", "-exp-mon", "-exp-year", "-cvv", "-name", "f", "k", "d"},
		Long:      `Encrypts bank cards data`,
		Example:   `  save card --num 2222-4444-5555-1111 --exp 10/29 --cvv 123 --owner "Max Space"`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				data = &card.Model{
					Common: model.Common{},
					Data:   &card.Data{},
				}
				flagLnk = map[string]any{
					"key":         &data.Key,
					"file":        &data.FileName,
					"description": &data.Description,
					"num":         &data.Data.Number,
					"exp":         &data.Data.Exp,
					"cvv":         &data.Data.CVV,
					"owner":       &data.Data.Name,
				}
				err error
			)

			for flag, dataValue := range flagLnk {
				if cmd.Flags().Changed(flag) {
					flagValue, er := cmd.Flags().GetString(flag)
					if er == nil {
						if dv, ok := dataValue.(*string); ok {
							*dv = flagValue
						}
						switch db := dataValue.(type) {
						case *string:
							*db = flagValue
						case model.Settable:
							db.Set(flagValue)
						}
					} else {
						err = errors.Join(err, er)
					}
				}
			}
			err = errors.Join(err, data.Validate())

			if err != nil {
				fmt.Println("Validate: Error: ", err)
			}

			// todo is draft yet
			fmt.Println(data.Data)
			// err := srv.Save(data)
		},
	}
	commonFlags(cmd)

	cmd.Flags().StringP("num", "n", "", "long card number 0000-0000-0000-0000")
	cmd.Flags().StringP("exp", "e", "", "expiry           MM/YY")
	cmd.Flags().StringP("cvv", "c", "", "cvv value        000")
	cmd.Flags().StringP("owner", "o", "", "owner, card holder     Firstname Lastname")

	return
}

func init() {
	var saveCmd = &cobra.Command{
		Use:       "save [command]",
		Short:     "Save data",
		ValidArgs: []string{},
		Args:      cobra.NoArgs,
		Long:      `Encrypts and save data`,
	}

	// saveCmd.AddCommand(shell.New(saveCmd, nil))

	saveCmd.AddCommand(saveAuthCmd(), saveTextCmd(), saveBinCmd(), saveCardCmd())

	rootCmd.AddCommand(saveCmd)

}
