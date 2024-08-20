package main

import (
	"gophKeeper/internal/client/cmd"
	"os"
)

func main() {
	if err := cmd.NewApp().Execute(); err != nil {
		os.Exit(1)
	}
}
