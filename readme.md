
# go-rest-api-cli

A small, cross-platform Go CLI for making REST API calls.

- ✅ Works on **Windows** and **Linux**
- ✅ **No external dependencies** – only Go standard library → great for **air-gapped** environments
- ✅ JSON payloads from file and inline, with **merge & override**
- ✅ **Profiles** for base URL, default headers, and auth
- ✅ **Hashing** binary files (MD5 / SHA-1 / SHA-256) and injecting into JSON
- ✅ Flexible output options: pretty JSON, raw, JSON-only, write to file
- ✅ Basic **retry logic**
- ✅ Clean architecture with **Command**, **Factory**, **Strategy** patterns
- ✅ Unit tests for core modules

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Project layout

```text
.
├── go.mod
├── main.go
└── internal
    ├── auth
    │   └── auth.go
    ├── command
    │   ├── call.go
    │   ├── command.go
    │   ├── headers.go
    │   ├── help.go
    │   ├── inspect.go
    │   └── profile.go
    ├── config
    │   └── config.go
    ├── httpclient
    │   └── factory.go
    └── payload
        └── json.go
```

(Plus `*_test.go` files next to some modules for unit tests.)

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Module overview

### `main.go`

- Entry point for the CLI.
- Creates a `command.Registry`, registers:
  - `call`
  - `profile`
  - `inspect`
  - `help`
- Dispatches based on `os.Args[1]` (the subcommand name).
- Handles global errors and exit codes.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### `internal/command`

Implements the **Command pattern** for subcommands.

#### `command.go`

- Defines the `Command` interface:

  ```go
  type Command interface {
      Name() string
      Description() string
      Run(args []string) error
  }
  ```

- `Registry` stores and retrieves commands by name.

#### `headers.go`

- Implements `HeaderFlag`, a `flag.Value` for repeated `--header` flags.
- Parses strings like `"Key: Value"` into a `map[string]string`.

#### `call.go`

The core **API call** command.

Responsibilities:

- Parse CLI flags:
  - HTTP method, URL
  - Profiles (`--profile`)
  - JSON payload (`--json-file`, `--data`)
  - Hashing (`--hash-file`, `--hash-algo`, `--hash-field`, `--hash-prefix-0x`)
  - Headers (`--header`)
  - Auth (`--auth`, `--user`, `--pass`, `--token`)
  - Output modes (`--pretty`, `--raw`, `--json-only`, `--out`)
  - Retry (`--retries`, `--retry-delay`)
- Load profile from config (if used).
- Merge:
  - Profile base URL + relative `--url`
  - Profile headers + CLI `--header`
  - Profile auth + CLI auth overrides
- Load & merge JSON from file + inline.
- Compute file hash (optional) and inject into JSON.
- Build request/client via `httpclient.Factory`.
- Perform request with retry logic.
- Format and print/save response.

#### `profile.go`

Manages saved profiles.

Commands:

- `profile add` – create/update a profile with:
  - `--name NAME`
  - `--base-url URL`
  - `--auth none|basic|bearer`
  - optional `--user`, `--pass`, `--token`
  - `--header "Key: Value"` (repeatable)
- `profile list` – list all profiles (basic info).
- `profile remove` – delete a profile by name.

Uses `internal/config` for persistence.

#### `inspect.go`

Pretty printing of stored profiles.

Commands:

- `inspect profiles` – show all profiles with details.
- `inspect profile --name NAME` – show a single profile in detail (base URL, auth, headers).

#### `help.go`

- Prints:
  - Tool description
  - Available commands
  - Quick examples

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### `internal/config`

Handles storage of **profiles** in a config file.

#### `config.go`

- Types:
  - `Profile` – base URL, headers, auth type, user/pass/token.
  - `Config` – root struct with `Profiles map[string]Profile`.
- Determines config file path:
  - Uses `os.UserConfigDir()` if possible.
  - Falls back to `~/.go-rest-api-cli/config.json`.
- `Load()`:
  - Returns a `Config` instance.
  - If file does not exist, returns an empty `Config` with initialized map.
- `Save()`:
  - Creates parent folders if needed.
  - Writes `config.json` with pretty JSON.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### `internal/auth`

Implements **auth strategies** (Strategy pattern).

#### `auth.go`

- `Strategy` interface with `Apply(req *http.Request)`.
- Implementations:
  - `NoAuth` – no changes to request.
  - `Basic` – sets `Authorization: Basic ...` using `req.SetBasicAuth`.
  - `Bearer` – sets `Authorization: Bearer <token>`.

Used by `call` command (and indirectly by `httpclient.Factory`).

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### `internal/httpclient`

