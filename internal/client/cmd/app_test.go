package cmd

import (
	"bytes"
	"fmt"
	cfg "gophKeeper/internal/client/config"
	"gophKeeper/internal/client/service"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// https://github.com/spf13/cobra/issues/1790#issuecomment-2121139148

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = ExecuteCommandC(root, args...)
	return output, err
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

type appTestSuite struct {
	suite.Suite
	db                         *sqlx.DB
	app                        *app
	srv                        *service.Service
	oldStdin, stdin, stdinPipe *os.File
	// user                       string
	// userBak                    string
	// pass                       string
}

var testDataPath string = filepath.Join("..", "..", "..", "testdata")

func (s *appTestSuite) SetupSuite() {
	cfg.Glob.Set("config_path", s.T().TempDir())
	s.app = NewApp()
	var err error
	s.stdin, s.stdinPipe, err = os.Pipe()
	require.NoError(s.T(), err)
	s.oldStdin, os.Stdin = os.Stdin, s.stdin
}

func (s *appTestSuite) input(str ...string) {
	input := []byte(strings.Join(str, "\n") + "\n")
	_, err := s.stdinPipe.Write(input)
	require.NoError(s.T(), err)
}

func (s *appTestSuite) TearDownSuite() {

	require.NoError(s.T(), os.RemoveAll(s.T().TempDir()))

	// restore stdin
	os.Stdin = s.oldStdin
	err := s.stdinPipe.Close()
	require.NoError(s.T(), err)
	err = s.stdin.Close()
	require.NoError(s.T(), err)
}

func TestApp(t *testing.T) {
	suite.Run(t, new(appTestSuite))
}

func (s *appTestSuite) Test_App() {
	t := s.T()
	tests := []struct {
		name        string
		commands    [][]string
		outStr      [][]string
		pass        string
		passConfirm string
	}{
		{
			name:     "no args",
			commands: [][]string{{""}},
		}, {
			name:     "help",
			commands: [][]string{{"--help"}},
			outStr:   [][]string{{"Available Commands"}},
		}, {
			name:     "profile",
			commands: [][]string{{"profile"}},
			outStr:   [][]string{{"Current profile"}},
		}, {
			name:     "profile --help",
			commands: [][]string{{"profile", "--help"}},
			outStr: [][]string{{"Available Commands",
				"list        list of profiles",
				"use         switch to another profile"}},
		}, {
			name:     "profile list",
			commands: [][]string{{"profile", "list"}},
			outStr:   [][]string{{"Available profiles", "- default"}},
		}, {
			name:     "profile use",
			commands: [][]string{{"profile", "use"}},
			outStr:   [][]string{{"Error: accepts 1 arg(s), received 0"}},
		}, {
			name: "profile use test, save, list and view",
			commands: [][]string{
				{"profile", "use", "testName"},
				{"save", "text", "-t", "some text data", "-d", "some description", "-k", "text-key-1"},
				{"list"},
				{"view", "text-key-1"},
			},
			outStr: [][]string{
				{"Switching to profile..  " + "testName"},
				{"Data saved successfully"},
				{"text-key-1", "some description", "Total:"},
				{"text-key-1", "some description", "some text data"},
			},
			pass:        "somePass",
			passConfirm: "somePass",
		}, {
			name: "profile use default, save, list and view",
			commands: [][]string{
				{"profile", "use", "default"},
				{"save", "card", "--num", "0000-0000-0000-0001", "--cvv", "222", "-k", "card-key-1"},
				{"save", "card", "--num", "0000-0000-0000-0000", "--cvv", "222", "-k", "card-key-1"},
				{"list"},
				{"view", "card-key-1"},
			},
			outStr: [][]string{
				{"Switching to profile..  ", "default"},
				{"Error:Field validation for 'Number'"},
				{"Data saved successfully"},
				{"card-key-1", "Total:"},
				{"card-key-1", "0000 0000 0000 0000", "222"},
			},
			pass:        "somePass",
			passConfirm: "somePass",
		}, {
			name: "profile use test2, save, list and view",
			commands: [][]string{
				{"profile", "use", "test2"},
				{"config", "save"},
			},
			outStr: [][]string{
				{"Switching to profile..  ", "test2"},
				{"Saving global config.. success", "Saving user config.. not changed"},
			},
			pass:        "somePass",
			passConfirm: "somePass",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.pass != "" {
				s.input(tt.pass)
			}
			if tt.passConfirm != "" {
				s.input(tt.passConfirm)
			}
			for i, cmd := range tt.commands {
				consoleOutput, _ := ExecuteCommand(s.app.root, cmd...)
				fmt.Println(consoleOutput)

				// require.NoError(t, err)
				if i < len(tt.outStr) {
					for _, wantOut := range tt.outStr[i] {
						require.Contains(t, consoleOutput, wantOut)
					}
				}
			}
		})
	}
}
