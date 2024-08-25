/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/model/type/auth"
	"gophKeeper/internal/client/model/type/bin"
	"gophKeeper/internal/client/model/type/card"
	"gophKeeper/internal/client/model/type/text"
	"gophKeeper/internal/helper"

	"github.com/spf13/cobra"
)

func generateSaveFlags(dst any, cmd *cobra.Command, debug *bool) (err error) {
	if debug != nil {
		cmd.Flags().BoolVarP(debug, "debug", "", *debug, "debug flag")
	}
	if common, ok := dst.(model.Base); ok {
		err = helper.GenerateFlags(common.GetBase(), cmd.Flags())
	}
	if data, ok := dst.(model.Data); ok {
		err = errors.Join(err, helper.GenerateFlags(data.GetDst(), cmd.Flags()))
	}
	return
}

func (a *app) saveDataRun(data model.Model) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		defer data.Reset()
		data.GetKey()
		if d, ok := data.GetDst().(model.Sanitisable); ok {
			d.Sanitize()
		}

		err := a.Srv().Save(data)
		if err != nil {
			cmd.Println(err.Error())
			return
		}
		cmd.Println("Data saved successfully")
	}
}

// todo; clean ?
/** /
func validSaveArgs(m model.Data) (validArgs []string) {
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
/**/

// saveAuthCmd
// Cobra command for save auth model
func (a *app) saveAuthCmd() (cmd *cobra.Command) {
	debug := false
	data := auth.New()

	cmd = &cobra.Command{
		Use:   "auth",
		Short: "Save auth data",
		// ValidArgs: validSaveArgs(&auth.Data{}),
		Long: `Encrypts login/password pairs.`,
		Example: `  save auth -l login -p password
  save auth -l login -p password -d site.com
  save auth -l login -p password -k "my-key-name" -d site.com
`,
		Run: a.saveDataRun(data),
	}
	err := generateSaveFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save auth generateSaveFlags error: %s\n", err)
	}

	return
}

// saveTextCmd
// Cobra command for save text model
func (a *app) saveTextCmd() (cmd *cobra.Command) {
	debug := false
	data := text.New()
	cmd = &cobra.Command{
		Use:   "text [flags]",
		Short: "Save text data",
		// // Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		// ValidArgs: validSaveArgs(&text.Data{}),
		Long: `Encrypts text data`,
		Example: `  save text -f filename
  save text -k custom-key -d description -s
`,
		Run: a.saveDataRun(data),
	}
	err := generateSaveFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save text generateSaveFlags error: %s\n", err)
	}

	return
}

// saveBinCmd
// Cobra command for save bin model
func (a *app) saveBinCmd() (cmd *cobra.Command) {
	debug := false
	data := bin.New()

	cmd = &cobra.Command{
		Use:   "bin [flags]",
		Short: "Save binary data",
		// ValidArgs: validSaveArgs(&auth.Data{}),
		Long:    `Encrypts bin data`,
		Example: `   save bin -f filename`,
		Run:     a.saveDataRun(data),
	}
	err := generateSaveFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save bin generateSaveFlags error: %s\n", err)
	}
	return
}

// saveCardCmd
// Cobra command for save card model
func (a *app) saveCardCmd() (cmd *cobra.Command) {
	debug := false
	data := card.New()

	cmd = &cobra.Command{
		Use:   "card [flags]",
		Short: "Save card data",
		// Args:      cobra.MatchAll(cobra.RangeArgs(0, 4), cobra.OnlyValidArgs),
		// ValidArgs: validSaveArgs(&card.Data{}),
		Long:    `Encrypts bank cards data`,
		Example: `  save card --num 2222-4444-5555-1111 --exp 10/29 --cvv 123 --owner "Max Space"`,
		Run:     a.saveDataRun(data),
	}
	err := generateSaveFlags(data, cmd, &debug)
	if err != nil {
		cmd.Printf("save card generateSaveFlags error: %s\n", err)
	}

	return
}

// addSaveCmd
// Cobra commands for save data operation
func (a *app) addSaveCmd() *app {
	var saveCmd = &cobra.Command{
		Use:   "save [command]",
		Short: "Save data",
		// ValidArgs: []string{},
		Args: cobra.NoArgs,
		Long: `Encrypts and save data`,
	}

	saveCmd.AddCommand(a.saveAuthCmd(), a.saveTextCmd(), a.saveBinCmd(), a.saveCardCmd())

	a.root.AddCommand(saveCmd)
	return a
}
