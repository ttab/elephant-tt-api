package everysportconnect_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/ttab/elephant-tt-api/everysport"
	"github.com/ttab/elephant-tt-api/everysport/everysportconnect"
)

// The Connect clients are drop-in replacements for the Twirp ones, so each
// one implements the plain service interface. A regeneration that renames or
// drops an adapter fails to compile here.
var (
	_ everysport.Manage = everysportconnect.NewManageServiceClient(nil, "")
)

// The handler adapters take the plain implementation and return the mount
// path together with the handler, which is the pair the API server registers.
var (
	_ func(everysport.Manage, ...connect.HandlerOption) (string, http.Handler) = everysportconnect.NewManageServiceHandler
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
			Service: everysportconnect.ManageName,
			Handler: func() (string, http.Handler) {
				return everysportconnect.NewManageServiceHandler(nil)
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
