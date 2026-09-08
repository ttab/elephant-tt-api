package baboonconnect_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/ttab/elephant-tt-api/baboon"
	"github.com/ttab/elephant-tt-api/baboon/baboonconnect"
)

// The Connect clients are drop-in replacements for the Twirp ones, so each
// one implements the plain service interface. A regeneration that renames or
// drops an adapter fails to compile here.
var (
	_ baboon.Assets = baboonconnect.NewAssetsServiceClient(nil, "")
	_ baboon.Print  = baboonconnect.NewPrintServiceClient(nil, "")
)

// The handler adapters take the plain implementation and return the mount
// path together with the handler, which is the pair the API server registers.
var (
	_ func(baboon.Assets, ...connect.HandlerOption) (string, http.Handler) = baboonconnect.NewAssetsServiceHandler
	_ func(baboon.Print, ...connect.HandlerOption) (string, http.Handler)  = baboonconnect.NewPrintServiceHandler
)

// TestHandlerPaths pins the Connect mount paths. They carry no "/twirp"
// prefix, so the two protocols coexist on one server, and an ingress rule
// written against them is only correct for as long as this holds.
func TestHandlerPaths(t *testing.T) {
	cases := []struct {
		Service string
		Handler func() (string, http.Handler)
	}{
		{
			Service: baboonconnect.AssetsName,
			Handler: func() (string, http.Handler) {
				return baboonconnect.NewAssetsServiceHandler(nil)
			},
		},
		{
			Service: baboonconnect.PrintName,
			Handler: func() (string, http.Handler) {
				return baboonconnect.NewPrintServiceHandler(nil)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Service, func(t *testing.T) {
			path, handler := c.Handler()

			want := "/" + c.Service + "/"
			if path != want {
				t.Errorf("got the mount path %q, wanted %q", path, want)
			}

			if handler == nil {
				t.Error("got a nil handler")
			}
		})
	}
}
