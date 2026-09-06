package genaiconnect_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/ttab/elephant-tt-api/genai"
	"github.com/ttab/elephant-tt-api/genai/genaiconnect"
)

// The Connect clients are drop-in replacements for the Twirp ones, so each
// one implements the plain service interface. A regeneration that renames or
// drops an adapter fails to compile here.
var (
	_ genai.Generate = genaiconnect.NewGenerateServiceClient(nil, "")
)

// The handler adapters take the plain implementation and return the mount
// path together with the handler, which is the pair the API server registers.
var (
	_ func(genai.Generate, ...connect.HandlerOption) (string, http.Handler) = genaiconnect.NewGenerateServiceHandler
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
			Service: genaiconnect.GenerateName,
			Handler: func() (string, http.Handler) {
				return genaiconnect.NewGenerateServiceHandler(nil)
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