Factory for building HTTP clients/requests.

#### `factory.go`

- `Config` struct:
  - Method, URL, Headers, Body
  - Timeout
  - Auth strategy
  - `SkipTLSVerify` flag
- `Factory.Build(cfg Config)`:
  - Creates `*http.Request` with method, URL, body.
  - Sets headers.
  - Applies auth strategy.
  - Builds `*http.Client` with given timeout.
  - Attaches `*http.Transport` with optional `InsecureSkipVerify`.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### `internal/payload`

Helpers for JSON payload and hashing.

#### `json.go`

- `LoadJSONFile(path)` – loads JSON file into `map[string]interface{}`.
- `ParseJSONInline(string)` – parses inline JSON into `map[string]interface{}`.
- `Merge(fileMap, inlineMap)` – merges two maps, where `inlineMap` overrides keys from `fileMap`.
- `NormalizeHashAlgo(algo)` – normalizes names like `"SHA-256"` → `"sha256"`.
- `ComputeFileHash(path, algo)`:
  - Supports `md5`, `sha-1`, `sha-256`.
  - Returns hex string hash.

`call` uses this module to build the final request body and compute hashes.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Build & install

From project root:

### Linux binary

```bash
GOOS=linux GOARCH=amd64 go build -o go-rest-api-cli .
```

### Windows binary

**From any OS** (cross-compile):

```bash
GOOS=windows GOARCH=amd64 go build -o go-rest-api-cli.exe .
```

Then copy the binary to your target machine.

> **Note:**  
> - Run `./go-rest-api-cli` in Linux/macOS shells.  
> - Run `go-rest-api-cli.exe` (or `.\go-rest-api-cli.exe`) in Windows `cmd` / PowerShell.  
> - Don’t run a Windows `.exe` inside WSL bash – you’ll get an “Exec format error”.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Commands overview

General syntax:

```bash
go-rest-api-cli <command> [flags...]
```

Available commands:

- `call` – execute a REST API call
- `profile` – manage saved profiles
- `inspect` – inspect stored profiles
- `help` – show help and examples

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## `call` command

### Key flags

- `--method` – HTTP method (default: `GET`)
- `--url` – request URL (**required**)
  - absolute: `https://api.example.com/v1/users`
  - or relative: `/v1/users` when using `--profile`
- `--profile` – profile name (base URL, headers, auth)
- `--header "Key: Value"` – extra headers (repeatable)
- `--json-file` – JSON file path for payload (merged base)
- `--data` – inline JSON to merge/override file payload
- `--auth` – `none|basic|bearer`
- `--user`, `--pass`, `--token` – auth parameters
- `--timeout` – in seconds (default `30`)
- `--insecure` – skip TLS verification (lab only)
- `--pretty` – pretty-print JSON responses
- `--raw` – print only body
- `--json-only` – print only JSON body
- `--out` – write response body to file
- `--retries` – number of retries (network/5xx)
- `--retry-delay` – delay between retries (seconds)
- `--hash-file` – path to file to hash
- `--hash-algo` – `md5|sha-1|sha-256` (default `sha-256`)
- `--hash-field` – JSON field name for hash (default `file_hash`)
- `--hash-prefix-0x` – if set, prefix hash with `0x`

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Examples

Below: for each scenario you get **Bash**, **PowerShell**, and **cmd** examples.

### 1. Simple GET call (Agify API)

#### Bash (Linux/macOS)

```bash
./go-rest-api-cli call \
  --method GET \
  --url "https://api.agify.io/?name=meelad" \
  --pretty
```

#### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method GET `
  --url "https://api.agify.io/?name=meelad" `
  --pretty
```

#### cmd.exe

```bat
go-rest-api-cli.exe call --method GET --url "https://api.agify.io/?name=meelad" --pretty
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### 2. POST with inline JSON only

We’ll POST to `https://api.restful-api.dev/objects`.

Inline JSON:

```json
{"name":"Dev Object","data":{"env":"dev","owner":"ido"}}
```

#### Bash

```bash
./go-rest-api-cli call \
  --method POST \
  --url "https://api.restful-api.dev/objects" \
  --data '{"name":"Dev Object","data":{"env":"dev","owner":"ido"}}' \
  --pretty
```

#### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method POST `
  --url "https://api.restful-api.dev/objects" `
  --data '{"name":"Dev Object","data":{"env":"dev","owner":"ido"}}' `
  --pretty
```

#### cmd.exe

> In `cmd.exe` you must escape double quotes inside JSON:

```bat
go-rest-api-cli.exe call ^
  --method POST ^
  --url "https://api.restful-api.dev/objects" ^
  --data "{\"name\":\"Dev Object\",\"data\":{\"env\":\"dev\",\"owner\":\"ido\"}}" ^
  --pretty
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### 3. POST with JSON file + inline override

