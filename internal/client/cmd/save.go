/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/auth"
	"gophKeeper/internal/client/model/type/bin"
	"gophKeeper/internal/client/model/type/card"
	"gophKeeper/internal/client/model/type/text"

	"github.com/spf13/cobra"
)

func saveDataRun(data model.Model, save func(data model.Model) (err error)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		defer data.Reset()
		data.GetKey()
		if d, ok := data.GetDst().(model.Sanitisable); ok {
			d.Sanitize()
		}
		if d, ok := data.(model.GetFile); ok {
			err := d.GetFile()
			if err != nil {
				fmt.Println("GetFile: Error: ", err)
			}
		}
		// err := data.Validate()
		// if err != nil {
		// 	fmt.Println("Validate: Error: ", err)
		// }

		// todo is draft yet
		// fmt.Println(data.Data)

		err := save(data)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func validArgs(m model.Data) (validArgs []string) {
	validArgsCommon, err := GenFlags(&model.Common{})
	if err != nil {
		fmt.Println(err)
		return
	}
	validArgsData, err := GenFlags(m)
	if err != nil {
		fmt.Println(err)
		return
	}

	validArgs = append(validArgsCommon, validArgsData...)
	return
}

func (a *app) saveAuthCmd() (cmd *cobra.Command) {
	debug := false
	data := auth.New()

	cmd = &cobra.Command{
		Use:       "auth",
		Short:     "Save auth data",
		ValidArgs: validArgs(&auth.Data{}),
		Long:      `Encrypts login/password pairs.`,
		Example: `  save auth -l login -p password
  save auth -l login -p password -d site.com
  save auth -l login -p password -k "my-key-name" -d site.com
`,
		Run: saveDataRun(data, a.Srv().Save),
	}
	err := modelGenerateFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save auth modelGenerateFlags error: %s\n", err)
	}

	return
}

func (a *app) saveTextCmd() (cmd *cobra.Command) {
	debug := false
	data := text.New()
	cmd = &cobra.Command{
		Use:   "text [flags]",
		Short: "Save text data",
		// Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		ValidArgs: validArgs(&text.Data{}),
		Long:      `Encrypts text data`,
		Example: `  save text -f filename
  save text -k custom-key -d description -s
`,
		Run: saveDataRun(data, a.Srv().Save),
	}
	err := modelGenerateFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save text modelGenerateFlags error: %s\n", err)
	}

	return
}

func (a *app) saveBinCmd() (cmd *cobra.Command) {
	debug := false
	data := bin.New()

	cmd = &cobra.Command{
		Use:       "bin [flags]",
		Short:     "Save binary data",
		ValidArgs: validArgs(&auth.Data{}),
		Long:      `Encrypts bin data`,
		Example:   `   save bin -f filename`,
		Run:       saveDataRun(data, a.Srv().Save),
	}
	err := modelGenerateFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save bin modelGenerateFlags error: %s\n", err)
	}
	return
}

func (a *app) saveCardCmd() (cmd *cobra.Command) {
	debug := false
	data := card.New()

	cmd = &cobra.Command{
		Use:       "card [flags]",
		Short:     "Save card data",
		Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		ValidArgs: validArgs(&card.Data{}),
		Long:      `Encrypts bank cards data`,
		Example:   `  save card --num 2222-4444-5555-1111 --exp 10/29 --cvv 123 --owner "Max Space"`,
		Run:       saveDataRun(data, a.Srv().Save),
	}
	err := modelGenerateFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save card modelGenerateFlags error: %s\n", err)
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
