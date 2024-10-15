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

func (a *app) validateServerConfigSet(cmd *cobra.Command) (ok bool) {
	if ok = cfg.User.Get("server") != nil; !ok {
		cmd.Println(`
The address of the synchronization server is not set, please set it by command
  config user --server <address:port>`)
		cmd.Println()
	}
	return
}

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
		cmd.PrintErrf("failed to get encription token: %v\n", err)
		return
	}
	syncToken, err = crypt.Decode(encryptedSyncTokenB, cryptToken)
	if err != nil {
		cmd.PrintErrf("failed to decrypt synchronization token: %v\n", err)
		return
	}

	return
}

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

func (a *app) syncRegisterCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := cfg.UserLoad()
		if err != nil {
			cmd.PrintErrf("failed to load config: %v\n", err)
		}

		if !a.validateServerConfigSet(cmd) {
			return
		}

		if cfg.User.Get("sync.token") != nil {
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

		pass, err := password.GetRawPass(false, "Please enter the server synchronization password: ")
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
			var syncSrv sync.SyncService
			ctx, cancel := context.WithTimeout(cmd.Context(), cfg.User.GetDuration("sync.timeout.sync"))
			defer cancel()
			ctx, syncSrv, err = sync.NewSyncService(ctx, cfg.User.GetString("server"), syncToken, a.Srv())
			if err != nil {
				cmd.PrintErrf("sychronization failed: %v\n", err)
				return
			}
			defer syncSrv.Close()
			var updated bool
			updated, err = syncSrv.SyncUser(ctx, "")
			if err != nil {
				cmd.PrintErrf("sychronization failed: %v\n", err)
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

		var syncSrv sync.SyncService
		ctx, syncSrv, err = sync.NewSyncService(ctx, cfg.User.GetString("server"), syncToken, a.Srv())
		if err != nil {
			cmd.PrintErrf("prepare sychronization failed: %v\n", err)
			return
		}
		defer syncSrv.Close()

		var updated bool
		updated, err = syncSrv.SyncUser(ctx, "")
		if err != nil {
			cmd.PrintErrf("user sychronization failed: %v\n", err)
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
			cmd.PrintErrf("data sychronization failed: %v\n", err)
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

func (a *app) syncPasswordCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := cfg.UserLoad()
		if err != nil {
			cmd.PrintErrf("failed to load config: %v\n", err)
		}

		if !a.validateServerConfigSet(cmd) {
			return
		}
		// todo
		cmd.Println(`
change password Not implemented yet `)
		cmd.Println()
	}
}

func (a *app) syncDeleteCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := cfg.UserLoad()
		if err != nil {
			cmd.PrintErrf("failed to load config: %v\n", err)
		}
		if !a.validateServerConfigSet(cmd) {
			return
		}

		// todo
		cmd.Println(`
delete account Not implemented yet `)
		cmd.Println()

	}
}
