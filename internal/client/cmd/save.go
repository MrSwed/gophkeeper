/*
This package provides commands for saving various types of data,
such as authentication data, text, binary data, and card data.
It uses the Cobra library to define commands and subcommands.

Main functionalities include:

- Saving authentication data.
- Saving text data.
- Saving binary data.
- Saving card data.
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

// generateSaveFlags generates flags for a given model.
// It takes a destination object, a command, and a debug flag as input.
// It generates flags for the base and destination objects of the model.
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

// saveDataRun takes a model as input and returns a function that is used to save the data.
// The returned function first resets the model, gets the key for the data, sanitizes the destination object if it is sanitizable, and then saves the data using the "Save" method of the server.
// If any of these steps fail, an error message is printed.
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

// saveAuthCmd returns a Cobra command for saving the authentication model.
// The command encrypts login/password pairs and saves them.
// Examples of using the command include:
//
//	save auth -l login -p password
//	save auth -l login -p password -d site.com
//	save auth -l login -p password -k "my-key-name" -d site.com
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

// saveTextCmd returns a Cobra command for saving the text model.
// The command encrypts text data and saves it.
// Examples of using the command include:
//
//	save text -f filename
//	save text -k custom-key -d description -s
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

// saveBinCmd returns a Cobra command for saving the binary model.
// The command encrypts binary data and saves it.
// Example of using the command:
//
//	save bin -f filename
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

// saveCardCmd returns a Cobra command for saving the card model.
// The command encrypts bank card data and saves it.
// Example of using the command:
//
//	save card --num 2222-4444-5555-1111 --exp 10/29 --cvv 123 --owner "Max Space"
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

// addSaveCmd adds commands for the data saving operation to the root command.
// Subcommands include commands for saving authentication data,
// text data, binary data, and card data.
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
