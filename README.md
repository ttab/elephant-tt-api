# TT API declarations for the Elephant

Protobuf API declarations for the TT specific parts of Elephant. Each service
is defined in a `service.proto` file and shipped with generated Go code for two
protocols: [Connect](https://connectrpc.com/), which a service also answers as
gRPC and gRPC-Web to callers inside the cluster, and
[Twirp](https://github.com/twitchtv/twirp), which is what the platform served
before Connect and is still served everywhere.

## The APIs

| Package | Services | Purpose | Declaration |
| --- | --- | --- | --- |
| `ttab.baboon` | `Assets`, `Print` | Print rendering, font and ICC profile management | [proto](baboon/service.proto) |
| `ttab.eidos` | `Tagger` | NLP document tagging and entity detection | [proto](eidos/service.proto) |
| `ttab.everysport` | `Manage` | Sports calendar event import | [proto](everysport/service.proto) |
| `ttab.genai` | `Generate` | GenAI document creation | [proto](genai/service.proto) |
| `ttab.ntb` | `Metadata`, `Media`, `Nynorsk` | NTB metadata, media archive search and Nynorsk translation | [proto](ntb/service.proto) |
| `ttab.wires` | `RssFeed` | RSS feed management | [proto](wires/service.proto) |

`baboon`, `eidos` and `genai` carry NewsDoc documents and blocks in their
requests and responses, and import `newsdoc/newsdoc.proto` from
[`elephant-api`](https://github.com/ttab/elephant-api) for them.

## Using the APIs

### Go

```bash
go get github.com/ttab/elephant-tt-api@latest
```

The module's `go` directive is 1.27.1, which is a floor for anything that
compiles it: a consumer on an older toolchain downloads one under
`GOTOOLCHAIN=auto` and fails under `GOTOOLCHAIN=local`.

Each service package holds the messages and the plain service interface, which
is the contract both protocols are expressed in:

```go
Translate(ctx context.Context, req *ntb.TranslateRequest) (*ntb.TranslateResponse, error)
```

#### Connect clients

The Connect clients live in a `<package>connect` subpackage and return that
same plain interface, so they are drop-in replacements for the Twirp clients:

```go
import (
	"github.com/ttab/elephant-tt-api/ntb"
	"github.com/ttab/elephant-tt-api/ntb/ntbconnect"
)

// client is an *http.Client that carries the bearer token, usually one built
// by oauth2.NewClient.
var nynorsk ntb.Nynorsk = ntbconnect.NewNynorskServiceClient(
	client, baseURL)
```

`New<Service>ServiceClient` takes the base URL of the server, not a per-service
path. The constructors, one pair per service:

| Package | Clients and handlers |
| --- | --- |
| `baboon/baboonconnect` | `Assets`, `Print` |
| `eidos/eidosconnect` | `Tagger` |
| `everysport/everysportconnect` | `Manage` |
| `genai/genaiconnect` | `Generate` |
| `ntb/ntbconnect` | `Metadata`, `Media`, `Nynorsk` |
| `wires/wiresconnect` | `RssFeed` |

`New<Service>ServiceHandler(svc, opts...)` is the server side. It takes an
implementation of the plain interface and returns the mount path together with
the handler, which is the pair `elephantine`'s API server registers.

The same packages also carry connect-go's own generated `New<Service>Client`
and `New<Service>Handler`, which speak in `*connect.Request[T]` and
`*connect.Response[T]`. The `Service` infix is what distinguishes the plain
adapters from them. Use the adapters unless you need per-call access to
headers.

#### Twirp clients

`New<Service>ProtobufClient` and `New<Service>JSONClient` are unchanged and
still generated. Nothing about them, or about the `/twirp/` paths, has changed.

### Other languages

Both protocols speak JSON over HTTP `POST` and are served side by side, on
different paths:

| Protocol | Path | Content types |
| --- | --- | --- |
| Connect | `/ttab.ntb.Nynorsk/Translate` | `application/json`, `application/proto` |
| Twirp | `/twirp/ttab.ntb.Nynorsk/Translate` | `application/json`, `application/protobuf` |

The Connect paths carry no prefix, so the two families never overlap and one
server mounts both. Connect clients send a `Connect-Protocol-Version: 1`
header, and `Connect-Timeout-Ms` sets a deadline; the servers do not require
either, so a plain `curl` or `fetch` works.

A service answers gRPC and gRPC-Web on the same Connect paths, selected by
content type, but only to callers inside the cluster: the ingress speaks
HTTP/1.1 to its targets, so neither protocol is reachable from outside and
neither is offered to customers. A browser client uses Connect — gRPC-Web
carries its errors in trailers, which a browser cannot read cross-origin, and
the fleet deliberately does not expose them.

There is no OpenAPI specification: the `.proto` file is the declaration, and a
non-Go consumer generates its client from it with its language's Connect or
protobuf tooling.

#### JSON field names differ between the two protocols

A Twirp response spells its fields the way the `.proto` declares them, and a
Connect response spells them in lowerCamelCase:

```json
{"assignee_uuid": "…", "commission_code": "…"}   // Twirp
{"assigneeUuid": "…", "commissionCode": "…"}     // Connect
```

Twirp marshals with `protojson` and `UseProtoNames: true`; Connect's JSON codec
is `protojson` with its default options. Nothing else about the encoding
differs — both omit unpopulated fields, render an enum as its name and a
timestamp as an RFC 3339 string — and requests are unaffected, because
`protojson` unmarshalling accepts both spellings on both stacks.

This only reaches a caller that reads JSON responses by hand. The generated
clients, Go and TypeScript alike, parse into generated types and see nothing.
A `fetch` or `curl` caller that changes only the path prefix reads `undefined`
for every multi-word field, which is the failure to look for.

### Errors

The two protocols share the error codes but not the body. Twirp:

```json
{"code": "not_found", "msg": "no such feed", "meta": {"uuid": "..."}}
```

Connect:

```json
{
  "code": "not_found",
  "message": "no such feed",
  "details": [{"type": "elephantine.rpc.ErrorMeta", "value": "<base64 Any>"}]
}
```

`msg` is `message`, and Connect has no free-form meta map in the body. The
key/value metadata the Twirp errors carry travels as an
`elephantine.rpc.ErrorMeta` error detail instead, which Go callers read with
`rpc.Meta(err)` from [`elephantine/rpc`](https://github.com/ttab/elephantine)
and other clients read with their Connect implementation's `findDetails`. The
message is byte for byte the same on both stacks.

The HTTP status differs for three codes: `canceled` is 499 rather than 408,
`deadline_exceeded` is 504 rather than 408, and `failed_precondition` is
**400** where Twirp sends 412. Read the code from the body rather than the
status. `elephantine`'s [`docs/connect.md`](https://github.com/ttab/elephantine/blob/main/docs/connect.md)
is the fleet reference for the rest of the detail.

This module declares the messages and nothing else: it does not depend on
`elephantine`, so the generated code imports only `connectrpc.com/connect` and
the message packages. Header propagation, error helpers and interceptors come
from `elephantine/rpc`.

## Working in this repo

The Protobuf, Connect and Twirp artifacts are generated through
[mage](https://magefile.org/) targets from
[`ttab/mage`](https://github.com/ttab/mage). The compiler is
[buf](https://buf.build/) and every plugin is pinned there and resolved from
the module proxy at generation time — there is no Docker image and nothing is
installed or taken off `PATH`. mage also pins the Go toolchain the plugins run
under, so the output does not depend on the toolchain the machine happens to
have. What it does depend on is the network: generation queries the proxy on
every run and fails with `GOPROXY=off` even when the module cache is warm. A
generator version moves when `ttab/mage` is bumped. Run all targets from the
repository root.

| Target | Purpose |
| --- | --- |
| `mage rpc:generate` | Regenerate the Go, Connect and Twirp artifacts for every service. |
| `mage rpc:stub <app> <service> <method>` | Scaffold a new service proto. |
| `mage rpc:vendorProto github.com/ttab/elephant-api newsdoc/newsdoc.proto` | Refresh the vendored copy of the NewsDoc declaration. |

Per service directory, generation writes `service.pb.go` (messages),
`service.twirp.go` (the Twirp clients, server, and the plain service
interface), and `<package>connect/service.connect.go` plus
`<package>connect/service.elephant.go` (the Connect clients and handlers, and
the adapters that put them on the plain interface). Nothing but Go is
generated.

buf compiles only what is inside the repository, so `newsdoc/newsdoc.proto` is
vendored into `rpc/vendor/` rather than resolved out of the `elephant-api`
module directory. The vendored file keeps the path it has in the repository it
came from, so the `import` lines in the service protos are unchanged, and it is
compiled but never generated for — its Go code comes from `elephant-api`, which
is where the services import it from. `buf.yaml` is written by the target and
is what makes the vendor directory a module root of its own. Refresh the copy
after bumping the `elephant-api` dependency; the target is idempotent, so a
`git diff` afterwards is the drift.

To change an API, edit its `service.proto`, regenerate, and commit the proto
together with the regenerated files.

## Releasing

A release is a git tag and nothing else — there is no version file and nothing
is stamped with the version. From a clean tree on `main`, with the generated
files up to date (`mage rpc:generate` produces no diff), cut `vX.Y.Z`:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

[`CHANGELOG.md`](CHANGELOG.md) documents what each release changed, and a
change a consumer would want to know about gets its entry in the same commit
as the change itself.

### Generators

The generators this module is built with are released: `ttab/mage` v0.13.1
runs buf and the plugins at pinned versions and pins `protoc-gen-elephant-rpc`
to elephantine v0.29.0, so the generated code in a tag is reproducible from
released code. Keep it that way: a `ttab/mage` requirement that is a
pseudo-version of a branch is a reason not to tag, since the branch can be
force-pushed out from under the tag.

## License

Licensed under MIT.
