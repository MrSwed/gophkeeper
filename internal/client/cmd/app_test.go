package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	cfg "gophKeeper/internal/client/config"
	serverApp "gophKeeper/internal/server/app"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	waitPortInterval    = 100 * time.Millisecond
	waitPortConnTimeout = 50 * time.Millisecond
)

type appTestSuite struct {
	suite.Suite
	client appTestClient
	server appTestServer
}

type appTestClient struct {
	oldStdin, stdin, stdinPipe, oldStdOut, stdout, stdoutPipe *os.File
	outC                                                      chan string
}

type appTestServer struct {
	ctx     context.Context
	stop    context.CancelFunc
	address string
	pgCont  *postgres.PostgresContainer
	osArgs  []string
}

func (s *appTestSuite) getServerAddress() string {
	s.maybeSetupServer()
	return s.server.address
}

func (s *appTestSuite) SetupSuite() {
	s.setupClient()
}

func (s *appTestSuite) setupClient() {
	cfg.Glob.Viper.Set("config_path", s.T().TempDir())
	var err error
	s.client.stdin, s.client.stdinPipe, err = os.Pipe()
	require.NoError(s.T(), err)
	s.client.oldStdin, os.Stdin = os.Stdin, s.client.stdin
}

func (s *appTestSuite) maybeSetupServer() {
	if s.server.pgCont != nil {
		return
	}
	var (
		err error
	)
	s.server.osArgs = os.Args
	os.Args = os.Args[0:1]
	s.server.ctx, s.server.stop = context.WithCancel(context.Background())

	s.server.pgCont, err = createPostgresContainer(s.server.ctx)
	require.NoError(s.T(), err)
	databaseDSN, err := s.server.pgCont.ConnectionString(s.server.ctx, "sslmode=disable")
	require.NoError(s.T(), err)
	s.server.address = net.JoinHostPort("", fmt.Sprintf("%d", rand.Intn(200)+30000))

	s.T().Setenv("DATABASE_DSN", databaseDSN)
	s.T().Setenv("GRPC_ADDRESS", s.server.address)
	s.T().Setenv("GRPC_OPERATION_TIMEOUT", "5000s")

	go serverApp.RunApp(s.server.ctx, nil, nil, serverApp.BuildMetadata{Version: "testing..", Date: time.Now().String(), Commit: ""})
	require.NoError(s.T(), waitGRPCPort(s.server.ctx, s.server.address))

	// db, err := sqlx.Connect("postgres", databaseDSN)
	// predefined, err := os.ReadFile(filepath.Join("../../../", "testdata", "server.sql"))
	// require.NoError(s.T(), err)
	// _, err = db.Exec(string(predefined))
	// require.NoError(s.T(), err)
}

func (s *appTestSuite) TearDownSuite() {
	s.tearDownClient()
	s.tearDownServer()
}

func (s *appTestSuite) tearDownClient() {
	// restore stdin
	// require.NoError(s.T(), os.RemoveAll(s.T().TempDir()))

	os.Stdin = s.client.oldStdin
	err := s.client.stdinPipe.Close()
	require.NoError(s.T(), err)
	err = s.client.stdin.Close()
	require.NoError(s.T(), err)
}

func (s *appTestSuite) tearDownServer() {
	if s.server.pgCont != nil {
		if err := s.server.pgCont.Terminate(s.server.ctx); err != nil {
			log.Fatalf("error terminating postgres container: %s", err)
		}

		s.server.stop()
		os.Args = s.server.osArgs
	}
}

func TestApp(t *testing.T) {
	suite.Run(t, new(appTestSuite))
}

