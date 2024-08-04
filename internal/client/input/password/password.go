package password

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

func readPassword() (pw []byte, e error) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		pw, e = term.ReadPassword(int(os.Stdin.Fd()))
		return
	}

	var b [1]byte
	for {
		n, err := os.Stdin.Read(b[:])
		// terminal.ReadPassword discards any '\r', so do the same
		if n > 0 && b[0] != '\r' {
			if b[0] == '\n' {
				return pw, nil
			}
			pw = append(pw, b[0])
			if len(pw) > 1024 {
				err = errors.New("password too long")
			}
		}
		if err != nil {
			// terminal.ReadPassword accepts EOF-terminated passwords
			// if non-empty, so do the same
			if err == io.EOF && len(pw) > 0 {
				err = nil
			}
			return
		}
	}
}

func GetRawPass(confirm bool) (pass string, err error) {
	var b []byte
	if confirm {
		fmt.Print("Please enter new password: ")
	} else {
		fmt.Print("Please enter password: ")
	}

	b, err = readPassword()
	if err == nil {
		if confirm {
			try := 0
			fmt.Println()
			fmt.Print("Please confirm you password: ")
		RepeatPass:
			b2, err2 := readPassword()
			if err2 != nil || string(b) != string(b2) {
				err = errors.New("password confirm error")
				fmt.Println("password confirm error")
				try++
				if try > 3 {
					return
				}
				goto RepeatPass
			}
		}
		pass = string(b)
	}
	return
}
