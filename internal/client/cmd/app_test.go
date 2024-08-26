package cmd

import (
	"bytes"
	"fmt"
	cfg "gophKeeper/internal/client/config"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type appTestSuite struct {
	suite.Suite
	oldStdin, stdin, stdinPipe, oldStdOut, stdout, stdoutPipe *os.File
	outC                                                      chan string
}

var testDataPath string = filepath.Join("..", "..", "..", "testdata")

// executeCommand
// https://github.com/spf13/cobra/issues/1790#issuecomment-2121139148
func (s *appTestSuite) executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	a := NewApp()
	a.root.SetOut(buf)
	a.root.SetErr(buf)
	a.root.SetArgs(args)

	err := a.Execute()

	return buf.String(), err
}

func (s *appTestSuite) SetupSuite() {
	cfg.Glob.Viper.Set("config_path", s.T().TempDir())
	var err error
	s.stdin, s.stdinPipe, err = os.Pipe()
	require.NoError(s.T(), err)
	s.oldStdin, os.Stdin = os.Stdin, s.stdin
}

func (s *appTestSuite) input(str ...string) {
	if len(str) > 0 {
		input := []byte(strings.Join(str, "\n") + "\n")
		_, err := s.stdinPipe.Write(input)
		require.NoError(s.T(), err)
	}
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
				6:  {"User testName configuration loaded", "Please enter password", "Please enter new password", "Please confirm you password"},
				7:  {"text-key-1", "some description", "some text data"},
				8:  {"Switching to profile..  " + "someOtherUser"},
				10: {"Current profile someOtherUser", "User someOtherUser configuration loaded", "Please enter new password", "Please confirm you password"},
				11: {"Current profile someOtherUser", "User someOtherUser configuration loaded", "Please enter password", "Please enter new password", "Please confirm you password"},
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
				2: {"somePass", "somePass"},
				3: {"somePass"},
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
				{"profile", "use", "newNameConfig"},
				{"config", "save"},
				{"config"},
				{"config", "user"},
				{"save", "text", "-t", "Some text data save", "-k", "Test-key-1"},
				{"config", "user"},
			},
			wantStrOut: [][]string{
				{"Switching to profile..  ", "newNameConfig"},
				{"Saving global config.. not changed", "Saving user config.. not changed"},
				{"Global configuration:", `"autosave"`, `"config_path"`, `"loaded_at"`, `"profile"`, `newNameConfig`},
				{"User configuration:", `"db_file"`, `"name"`, `newNameConfig`},
				{"Data saved successfully"},
				{"User configuration:", `"db_file"`, `"loaded_at"`, `"name"`, `newNameConfig`, `packed_key`},
			},
			wantNoStrOut: [][]string{
				{},
				{},
				{},
				{"packed_key", "encryption_key"},
				{},
				{"encryption_key"},
			},
			inputs: [][]string{
				{},
				{},
				{},
				{},
				{"somePass", "somePass"},
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

				consoleOutput, _ := s.executeCommand(cmd...)

				err = w.Close()
				require.NoError(t, err)
				out, err := io.ReadAll(r)
				require.NoError(t, err)
				os.Stdout = rescueStdout
				if i < len(tt.wantStrOut) {
					for _, wantOut := range tt.wantStrOut[i] {
						require.Contains(t, consoleOutput+string(out), wantOut, append(cmd, fmt.Sprintf(" : %d", i)))
					}
				}
				if i < len(tt.wantNoStrOut) {
					for _, wantNoOut := range tt.wantNoStrOut[i] {
						require.NotContains(t, consoleOutput, wantNoOut, cmd)
					}
				}
			}
		})
	}
}
