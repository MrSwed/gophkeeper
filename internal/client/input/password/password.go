package password

import (
	"errors"
	"fmt"
	errs "gophKeeper/internal/client/errors"
	"io"
	"os"

	"golang.org/x/term"
)

const (
	PromptPassword        = "Please enter password: "
	PromptNewPassword     = "Please enter new password: "
	PromptConfirmPassword = "Please confirm you password: "
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

// GetRawPass get password from user
//
//	if confirm is true, ask for confirm password
//	if prompts is not empty, use it as prompts for password and confirm password
func GetRawPass(confirm bool, prompts ...string) (pass string, err error) {
	var b []byte
	if len(prompts) == 0 {
		if confirm {
			prompts = append(prompts, PromptPassword)
		} else {
			prompts = append(prompts, PromptNewPassword)
		}
		if confirm && len(prompts) == 1 {
			prompts = append(prompts, PromptConfirmPassword)
		}
	}
	fmt.Print(prompts[0])
	b, err = readPassword()
	fmt.Println()
	if err == nil {
		if confirm {
			try := 0
			fmt.Print(prompts[1])
		RepeatPass:
			b2, err2 := readPassword()
			fmt.Println()
			if err2 != nil || string(b) != string(b2) {
				err = errors.Join(err2, errs.ErrPasswordConfirm)
				fmt.Println("\n", err.Error())
				try++
				if try > 3 {
					return
				}
				goto RepeatPass
			}
		}
		err = nil
		pass = string(b)
	}
	return
}
