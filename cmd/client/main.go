package main

import (
	"gophKeeper/internal/client/cmd"
	"os"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	err := cmd.NewApp(cmd.BuildMetadata{
		Version: buildVersion,
		Date:    buildDate,
		Commit:  buildCommit,
	}).Execute()
	if err != nil {
		os.Exit(1)
	}
}
