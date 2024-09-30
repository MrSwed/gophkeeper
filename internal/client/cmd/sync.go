package cmd

import (
	"context"
	cfg "gophKeeper/internal/client/config"
	"gophKeeper/internal/client/crypt"
	"gophKeeper/internal/client/input/password"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/sync"
	"time"

	"github.com/spf13/cobra"
)

func isServerConfigSet(cmd *cobra.Command) (ok bool) {
	if ok = cfg.User.Get("server") != nil; !ok {
		cmd.Println(`
The address of the synchronization server is not set, please set it by command
  config user --server <address:port>`)
		cmd.Println()
	}
	return
}

func (a *app) addSyncCmd() *app {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync with remote server",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
			}
			if cfg.User.Get("sync_status") == nil {
				cmd.Println(`
Synchronization has not been performed yet. 
You can run synchronization by command 
  sync now`)
				cmd.Println()
				// _ = cmd.Usage()
				return
			}
			cmd.Println("Synchronization status:", cfg.User.Get("sync_status"))
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "now",
		Short: "Sync now",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
			}
			if cfg.User.Get("sync_token") == nil {
				cmd.Println(`
This client is not registered on the server yet, please run the server registration command 
  sync register`)
				cmd.Println()
				return
			}
			if !isServerConfigSet(cmd) {
				return
			}

			cmd.Println(`Start synchronization with server`, time.Now().Format(time.DateTime))
			// todo sync here
			cmd.Println(`
sync Not implemented yet `)
			cmd.Println()

		},
	}, &cobra.Command{
		Use:   "register",
		Short: "Register at remote server",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
			}

			if cfg.User.Get("sync_token") != nil {
				cmd.Println(`
This client is already registered on the server. You can run synchronization.`)
				cmd.Println()
				return
			}

			if !isServerConfigSet(cmd) {
				return
			}

			if cfg.User.Get("email") == nil {
				cmd.Println(`
The email is not specified in the configuration, please set it by command
  config user --email <email@example.com>`)
				return
			}

			cmd.Printf(`
Registering this client at server with email %s. 
If you have not yet registered on the server with your email, come up with a new synchronization password, a new account will be created for you, after which this client will be able to synchronize.
`, cfg.User.Get("email"))
			cmd.Println()

			pass, err := password.GetRawPass(false, "Please enter the server synchronization password: ")
			if err != nil {
				cmd.PrintErrf("failed to get password: %v\n", err)
				return
			}
			req := model.RegisterClientRequest{
				Email:    cfg.User.GetString("email"),
				Password: pass,
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), cfg.User.GetDuration("timeout"))
			defer cancel()
			syncToken, err := sync.RegisterClient(ctx, cfg.User.GetString("server"), req)
			if err != nil {
				cmd.PrintErrf("failed to register client: %v\n", err)
				return
			}
			cmd.Println(`
The synchronization token has been successfully received, we save it in the config...`)
			cmd.Println()

			cryptToken, err := a.Srv().GetToken()
			if err != nil {
				cmd.PrintErrf("failed to get synchronization token: %v\n", err)
				return
			}
			encryptedSyncToken, err := crypt.Encode(syncToken, cryptToken)
			if err != nil {
				cmd.PrintErrf("failed to crypt synchronization token: %v\n", err)
				return
			}
			cfg.User.Set("sync_token", encryptedSyncToken)
			cmd.Println(`
Congratulations! The client is successfully registered on the server, the synchronization token is saved in the settings.`)
			cmd.Println()
		},
	}, &cobra.Command{
		Use:   "password",
		Short: "Change server password",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
			}

			if !isServerConfigSet(cmd) {
				return
			}
			// todo
			cmd.Println(`
change password Not implemented yet `)
			cmd.Println()
		},
	}, &cobra.Command{
		Use:   "delete",
		Short: "Delete server account",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
			}
			if !isServerConfigSet(cmd) {
				return
			}

			// todo
			cmd.Println(`
delete account Not implemented yet `)
			cmd.Println()

		},
	})

	a.root.AddCommand(cmd)

	return a
}
