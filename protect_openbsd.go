//+build openbsd

package main

import (
	"golang.org/x/sys/unix"
)

func unveil(path string, flags string) {
	unix.Unveil(path, flags)
}

func unveilBlock() {
	unix.UnveilBlock()
}

func plegde(promises string) {
	unix.PledgePromises(promises)

}