`payload.json`:

```json
{
  "name": "Base Object",
  "data": {
    "env": "prod",
    "version": 1
  }
}
```

We override `name` and add extra data via inline JSON.

#### Bash

```bash
./go-rest-api-cli call \
  --method POST \
  --url "https://api.restful-api.dev/objects" \
  --json-file "payload.json" \
  --data '{"name":"Overridden Name","data":{"extra":"from-inline"}}' \
  --pretty
```

#### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method POST `
  --url "https://api.restful-api.dev/objects" `
  --json-file "payload.json" `
  --data '{"name":"Overridden Name","data":{"extra":"from-inline"}}' `
  --pretty
```

#### cmd.exe

```bat
go-rest-api-cli.exe call ^
  --method POST ^
  --url "https://api.restful-api.dev/objects" ^
  --json-file "payload.json" ^
  --data "{\"name\":\"Overridden Name\",\"data\":{\"extra\":\"from-inline\"}}" ^
  --pretty
```

Effective payload sent:

```json
{
  "name": "Overridden Name",
  "data": {
    "env": "prod",
    "version": 1,
    "extra": "from-inline"
  }
}
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### 4. Profiles: add, list, inspect, use

#### 4.1 Create a profile

Profile name: `restful`  
Base URL: `https://api.restful-api.dev`

##### Bash

```bash
./go-rest-api-cli profile add \
  --name restful \
  --base-url https://api.restful-api.dev \
  --auth none \
  --header "X-Env: dev"
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe profile add `
  --name restful `
  --base-url https://api.restful-api.dev `
  --auth none `
  --header "X-Env: dev"
```

##### cmd.exe

```bat
go-rest-api-cli.exe profile add ^
  --name restful ^
  --base-url https://api.restful-api.dev ^
  --auth none ^
  --header "X-Env: dev"
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 4.2 List profiles

##### Bash

```bash
./go-rest-api-cli profile list
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe profile list
```

##### cmd.exe

```bat
go-rest-api-cli.exe profile list
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 4.3 Inspect all profiles

##### Bash

```bash
./go-rest-api-cli inspect profiles
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe inspect profiles
```

##### cmd.exe

```bat
go-rest-api-cli.exe inspect profiles
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 4.4 Inspect a single profile

##### Bash

```bash
./go-rest-api-cli inspect profile --name restful
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe inspect profile --name restful
```

##### cmd.exe

```bat
go-rest-api-cli.exe inspect profile --name restful
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 4.5 Use profile with relative URL

Call `GET https://api.restful-api.dev/objects/1` via profile.

##### Bash

```bash
./go-rest-api-cli call \
  --profile restful \
  --method GET \
  --url "/objects/1" \
  --pretty
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --profile restful `
  --method GET `
  --url "/objects/1" `
  --pretty
```

##### cmd.exe

```bat
go-rest-api-cli.exe call ^
  --profile restful ^
  --method GET ^
  --url "/objects/1" ^
  --pretty
```

The tool combines:

```text
base-url: https://api.restful-api.dev
url:      /objects/1
→ https://api.restful-api.dev/objects/1
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### 5. Hash a file and inject into JSON payload

We will:

1. Compute a hash for a file.
2. Inject that hash into the JSON payload under a custom field.
3. Update the JSON file on disk.
4. Send the resulting payload.

Assume:

- Binary file: `build/artifact.bin`
- JSON payload: `payload.json`

`payload.json` (before):

```json
{
  "name": "artifact upload",
  "build": 42
}
```

#### 5.1 SHA-256 hash into `checksum` field (with no prefix)

##### Bash

```bash
./go-rest-api-cli call \
  --method POST \
  --url "https://api.restful-api.dev/objects" \
  --json-file "payload.json" \
  --hash-file "build/artifact.bin" \
  --hash-algo "sha-256" \
  --hash-field "checksum" \
  --pretty
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method POST `
  --url "https://api.restful-api.dev/objects" `
  --json-file "payload.json" `
  --hash-file "build/artifact.bin" `
  --hash-algo "sha-256" `
  --hash-field "checksum" `
  --pretty
```

##### cmd.exe

```bat
go-rest-api-cli.exe call ^
  --method POST ^
  --url "https://api.restful-api.dev/objects" ^
  --json-file "payload.json" ^
  --hash-file "build/artifact.bin" ^
  --hash-algo "sha-256" ^
  --hash-field "checksum" ^
  --pretty
```

What happens:

