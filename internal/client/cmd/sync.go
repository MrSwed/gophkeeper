/*
This package provides command-line interface functionalities for synchronizing
with a remote server, managing user registration, and handling user data.

The package uses the Cobra library to define commands and subcommands for
interacting with the synchronization server.

Main functionalities include:

- Registering a client with the server.
- Performing data synchronization.
- Changing the server password.
- Deleting the user account from the server.
*/
package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	cfg "gophKeeper/internal/client/config"
	"gophKeeper/internal/client/crypt"
	"gophKeeper/internal/client/input/password"
	"gophKeeper/internal/client/model"
	"gophKeeper/internal/client/sync"
	"time"

	"github.com/spf13/cobra"
)

// validateServerConfigSet checks if the server configuration is set.
// It returns true if the server address is configured, otherwise it prints an error message.
func (a *app) validateServerConfigSet(cmd *cobra.Command) (ok bool) {
	if ok = cfg.User.Get("server") != nil; !ok {
		cmd.Println(`
The address of the synchronization server is not set, please set it by command
  config user --server <address:port>`)
		cmd.Println()
	}
	return
}

// getSyncToken retrieves the synchronization token from the configuration.
// It decrypts the token and returns it as a byte slice. If the token cannot be retrieved
// or decrypted, it prints an error message and returns nil.
func (a *app) getSyncToken(cmd *cobra.Command) (syncToken []byte) {
	encryptedSyncToken := cfg.User.GetString("sync.token")
	if encryptedSyncToken == "" {
		cmd.Println(`
This client is not registered on the server yet, please run the server registration command 
  sync register`)
		cmd.Println()

		return
	}
	encryptedSyncTokenB, err := hex.DecodeString(encryptedSyncToken)
	if err != nil {
		cmd.PrintErrf("failed to hex Decode String encryptedSyncToken: %v\n", err)
		return
	}
	cryptToken, err := a.Srv().GetToken()
	if err != nil {
		cmd.PrintErrf("failed to get encryption token: %v\n", err)
		return
	}
	syncToken, err = crypt.Decode(encryptedSyncTokenB, cryptToken)
	if err != nil {
		cmd.PrintErrf("failed to decrypt synchronization token: %v\n", err)
		return
	}

	return
}

// addSyncCmd adds the sync command and its subcommands to the root command.
// Subcommands include:
// - register: Registers the client with the remote server.
// - now: Performs immediate synchronization with the server.
// - password: Changes the synchronization password on the server.
// - delete: Deletes the user account from the server.
func (a *app) addSyncCmd() *app {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync with remote server",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.UserLoad()
			if err != nil {
				cmd.PrintErrf("failed to load config: %v\n", err)
				return
			}
			if cfg.User.Get("sync.status") == nil {
				cmd.Println(`
Synchronization has not been performed yet. 
You can run synchronization by command 
  sync now`)
				cmd.Println()
				// _ = cmd.Usage()
				return
			}
			syncInfo, err := json.MarshalIndent(cfg.User.Get("sync.status"), "", " ")
			if err != nil {
				cmd.PrintErrf("failed to marshal sync.status: %v\n", err)
				cmd.Println("Synchronization status raw:", cfg.User.Get("sync.status"))
				return
			}
			cmd.Println("Synchronization status:", string(syncInfo))
		},
	}
	syncCmd.AddCommand(
		&cobra.Command{
			Use:   "register",
			Short: "Register at remote server",
			Run:   a.syncRegisterCmd(),
		},
		&cobra.Command{
			Use:   "now",
			Short: "Sync now",
			Run:   a.syncNowCmd(),
		},
		&cobra.Command{
			Use:   "password",
			Short: "Change server password",
			Run:   a.syncPasswordCmd(),
		},
		&cobra.Command{
			Use:   "delete",
			Short: "Delete server account",
			Run:   a.syncDeleteCmd(),
		})
	a.root.AddCommand(syncCmd)

	return a
}

