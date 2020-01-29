package main

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func prompt(p string) (string, error) {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	t := terminal.NewTerminal(os.Stdin, p)
	return t.ReadPassword(p)
}