- Computes `sha-256(build/artifact.bin)` → e.g. `"b94d27b9..."`
- Prints:

  ```text
  Computed hash (sha-256) for build/artifact.bin: b94d27b9...
  Updated JSON file payload.json with field checksum
  ```

- `payload.json` is updated to:

  ```json
  {
    "name": "artifact upload",
    "build": 42,
    "checksum": "b94d27b9..."
  }
  ```

- This updated JSON is sent as the request body.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 5.2 MD5 with `0x` prefix into `md5_sum`

##### Bash

```bash
./go-rest-api-cli call \
  --method POST \
  --url "https://api.restful-api.dev/objects" \
  --json-file "payload.json" \
  --hash-file "build/artifact.bin" \
  --hash-algo "md5" \
  --hash-field "md5_sum" \
  --hash-prefix-0x \
  --pretty
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method POST `
  --url "https://api.restful-api.dev/objects" `
  --json-file "payload.json" `
  --hash-file "build/artifact.bin" `
  --hash-algo "md5" `
  --hash-field "md5_sum" `
  --hash-prefix-0x `
  --pretty
```

##### cmd.exe

```bat
go-rest-api-cli.exe call ^
  --method POST ^
  --url "https://api.restful-api.dev/objects" ^
  --json-file "payload.json" ^
  --hash-file "build/artifact.bin" ^
  --hash-algo "md5" ^
  --hash-field "md5_sum" ^
  --hash-prefix-0x ^
  --pretty
```

Now:

- Hash looks like `"0x5eb63bbb..."`.
- JSON is updated:

  ```json
  {
    "name": "artifact upload",
    "build": 42,
    "checksum": "b94d27b9...",
    "md5_sum": "0x5eb63bbb..."
  }
  ```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### 6. Output modes and saving response

Use a simple GET and play with output options:

#### 6.1 Default output (status + headers + body)

##### Bash

```bash
./go-rest-api-cli call \
  --method GET \
  --url "https://api.agify.io/?name=meelad"
```

#### 6.2 Pretty JSON only (no status/headers)

##### Bash

```bash
./go-rest-api-cli call \
  --method GET \
  --url "https://api.agify.io/?name=meelad" \
  --json-only \
  --pretty
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method GET `
  --url "https://api.agify.io/?name=meelad" `
  --json-only `
  --pretty
```

##### cmd.exe

```bat
go-rest-api-cli.exe call --method GET --url "https://api.agify.io/?name=meelad" --json-only --pretty
```

#### 6.3 Raw body, saved to file

##### Bash

```bash
./go-rest-api-cli call \
  --method GET \
  --url "https://api/agify.io/?name=meelad" \
  --raw \
  --out "agify_response.json"
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### 7. Retry logic

Example: flaky endpoint with retry attempts.

#### Bash

```bash
./go-rest-api-cli call \
  --method GET \
  --url "https://flaky.example.com/data" \
  --retries 3 \
  --retry-delay 2 \
  --pretty
```

#### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method GET `
  --url "https://flaky.example.com/data" `
  --retries 3 `
  --retry-delay 2 `
  --pretty
```

#### cmd.exe

```bat
go-rest-api-cli.exe call ^
  --method GET ^
  --url "https://flaky.example.com/data" ^
  --retries 3 ^
  --retry-delay 2 ^
  --pretty
```

Behavior:

- Total attempts: `retries + 1` → here: 4 attempts maximum.
- Retries on:
  - Network errors
  - HTTP 5xx status codes
- Sleeps `--retry-delay` seconds between attempts.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Running tests

If you added the `*_test.go` files as described:

### Run all tests

```bash
go test ./...
```

### Verbose output

```bash
go test -v ./...
```

### Single package

```bash
go test ./internal/payload
go test ./internal/auth
go test ./internal/httpclient
go test ./internal/command
```

### With coverage

```bash
go test ./... -cover
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Design patterns used

- **Command pattern**
  - `internal/command.Command` interface + `Registry`.
  - Subcommands: `CallCommand`, `ProfileCommand`, `InspectCommand`, `HelpCommand`.

- **Strategy pattern (Auth)**
  - `internal/auth.Strategy` interface.
  - Concrete strategies: `NoAuth`, `Basic`, `Bearer`.

- **Factory pattern (HTTP)**
  - `internal/httpclient.Factory` builds `*http.Request` and `*http.Client` from a config.

- **Configuration module**
  - `internal/config` isolates config file location & schema (profiles).

- **Separation of concerns**
  - CLI parsing & orchestration → `internal/command`
  - HTTP behavior → `internal/httpclient`
  - Auth logic → `internal/auth`
  - JSON & hashing → `internal/payload`
  - Persistence of profiles → `internal/config`



