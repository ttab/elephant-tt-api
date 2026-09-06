//go:build mage
// +build mage

package main

import (
	//mage:import rpc
	"github.com/ttab/mage/rpc"
)

func init() {
	// The services are dual stack: the /twirp/ paths are still served, so
	// the Twirp code is generated alongside the Connect code.
	rpc.Twirp = true
}
