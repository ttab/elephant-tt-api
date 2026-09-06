# Changelog

Everything from v0.5.0 onwards is documented here; earlier releases are not
reconstructed. The entries are derived from the release tags, and the linked
pull requests hold the detail.

## [v0.7.0] - Unreleased

**Build (Go 1.27):** the `go` directive moves from 1.25.7 to **1.27.1**. A
module's directive is a floor for everything that compiles it, so every
consumer of this module needs a Go 1.27.1 toolchain: with `GOTOOLCHAIN=auto`
the go command downloads one, with `GOTOOLCHAIN=local` or no network the build
fails outright. Move your own module to 1.27 before taking this release, and
check that CI does not pin an older toolchain — a workflow that reads
`go-version-file: go.mod` follows the bump on its own, one with a hardcoded
`go-version` does not.

**New protocol (Connect):** every service now ships Connect clients and
handlers alongside the Twirp ones, in a `<package>connect` subpackage —
`baboon/baboonconnect`, `eidos/eidosconnect`,
`everysport/everysportconnect`, `genai/genaiconnect`, `ntb/ntbconnect` and
`wires/wiresconnect`. `New<Service>ServiceClient(httpClient, baseURL)`
returns the same plain service interface `New<Service>ProtobufClient`
returns, so switching a Go client is one constructor call and nothing
downstream changes; `New<Service>ServiceHandler(svc, opts...)` takes an
implementation of that interface and returns the mount path and the handler.
That is `NewAssetsServiceClient` and `NewPrintServiceClient` for `baboon`,
`NewTaggerServiceClient` for `eidos`, `NewManageServiceClient` for
`everysport`, `NewGenerateServiceClient` for `genai`,
`NewMetadataServiceClient`, `NewMediaServiceClient` and
`NewNynorskServiceClient` for `ntb`, and `NewRssFeedServiceClient` for
`wires`. connect-go's own `New<Service>Client` and `New<Service>Handler`,
which speak in `*connect.Request[T]`, are generated too — the `Service` infix
is what tells the two apart. Nothing about the Twirp clients, the Twirp
server interfaces or the `/twirp/` paths changed.

The Connect paths are the standard `/<package>.<Service>/<Method>`, with no
prefix: `POST /ttab.ntb.Nynorsk/Translate` against
`POST /twirp/ttab.ntb.Nynorsk/Translate`. They do not overlap, so both
protocols are served by one server. gRPC and gRPC-Web are served on the
Connect paths as well, selected by content type, but **inside the cluster
only**: the fleet's ingress speaks HTTP/1.1 to its targets, so a gRPC client
outside the cluster cannot reach these services and is not meant to. Use
Connect from outside, gRPC service to service if you want it. An ingress rule
that routes on `/twirp/` needs a sibling rule before a service can serve
Connect. This module only declares the API — whether an environment answers on
the Connect paths is decided by each service's own release.

**Behaviour change (Connect error bodies):** a Connect error body is
`{"code":…,"message":…,"details":[…]}` where Twirp's is
`{"code":…,"msg":…,"meta":{…}}`. Connect has no free-form meta map, so the
key/value metadata travels as an `elephantine.rpc.ErrorMeta` error detail,
which Go callers read with `rpc.Meta(err)` from `elephantine/rpc` and other
clients read with `findDetails`. The codes and the messages are identical on
both stacks. The HTTP status is identical except for three codes: `canceled`
is 499 rather than 408, `deadline_exceeded` is 504 rather than 408, and
`failed_precondition` is **400** where Twirp sends 412. Anything keyed on 412
has to read the code from the body instead.

**Behaviour change (Connect JSON field names):** a Connect success body spells
its fields differently from a Twirp one. Twirp marshals with `protojson` and
`UseProtoNames: true`, so a JSON response carries the names the `.proto`
declares (`document_uuid`, `created_at`); Connect's JSON codec is `protojson`
with its default options, which spell the same fields in lowerCamelCase
(`documentUuid`, `createdAt`). Nothing else about the encoding differs: both
omit unpopulated fields, both render an enum as its name and a timestamp as an
RFC 3339 string, and requests are unaffected because `protojson` unmarshalling
accepts both spellings on both stacks. This reaches one kind of caller: one
that reads JSON responses by hand, with `fetch` or `curl`, and that changes
only the path prefix — it will read `undefined` for every multi-word field.
The generated clients, Go and TypeScript alike, parse into generated types and
see nothing. The difference is deliberate and is not papered over with a
`UseProtoNames` codec: every Connect runtime and proxy assumes the standard
encoding, and a service that deviates from it is one whose clients cannot be
generated from its `.proto` alone.

