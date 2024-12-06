package main

import (
	"os"

	"golang.org/x/term"
)

// noinspection GoUnresolvedReference
func prompt(p string) (string, error) {
	oldState, err := term.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer term.Restore(0, oldState)

	t := term.NewTerminal(os.Stdin, p)
	return t.ReadPassword(p)
}