// syncRegisterCmd returns a function that handles the registration of the client
// with the remote server. It prompts the user for their email and synchronization password,
// sends a registration request, and handles the response, including synchronization token
// storage and user data synchronization if necessary.
func (a *app) syncRegisterCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := cfg.UserLoad()
		if err != nil {
			cmd.PrintErrf("failed to load config: %v\n", err)
		}

		if !a.validateServerConfigSet(cmd) {
			return
		}

		if cfg.User.GetString("sync.token") != "" {
			cmd.Println(`
This client is already registered on the server. You can run synchronization.`)
			cmd.Println()
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
If you have not yet registered on the server with your email, come up with a new synchronization password, a new account will be created for you, after which this client will be able to synchronize`, cfg.User.Get("email"))
		cmd.Println()

		pass, err := password.GetRawPass(false, cfg.PromptSyncPassword)
		if err != nil {
			cmd.PrintErrf("failed to get password: %v\n", err)
			return
		}
		req := model.RegisterClientRequest{
			Email:    cfg.User.GetString("email"),
			Password: pass,
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), cfg.User.GetDuration("sync.timeout.register"))
		defer cancel()
		syncToken, err := sync.RegisterClient(ctx, cfg.User.GetString("server"), req)
		if err != nil {
			cmd.PrintErrf("failed to register client: %v\n", err)
			return
		}
		cmd.Println(`
The synchronization token has been successfully received...`)
		cmd.Println()

		// Check is local encryption key exist
		if cfg.User.GetString("packed_key") == "" {
			cmd.Println(`
Since there is no encryption token, first try to get user data from the server`)
			cmd.Println()
			var syncSrv sync.Service
			ctx, cancel := context.WithTimeout(cmd.Context(), cfg.User.GetDuration("sync.timeout.sync"))
			defer cancel()
			ctx, syncSrv, err = sync.NewSyncService(ctx, cfg.User.GetString("server"), syncToken, a.Srv())
			if err != nil {
				cmd.PrintErrf("synchronization failed: %v\n", err)
				return
			}
			defer syncSrv.Close()
			var updated bool
			updated, err = syncSrv.SyncUser(ctx, "")
			if err != nil {
				cmd.PrintErrf("synchronization failed: %v\n", err)
				return
			}
			if !updated {
				cmd.Println(`User sync finished, no new data received`)
			} else {
				cmd.Println(`User sync finished, user data updated from server`)
			}
			// Steel no encryption key
			// Notify about create new one
			if cfg.User.GetString("packed_key") == "" {
				cmd.Println(`
To save client token to config, needs to create an encryption key, 
please come up with a password for this..`)
				cmd.Println()
			}
		}

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
		cfg.User.Set("sync.token", hex.EncodeToString(encryptedSyncToken))
		cfg.User.Set("sync.status.token.created_at", time.Now())
		cmd.Println(`
Congratulations! The client is successfully registered on the server, the synchronization token is saved in the settings.`)
		cmd.Println()
	}
}

// syncNowCmd returns a function that handles immediate synchronization with the server.
// It retrieves the sync token, connects to the server, and synchronizes user data and
// other relevant information.
func (a *app) syncNowCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := cfg.UserLoad()
		if err != nil {
			cmd.PrintErrf("failed to load config: %v\n", err)
		}
		if !a.validateServerConfigSet(cmd) {
			return
		}
		syncToken := a.getSyncToken(cmd)
		if len(syncToken) == 0 {
			return
		}

		cmd.Println(time.Now().Format(time.DateTime), `Start synchronization with server`)

		ctx, cancel := context.WithTimeout(cmd.Context(), cfg.User.GetDuration("sync.timeout.sync"))
		defer cancel()

		var syncSrv sync.Service
		ctx, syncSrv, err = sync.NewSyncService(ctx, cfg.User.GetString("server"), syncToken, a.Srv())
		if err != nil {
			cmd.PrintErrf("prepare synchronization failed: %v\n", err)
			return
		}
		defer syncSrv.Close()

		var updated bool
		updated, err = syncSrv.SyncUser(ctx, "")
		if err != nil {
			cmd.PrintErrf("user synchronization failed: %v\n", err)
			return
		}
		if !updated {
			cmd.Println(`User sync finished, no new data received`)
		} else {
			cmd.Println(`User sync finished, user data updated from server`)
		}

		cmd.Println(time.Now().Format(time.DateTime), `User synchronization finished`)

		err = syncSrv.SyncData(ctx)
		if err != nil {
			cmd.PrintErrf("data synchronization failed: %v\n", err)
			return
		}
		cmd.Println(time.Now().Format(time.DateTime), `Data synchronization finished`)
		syncInfo, err := json.MarshalIndent(cfg.User.Get("sync.status"), "", " ")
		if err != nil {
			cmd.PrintErrf("failed to marshal sync.status: %v\n", err)
			cmd.Println("Synchronization status raw:", cfg.User.Get("sync.status"))
			return
		}
		cmd.Println("Synchronization status:", string(syncInfo))
	}
}

// syncPasswordCmd returns a function that handles changing the synchronization password
// on the server. It prompts the user for a new password, sends the request, and handles
// the response.
func (a *app) syncPasswordCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := cfg.UserLoad()
		if err != nil {
			cmd.PrintErrf("failed to load config: %v\n", err)
		}

		if !a.validateServerConfigSet(cmd) {
			return
		}

		syncToken := a.getSyncToken(cmd)
		if len(syncToken) == 0 {
			return
		}

		cmd.Println(time.Now().Format(time.DateTime), `connecting to server..`)

		ctx, cancel := context.WithTimeout(cmd.Context(), cfg.User.GetDuration("sync.timeout.sync"))
		defer cancel()

		var syncSrv sync.Service
		ctx, syncSrv, err = sync.NewSyncService(ctx, cfg.User.GetString("server"), syncToken, a.Srv())
		if err != nil {
			cmd.PrintErrf("prepare synchronization failed: %v\n", err)
			return
		}
		defer syncSrv.Close()

		var pass string
		pass, err = password.GetRawPass(true, cfg.PromptNewSyncPassword, cfg.PromptSyncConfirmPassword)
		if err != nil {
			cmd.PrintErrf("failed to get password: %v\n", err)
			return
		}

		var updated bool
		updated, err = syncSrv.SyncUser(ctx, pass)
		if err != nil {
			cmd.PrintErrf("user synchronization failed: %v\n", err)
			return
		}
		if !updated {
			cmd.Println(`User sync finished, no new data received`)
		} else {
			cmd.Println(`User sync finished, user data updated from server`)
		}

		cmd.Println(time.Now().Format(time.DateTime), `User synchronization finished`)
	}
}

// syncDeleteCmd returns a function that handles the deletion of the user account from
// the server. It connects to the server and sends a delete request, then removes the
// synchronization token from local configuration.
func (a *app) syncDeleteCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := cfg.UserLoad()
		if err != nil {
			cmd.PrintErrf("failed to load config: %v\n", err)
		}

		if !a.validateServerConfigSet(cmd) {
			return
		}

		syncToken := a.getSyncToken(cmd)
		if len(syncToken) == 0 {
			return
		}

		cmd.Println(time.Now().Format(time.DateTime), `connecting to server..`)

		ctx, cancel := context.WithTimeout(cmd.Context(), cfg.User.GetDuration("sync.timeout.sync"))
		defer cancel()

		var syncSrv sync.Service
		ctx, syncSrv, err = sync.NewSyncService(ctx, cfg.User.GetString("server"), syncToken, a.Srv())
		if err != nil {
			cmd.PrintErrf("prepare synchronization failed: %v\n", err)
			return
		}
		defer syncSrv.Close()

		// todo : need confirm ?
		err = syncSrv.DeleteUser(ctx)
		if err != nil {
			cmd.PrintErrf("Error delete user from server: %v\n", err)
			return
		}
		delete(cfg.User.Get("sync").(map[string]any), "token")

		cmd.Println(`
The user has been successfully deleted from the server, 
the client token has been deleted from the local configuration`)
	}
}