**Behaviour change (baboon descriptor):** `ttab.baboon`'s `go_package` option
declared `github.com/ttab/baboon/rpc/baboon`, a path that does not exist. It
is corrected to `github.com/ttab/elephant-tt-api/baboon`, which is where the
generated code has always lived and how every consumer already imports it, so
no Go import path changes. It is a change to the descriptor the service
embeds and serves through reflection, and it is the only field that differs
in any descriptor in this release — a consumer that generates its own client
from the descriptor rather than from `service.proto` gets the corrected
package hint.

**Removed (OpenAPI):** the OpenAPI 3 specifications under `docs/` are gone —
`baboon-openapi.json`, `eidos-openapi.json`, `everysport-openapi.json`,
`genai-openapi.json`, `ntb-openapi.json` and `wires-openapi.json`. They
described the Twirp paths and Twirp's error schema only, nobody generated a
client from them, and the generator that wrote them cannot run under buf. The
`.proto` files are the declaration a non-Go consumer generates from. With
nothing left to stamp a version into, a release is a plain git tag: there is
no `rpc:release` target and no "bump to vX.Y.Z" commit any more.

**Build change (generation):** the artifacts are generated with buf and
plugins pinned in `ttab/mage`, not with protoc in the `elephant-twirptools`
Docker image, so regenerating needs no Docker and installs nothing. The mage
targets are renamed to match: `mage rpc:generate` and `mage rpc:stub`,
replacing the `twirp:` ones, with `rpc.Twirp = true` in the magefile keeping
Twirp generation on. buf compiles only what is inside the repository, so
`newsdoc/newsdoc.proto` is now vendored into `rpc/vendor/` — where protoc was
handed the `elephant-api` module directory as a `--proto_path` — and
`mage rpc:vendorProto github.com/ttab/elephant-api newsdoc/newsdoc.proto`
refreshes it. The vendored copy keeps its original path, so the `import`
lines in `baboon`, `eidos` and `genai` are unchanged, and it is compiled but
never generated for: the Go code those services import still comes from
`elephant-api`. Generation runs the plugins under a Go toolchain `ttab/mage`
pins, not the one that happens to be on the machine, so the output no longer
depends on who regenerated it; it does need network access, since every plugin
is resolved from the module proxy on each run, and a `GOPROXY=off` build fails
even with a warm module cache. The `ttab/mage` pin is the released v0.13.1,
which pins `protoc-gen-elephant-rpc` to elephantine v0.29.0, so the generated
code in this release is reproducible from released generators.

Changes:

