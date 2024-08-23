package cmd

import (
	"bytes"
	cfg "gophKeeper/internal/client/config"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type appTestSuite struct {
	suite.Suite
	oldStdin, stdin, stdinPipe *os.File
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
			name: "profile use test, save, list and view",
			commands: [][]string{
				{"profile", "use", "testName"},
				{"save", "text", "-t", "some text data", "-d", "some description", "-k", "text-key-1"},
				{"list"},
				{"view", "text-key-1"},
				{"list"},
				{"profile", "list"},
			},
			wantStrOut: [][]string{
				{"Switching to profile..  " + "testName"},
				{"Data saved successfully"},
				{"text-key-1", "some description", "Total:"},
				{"text-key-1", "some description", "some text data"},
				{},
				{"Available profiles", "- testName"},
			},
			inputs: [][]string{
				{},
				{"somePass", "somePass"},
				{},
				{"somePass"},
			},
		}, {
			name: "profile use default, save, list and view",
			commands: [][]string{
				{"profile", "use", "default"},
				{"save", "card", "--num", "0000-0000-0000-0001", "--cvv", "222", "-k", "card-key-1"},
				{"save", "card", "--num", "0000-0000-0000-0000", "--cvv", "222", "-k", "card-key-1"},
				{"list"},
				{"view", "card-key-1"},
			},
			wantStrOut: [][]string{
				{"Switching to profile..  ", "default"},
				{"Error:Field validation for 'Number'"},
				{"Data saved successfully"},
				{"card-key-1", "Total:"},
				{"card-key-1", "0000 0000 0000 0000", "222"},
			},
			inputs: [][]string{
				{},
				{"somePass", "somePass"},
				{"somePass", "somePass"},
				{"somePass"},
			},
		}, {
			name: "profile use test2, save, list and view",
			commands: [][]string{
				{"profile", "use", "test2"},
				{"config", "save"},
				{"config", "user", "-e", "some@email.net"},
			},
			wantStrOut: [][]string{
				{"Switching to profile..  ", "test2"},
				{"Saving global config.. not changed", "Saving user config.. not changed"},
				{"some@email.net", "User configuration: set", "Success autosave config"},
			},
			inputs: [][]string{
				{},
				{"somePass", "somePass"},
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
				{"packed_key"},
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
				consoleOutput, _ := s.executeCommand(cmd...)
				// fmt.Println(consoleOutput)
				// require.NoError(t, err)
				if i < len(tt.wantStrOut) {
					for _, wantOut := range tt.wantStrOut[i] {
						require.Contains(t, consoleOutput, wantOut, cmd)
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