func (s *appTestSuite) Test_App() {
	t := s.T()
	tests := []struct {
		name         string
		commands     [][]string
		wantStrOut   [][]string
		wantNoStrOut [][]string
		inputs       [][]string
	}{
		{
			name:     "no args",
			commands: [][]string{{""}},
		}, {
			name:       "help",
			commands:   [][]string{{"--help"}},
			wantStrOut: [][]string{{"Available Commands"}},
		}, {
			name:       "profile",
			commands:   [][]string{{"profile"}},
			wantStrOut: [][]string{{"Current profile"}},
		}, {
			name:     "profile --help",
			commands: [][]string{{"profile", "--help"}},
			wantStrOut: [][]string{{"Available Commands",
				"list        list of profiles",
				"use         switch to another profile"}},
		}, {
			name:       "profile list",
			commands:   [][]string{{"profile", "list"}},
			wantStrOut: [][]string{{"No profiles yet"}},
		}, {
			name:       "profile use",
			commands:   [][]string{{"profile", "use"}},
			wantStrOut: [][]string{{"Error: accepts 1 arg(s), received 0"}},
		}, {
			name: "profile use test, save, list and view, change password",
			commands: [][]string{
				0:  {"profile", "use", "testName"},
				1:  {"save", "text", "-t", "some text data", "-d", "some description", "-k", "text-key-1"},
				2:  {"list"},
				3:  {"view", "text-key-1"},
				4:  {"list", "-k", "text"},
				5:  {"profile", "list"},
				6:  {"profile", "password"},
				7:  {"view", "text-key-1"},
				8:  {"profile", "use", "someOtherUser"},
				10: {"profile", "password"},
				11: {"profile", "password"},
				12: {"save", "text", "-t", "some text data", "-k", "text-key-1"},
			},
			wantStrOut: [][]string{
				0:  {"Switching to profile..  " + "testName"},
				1:  {"Data saved successfully"},
				2:  {"text-key-1", "some description", "Total:"},
				3:  {"text-key-1", "some description", "some text data"},
				4:  {"text-key-1", "some description", "Total:"},
				5:  {"Available profiles", "- testName"},
				6:  {"User testName configuration loaded", cfg.PromptMasterPassword, cfg.PromptNewMasterPassword, cfg.PromptConfirmMasterPassword},
				7:  {"text-key-1", "some description", "some text data"},
				8:  {"Switching to profile..  " + "someOtherUser"},
				10: {"Current profile someOtherUser", "User someOtherUser configuration loaded", cfg.PromptNewMasterPassword, cfg.PromptConfirmMasterPassword},
				11: {"Current profile someOtherUser", "User someOtherUser configuration loaded", cfg.PromptMasterPassword, cfg.PromptNewMasterPassword, cfg.PromptConfirmMasterPassword},
				12: {"Data saved successfully"},
			},
			inputs: [][]string{
				0:  {},
				1:  {"somePass", "somePass"},
				2:  {},
				3:  {"somePass"},
				4:  {},
				5:  {},
				6:  {"somePass", "newPass", "newPass"},
				7:  {},
				8:  {},
				10: {"newProfPass", "newProfPass"},
				11: {"newProfPass", "newNewProfPass", "newNewProfPass"},
			},
		}, {
			name: "profile use default, save, list, view, delete",
			commands: [][]string{
				0: {"profile", "use", "default"},
				1: {"save", "card", "--num", "0000-0000-0000-0001", "--cvv", "222", "-k", "card-key-1"},
				2: {"save", "card", "--num", "0000-0000-0000-0000", "--cvv", "222", "-k", "card-key-1"},
				3: {"list"},
				4: {"view"},
				5: {"view", "card-key-1"},
				6: {"delete", "card-key-1"},
				7: {"view", "card-key-1"},
				8: {"delete", "card-key-1"},
			},
			wantStrOut: [][]string{
				0: {"Switching to profile..  ", "default"},
				1: {"Error:Field validation for 'Number'"},
				2: {"Data saved successfully"},
				3: {"card-key-1", "Total:"},
				4: {"Usage:", "view <key name> [flags]"},
				5: {"card-key-1", "0000 0000 0000 0000", "222"},
				6: {"card-key-1 success deleted"},
				7: {"Record not exist: card-key-1"},
				8: {"Record not exist: card-key-1"},
			},
			inputs: [][]string{
				0: {},
				1: {"somePass", "somePass"},
				2: {"somePass"},
			},
		}, {
			name: "profile use test2, config",
			commands: [][]string{
				{"profile", "use", "test2"},
				{"config", "user", "-e", "some@email.net"},
			},
			wantStrOut: [][]string{
				{"Switching to profile..  ", "test2"},
				{"some@email.net", "User configuration: set", "Success autosave config"},
			},
		}, {
			name: "test configs",
			commands: [][]string{
				0: {"profile", "use", "newNameConfig"},
				1: {"config", "save"},
				2: {"config", "global"},
				3: {"config", "user"},
				4: {"config", "global", "-a"},
				5: {"config", "user", "-a"},
				6: {"save", "text", "-t", "Some text data save", "-k", "Test-key-1"},
				7: {"config", "user"},
			},
			wantStrOut: [][]string{
				0: {"Switching to profile..  ", "newNameConfig"},
				1: {"Saving global config.. not changed", "Saving user config.. not changed"},
				2: {"Global configuration:", `"autosave"`, `"config_path"`, `"loaded_at"`, `"profile"`, `newNameConfig`},
				3: {"User configuration:", `"db_file"`, `"name"`, `newNameConfig`},
				4: {"Global configuration: set `autosave` = `true`"},
				5: {"User configuration: set `autosave` = `true`"},
				6: {"Data saved successfully"},
				7: {"User configuration:", `"db_file"`, `"loaded_at"`, `"name"`, `newNameConfig`, `packed_key`},
			},
			wantNoStrOut: [][]string{
				0: {},
				1: {},
				2: {},
				3: {"packed_key", "encryption_key"},
				4: {},
				5: {},
				6: {},
				7: {"encryption_key"},
			},
			inputs: [][]string{
				0: {},
				1: {},
				2: {},
				3: {},
				4: {},
				5: {},
				6: {"somePass", "somePass"},
			},
		},
		{
			name: "sync new user",
			commands: [][]string{
				0:  {"profile", "use", "newSyncUser"},
				1:  {"save", "card", "--num", "0000-0000-0000-0000", "--cvv", "222", "-k", "card-key-1"},
				2:  {"save", "text", "--text", "some text data"},
				3:  {"sync"},
				4:  {"sync", "now"},
				5:  {"config", "user", "--server", s.getServerAddress()},
				6:  {"sync", "now"},
				7:  {"sync", "register"},
				8:  {"config", "user", "-e", "newSyncUser@email.localhost"},
				9:  {"sync", "register"},
				10: {"sync", "register"},
				11: {"sync", "register"},
				12: {"sync", "now"},
				13: {"sync"},
			},
			wantStrOut: [][]string{
				0:  {"Switching to profile..  ", "newSyncUser"},
				1:  {"Data saved successfully"},
				2:  {"Data saved successfully"},
				3:  {"Synchronization has not been performed yet.", "You can run synchronization by command"},
				4:  {"The address of the synchronization server is not set"},
				5:  {"User configuration: set `server` = `" + s.getServerAddress() + "`"},
				6:  {"This client is not registered on the server yet"},
				7:  {"The email is not specified in the configuration"},
				8:  {"User configuration: set `email` = `newSyncUser@email.localhost`"},
				9:  {"Registering this client at server with email newSyncUser@email.localhost", "If you have not yet registered on the server with your email, come up with a new synchronization password, a new account will be created for you, after which this client will be able to synchronize", "failed to register client", "validation", "password"},
				10: {"Registering this client at server with email newSyncUser@email.localhost", "If you have not yet registered on the server with your email, come up with a new synchronization password, a new account will be created for you, after which this client will be able to synchronize", "The synchronization token has been successfully received", "Congratulations! The client is successfully registered on the server, the synchronization token is saved in the settings."},
				11: {"This client is already registered on the server. You can run synchronization"},
				12: {"Start synchronization with server", "User sync finished, no new data received", "User synchronization finished", "Data synchronization finished", "Synchronization status"},
				13: {"Synchronization status:", "last_sync_at", "last_sync_at", `"updated": 2`},
			},
			wantNoStrOut: [][]string{
				0:  {},
				1:  {},
				2:  {},
				3:  {},
				4:  {},
				5:  {},
				6:  {},
				7:  {},
				8:  {},
				9:  {},
				10: {"This client is already registered on the server. You can run synchronization.", "The email is not specified in the configuration, please set it by command", "failed to get password", "failed to register client"},
				11: {},
				12: {"failed to load config", "prepare synchronization failed", "user synchronization failed", "failed to marshal sync.status", "User sync finished, user data updated from server", "data synchronization failed", "failed to hex Decode String encryptedSyncToken: invalid byte", "failed to get encryption token", "failed to decrypt synchronization token"},
			},
			inputs: [][]string{
				0:  {},
				1:  {"somePass", "somePass"},
				2:  {"somePass"},
				3:  {},
				4:  {},
				5:  {},
				6:  {},
				7:  {},
				8:  {},
				9:  {"simplePass"},
				10: {`Pa$$w0rd`},
				11: {},
				12: {"somePass"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, cmd := range tt.commands {
				if i < len(tt.inputs) {
					s.input(tt.inputs[i]...)
				}
				rescueStdout := os.Stdout
				r, w, err := os.Pipe()
				require.NoError(t, err)
				os.Stdout = w
				if len(cmd) > 0 && cmd[0] == "sync" {
					s.maybeSetupServer()
				}
				consoleOutput, _ := s.client.executeCommand(cmd...)

				err = w.Close()
				require.NoError(t, err)
				out, err := io.ReadAll(r)
				require.NoError(t, err)
				os.Stdout = rescueStdout
				if i < len(tt.wantStrOut) {
					for _, wantOut := range tt.wantStrOut[i] {
						require.Contains(t, consoleOutput+string(out), wantOut, fmt.Sprintf("cmd %d: %s", i, cmd))
					}
				}
				if i < len(tt.wantNoStrOut) {
					for _, wantNoOut := range tt.wantNoStrOut[i] {
						require.NotContains(t, consoleOutput+string(out), wantNoOut, fmt.Sprintf("cmd %d: %s", i, cmd))
					}
				}
			}
		})
	}
}

func (s *appTestSuite) input(str ...string) {
	if len(str) > 0 {
		input := []byte(strings.Join(str, "\n") + "\n")
		_, err := s.client.stdinPipe.Write(input)
		require.NoError(s.T(), err)
	}
}

// executeCommand
// https://github.com/spf13/cobra/issues/1790#issuecomment-2121139148
func (s *appTestClient) executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	a := NewApp(BuildMetadata{
		Version: "N/A",
		Date:    time.Now().UTC().Format(time.RFC3339),
		Commit:  "N/A",
	})
	a.root.SetOut(buf)
	a.root.SetErr(buf)
	a.root.SetArgs(args)

	err := a.Execute()
	// _ = a.Close()
	return buf.String(), err
}

func createPostgresContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx, "postgres:14-alpine",
		// postgres.WithInitScripts(
		// 	filepath.Join("../../../", "testdata", "server.sql"),
		// ),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	return pgContainer, nil
}
func waitGRPCPort(ctx context.Context, address string) error {
	if address == "" {
		return nil
	}
	ticker := time.NewTicker(waitPortInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			conn, _ := net.DialTimeout("tcp", address, waitPortConnTimeout)
			if conn != nil {
				_ = conn.Close()
				return nil
			}
		}
	}
}