- The module requires `connectrpc.com/connect` v1.20.0. It still does not
  depend on `elephantine`: the generated code imports only connect and the
  message packages, and the error helpers, header propagation and interceptors
  live in `elephantine/rpc`. (#41)
- The Go code is generated by protoc-gen-go v1.36.12, a newer version than the
  Docker image pinned, which rewrites every `.pb.go`: the embedded descriptor
  becomes a string constant rather than a byte slice, `unsafe` is imported,
  and the header records `protoc (unknown)` because buf reports no protoc
  version. The `service.twirp.go` files change only in the gzip encoding of
  their descriptor blob. The compiled descriptors are byte-identical to the
  ones the image produced apart from `baboon`'s `go_package`, so no message,
  field or method changed. (#41)
- Dependency upgrades: `google.golang.org/protobuf` to v1.36.12, matching the
  protoc-gen-go that writes the generated headers, and `elephant-api` to
  v0.24.2, which leaves `newsdoc/newsdoc.proto` and therefore every embedded
  descriptor unchanged. (#41)
- Each `<package>connect` package has a test that asserts the adapters still
  satisfy the plain service interfaces and still mount on the unprefixed
  paths, so a regeneration that renames or drops an adapter fails rather than
  compiling. (#41)
- `README.md` now covers both protocols: what each service package holds, how
  to build a Connect client, how the two error bodies and the two JSON field
  spellings differ and where the metadata went, which protocols are reachable
  from outside the cluster, how generation works under buf, and what has to be
  tagged before this module can be. (#41)

## [v0.6.1] - 2026-08-31

**Breaking (everysport):** `ImportSettings.merge_events` is removed, and
field number 4 and the name are reserved. The merge strategy now names the
shape of the documents that come out on its own, and `merge_setting` takes a
third value to say so: `"NoMerge"` is one event per calendar entry,
`"CategoryMerge"` one event per category and day, and `"CompetitionMerge"`
one event per competition and day. A client that set `merge_events` has to
pick the corresponding `merge_setting` instead.

**Breaking (everysport):** nothing in `ImportSettings` is inherited any
longer. A competition is imported with the settings on its own row, and a
sport's settings are the defaults its new competitions are created with, so
`merge_setting` and `sport` are both required when enabling an entry rather
than falling back to the sport's values when left empty. An entry enabled
with no topic names of its own publishes events linked to the generic
`"Sport"` topic and nothing else.

Changes:

- `Manage.ApplyDefaults` writes the defaults of a sport or a category to the
  rows below it, selecting which fields (`apply_newsvalue`,
  `apply_merge_setting`, `apply_sport`) and which competitions
  (`competition_states`) it touches. With `dry_run` set it writes nothing and
  reports what would change, so the preview shown before an apply comes from
  the same call that performs it. The response lists one `DefaultsChange` per
  affected row with the old and new value of every field that changes, plus
  counts of the categories, competitions and enabled competitions involved.
- `Manage.GetCategory` returns a single competition category, matching the
  `GetSport` and `GetCompetition` pair.
- Competition categories now carry import defaults. The new `ImportDefaults`
  message — `newsvalue`, `merge_setting` and `sport`, where the zero value
  means "unset, use the sport's" — appears as `CategoryConfig.defaults` and
  is replaced in full through the new `UpdateCategoryRequest.defaults` field.
  It is deliberately not an `ImportSettings`: a category imports nothing, so
  it has no state, and "unset" is a real value for a template where it is not
  one for a row that decides what a document says.
- `UpdateCategoryResponse` gained a `warnings` field, so saving a category's
  defaults reports the same `ConfigWarning` list the sport and competition
  updates do.

## [v0.6.0] - 2026-08-17

- `ttab.everysport`'s `Manage` service grew the whole import configuration
  surface, which was previously only `ManualImport`. Reads:
  `ListSports`, `ListCategories`, `ListCompetitions`, `GetSport`,
  `GetCompetition` and `ExportConfig`, which returns the configuration as
  CSV. Writes: `UpdateSport`, `UpdateCompetition` and `UpdateCategory`, each
  replacing the editable settings in full and returning a `ConfigWarning`
  list for problems that don't stop a save. Review:
  `ListPendingImports` lists the competitions whose calendar entries were
  skipped, with the date range and the count still missing, and
  `ReviewCompetitions` decides a batch of them at once. Entries the service
  discovers start in `STATE_PENDING` and import nothing until reviewed; the
  states are declared by the new `ConfigState` enum, and the settings, the
  sport, category and competition rows by `ImportSettings`, `SportConfig`,
  `CategoryConfig` and `CompetitionConfig`. (#40)
- `ManualImportRequest.competitions` re-imports several competitions by
  static ID in one call. It combines with the existing single `competition`
  field rather than replacing it — an event matches if its static ID is in
  either. (#40)

## [v0.5.2] - 2026-05-29

- `ttab.ntb`'s `Media.Search` takes five more filters on `SearchRequest`:
  `editorial_topics`, `distributor_names`, `events`, `subject_terms` and
  `persons`, each a repeated string. (#35)

## [v0.5.1] - 2026-04-23

- `ttab.ntb`'s `Metadata` service gained `ListAssignments`, which returns
  every assignment rather than resolving one by commission code the way
  `GetAssignment` does. Each entry is the new `Assignment` message, which
  carries `assignee_uuid`, `assignee_name` and `commission_code` alongside
  the planning and assignment UUIDs that `GetAssignmentResponse` returns.
  (#30)
- Dependency upgrades: elephant-api to v0.22.2.

## [v0.5.0] - 2026-03-27

- `ttab.ntb` gained the `Nynorsk` service, which translates documents between
  Scandinavian languages through Nynorskroboten. `Translate` takes a document
  as plain text, as an ordered list of fragments or as a named map, with the
  source and target languages as ISO 639-1 codes, the document format
  (`"html"`, `"txt"`, `"docx"`, `"json"`), a `post_edit` flag for AI-assisted
  post-editing, and language form preferences either by template name or
  individually through the `prefs` map. `GetLanguagePrefs`,
  `GetDocumentFormats` and `GetPreferenceTemplates` enumerate what those
  fields accept. `LangParameterTranslate` controls what happens to text
  inside HTML `lang` elements. The four methods are also documented in
  `ntb/docs/`. (#29)
