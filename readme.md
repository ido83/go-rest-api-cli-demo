# go-rest-api-cli

```text
#                                 _                _       _ _ 
#                                | |              (_)     | (_)
#       __ _  ___   _ __ ___  ___| |_   __ _ _ __  _   ___| |_ 
#      / _` |/ _ \ | '__/ _ \/ __| __| / _` | '_ \| | / __| | |
#     | (_| | (_) || | |  __/\__ \ |_ | (_| | |_) | || (__| | |
#      \__, |\___(_)_|  \___||___/\__(_)__,_| .__/|_(_)___|_|_|
#       __/ |                               | |                
#      |___/                                |_|                                            
                ..-*@@@@@@@+:.        .:+@@@@@@@*-..                                
                .*@@@@@@@@@@@@@+.    .+@@@@@@@@@@@@@*.                               
              .%@@@@*:...:+@@@@@-..-@@@@@+:...:+@@@@%.                              
              #@@@@.        =@@@@@@@@@@=        .@@@@#                              
              .@@@@=          .*@@@@@@*.          -@@@@:                             
              :@@@@-           .*@@@@*.           :@@@@-                             
              .@@@@+          .%@@@@@@%:          =@@@@.                             
              *@@@@-.     ..#@@@@**@@@@#..     .:@@@@*                              
              .*@@@@%=---=#@@@@#:..:#@@@@#=---=#@@@@#.                              
                .-%@@@@@@@@@@@%-.    .:#@@@@@@@@@@@@=.                               
                  ..-#@@@@@*:..        ..:*@@@@@#-..                                 
                                                                                                    
                                                                                                    
```

A small, cross-platform Go CLI for making REST API calls.

- ✅ Works on **Windows** and **Linux**
- ✅ **No external dependencies** – only Go standard library → great for **air-gapped** environments
- ✅ JSON payloads from file and inline, with **merge & override**
- ✅ **Profiles** for base URL, default headers, and auth
- ✅ **Auth strategies**: None, Basic, Bearer, **OAuth 2.0 Client Credentials**
- ✅ Dedicated **`token` command** to fetch, print, or save an OAuth2 token independently
- ✅ **Hashing** binary files (MD5 / SHA-1 / SHA-256) and injecting into JSON
- ✅ Optional hash formatting: `0x` prefix and uppercase hex
- ✅ Flexible output options: pretty JSON, raw, JSON-only, write to file
- ✅ Basic **retry logic**
- ✅ Clean architecture with **Command**, **Factory**, **Strategy** patterns
- ✅ Unit tests for core modules
- ✅ Built-in **version** command and ASCII banner

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
    │   ├── profile.go
    │   └── version.go
    ├── config
    │   └── config.go
    ├── httpclient
    │   └── factory.go
    ├── payload
    │   └── json.go
    └── version
        └── version.go
```

(Plus `*_test.go` files next to some modules for unit tests.)

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Module overview

### `main.go`

- Entry point for the CLI.
- Defines the ASCII art logo and banner.
- Prints banner with:
  - Tool name: `go-rest-api-cli`
  - Version, commit, and build date (from `internal/version`).
- Creates a `command.Registry`, registers:
  - `call`
  - `profile`
  - `inspect`
  - `help`
  - `version`
- Behavior:
  - No arguments → print banner + global help (via `help` command) and exit.
  - Unknown command → print banner + error + global help.
  - On command error → print error, banner, and command-specific help (fallback to global).

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
  - Hashing (`--hash-file`, `--hash-algo`, `--hash-field`, `--hash-prefix-0x`, `--hash-upper`)
  - Headers (`--header`)
  - Auth (`--auth`, `--user`, `--pass`, `--token`)
  - Output modes (`--pretty`, `--raw`, `--json-only`, `--out`)
  - Retry (`--retries`, `--retry-delay`)
  - Network behavior (`--timeout`, `--insecure`)
- Load profile from config (if used).
- Merge:
  - Profile base URL + relative `--url`
  - Profile headers + CLI `--header`
  - Profile auth + CLI auth overrides
- Load & merge JSON from file + inline.
- Compute file hash (optional) and inject into JSON.
  - Optionally rewrite the JSON file on disk with the new hash field.
  - Print the computed hash to the console.
- Build request/client via `httpclient.Factory`.
- Perform request with retry logic.
- Format and print/save response according to selected output strategy.

#### `profile.go`

Manages saved profiles.

Commands:

- `profile add` – create/update a profile with:
  - `--name NAME`
  - `--base-url URL`
  - `--auth none|basic|bearer|oauth2`
  - optional `--user`, `--pass`, `--token`
  - `--oauth2-token-url`, `--oauth2-client-id`, `--oauth2-client-secret`, `--oauth2-scopes`
  - `--header "Key: Value"` (repeatable)
- `profile list` – list all profiles (basic info).
- `profile remove` – delete a profile by name.

Uses `internal/config` for persistence.

#### `token.go`

Standalone command for fetching an OAuth2 access token without making an API call.

Flags:

- `--profile` – load OAuth2 settings from a saved profile
- `--oauth2-token-url` – token endpoint URL
- `--oauth2-client-id` – client ID
- `--oauth2-client-secret` – client secret
- `--oauth2-scopes` – space-separated scopes
- `--json` – print the full token response as JSON (`access_token`, `token_type`, `expires_in`, `scope`) instead of just the raw token string
- `--quiet` – suppress all console output (useful when only writing to file)
- `--out FILE` – write the token (or JSON) to a file; file is written with mode `0600`

CLI flags take precedence over profile defaults, so individual fields can be overridden at call time.

#### `inspect.go`

Pretty printing of stored profiles.

Commands:

- `inspect profiles` – show all profiles with details.
- `inspect profile --name NAME` – show a single profile in detail (base URL, auth, headers).

#### `help.go`

- Implements the `help` command.
- Prints:
  - Tool description
  - Available commands
  - Quick examples
- Supports:
  - `go-rest-api-cli help` – global help
  - `go-rest-api-cli help call` – command-specific help, etc.

#### `version.go`

- Implements the `version` command.
- Commands:
  - `version` – prints full version string with commit and build date.
  - `version --short` – prints only the semantic version.
- Uses `internal/version` for version metadata.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### `internal/config`

Handles storage of **profiles** in a config file.

#### `config.go`

- Types:
  - `Profile` – base URL, headers, auth type, user/pass/token, OAuth2 fields (`oauth2_token_url`, `oauth2_client_id`, `oauth2_client_secret`, `oauth2_scopes`).
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
  - `OAuth2ClientCredentials` – performs the **OAuth 2.0 Client Credentials** grant:
    - `FetchTokenFull()` POSTs `grant_type=client_credentials` + `client_id` + `client_secret` (+ optional `scope`) to the token endpoint and returns a `*TokenResponse` (`access_token`, `token_type`, `expires_in`, `scope`).
    - `FetchToken()` calls `FetchTokenFull()` and caches only the `access_token` for `Apply()`.
    - `Apply()` sets `Authorization: Bearer <access_token>` on the request.
    - Returns structured errors for missing config, HTTP failures, and server-side OAuth2 errors (`error` / `error_description`).

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
- `NormalizeHashAlgo(algo)` – normalizes names like `"SHA-256"` → `"sha-256"`.
- `ComputeFileHash(path, algo)`:
  - Supports `md5`, `sha-1`, `sha-256`.
  - Returns hex string hash (lowercase – can be uppercased by the `call` command).
- The `call` command uses this module to build the final request body and compute hashes.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### `internal/version`

#### `version.go`

- Stores build-time metadata:
  - `Version` – semantic version (default `dev`).
  - `Commit` – git commit hash (default `none`).
  - `Date` – build date (default `unknown`).
- `Full()` – returns a formatted string:

  ```text
  v1.0.0 (commit a1b2c3d4, built 2026-01-20T18:45:00Z)
  ```

- Values are intended to be overridden using `go build -ldflags`.

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

### Build with embedded version info

Embed version, commit, and build date (for the `version` command and banner):

#### Linux binary with version metadata

```bash
go build \
  -ldflags "\
    -X go-rest-api-cli/internal/version.Version=v1.0.0 \
    -X go-rest-api-cli/internal/version.Commit=$(git rev-parse HEAD) \
    -X go-rest-api-cli/internal/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o go-rest-api-cli .
```

#### Windows binary with version metadata (built from Linux/macOS)

```bash
GOOS=windows GOARCH=amd64 go build \
  -ldflags "\
    -X go-rest-api-cli/internal/version.Version=v1.0.0 \
    -X go-rest-api-cli/internal/version.Commit=$(git rev-parse HEAD) \
    -X go-rest-api-cli/internal/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o go-rest-api-cli.exe .
```

### Verify installation

Linux/macOS:

```bash
./go-rest-api-cli version
./go-rest-api-cli
```

Windows (PowerShell):

```powershell
.\go-rest-api-cli.exe version
.\go-rest-api-cli.exe
```

- Running without arguments prints:
  - ASCII banner
  - Application name
  - Version, commit, build date
  - Global help
- `version` prints detailed version info, `version --short` prints just the version.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## Commands overview

General syntax:

```bash
go-rest-api-cli <command> [flags...]
```

Available commands:

- `call` – execute a REST API call
- `token` – fetch an OAuth2 access token (print / save / suppress)
- `profile` – manage saved profiles
- `inspect` – inspect stored profiles
- `version` – show version information
- `help` – show help and examples

Running the binary with **no command** or with an **unknown command** prints the ASCII banner and help.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

## `call` command

### Key flags

- `--method` – HTTP method (default: `GET`)
- `--url` – request URL (**required**)
  - absolute: `https://api.example.com/v1/users`
  - or relative: `/v1/users` when using `--profile`
- `--profile` – profile name (base URL, header defaults, auth defaults)
- `--header "Key: Value"` – extra headers (repeatable)
- `--json-file` – JSON file path for payload (merged base)
- `--data` – inline JSON to merge/override file payload
- `--auth` – `none|basic|bearer|oauth2`
- `--user`, `--pass`, `--token` – auth parameters (basic / bearer)
- `--oauth2-token-url` – OAuth2 token endpoint URL
- `--oauth2-client-id` – OAuth2 client ID
- `--oauth2-client-secret` – OAuth2 client secret
- `--oauth2-scopes` – space-separated OAuth2 scopes (e.g. `"read write"`)
- `--timeout` – timeout in seconds (default `30`)
- `--insecure` – skip TLS verification (lab only)
- `--pretty` – pretty-print JSON responses
- `--raw` – print only body
- `--json-only` – print only JSON body (when Content-Type is JSON)
- `--out` – write response body to file
- `--retries` – number of retries (network/5xx)
- `--retry-delay` – delay between retries (seconds)
- `--hash-file` – path to file to hash
- `--hash-algo` – `md5|sha-1|sha-256` (default `sha-256`)
- `--hash-field` – JSON field name for hash (default `file_hash`)
- `--hash-prefix-0x` – if set, prefix hash with `0x`
- `--hash-upper` – if set, render hex in uppercase (`A-F`)

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
3. Optionally update the JSON file on disk.
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

- Computes `sha-256(build/artifact.bin)` → e.g. `"b94d27b9..."`.
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

#### 5.2 MD5 with `0x` prefix and uppercase into `md5_sum`

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
  --hash-upper \
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
  --hash-upper `
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
  --hash-upper ^
  --pretty
```

Now:

- Hash looks like `"0x5EB63BBBE01EEED093CB22BB8F5ACDC3"`.
- JSON is updated:

  ```json
  {
    "name": "artifact upload",
    "build": 42,
    "checksum": "b94d27b9...",
    "md5_sum": "0x5EB63BBBE01EEED093CB22BB8F5ACDC3"
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
  --url "https://api.agify.io/?name=meelad" \
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

### 8. OAuth 2.0 Client Credentials

The `--auth oauth2` strategy fetches a Bearer token from a token endpoint before the API call.
The token is obtained using the **Client Credentials** grant (machine-to-machine, no user interaction).

#### 8.1 One-off call with inline OAuth2 flags

##### Bash (Linux/macOS)

```bash
./go-rest-api-cli call \
  --method GET \
  --url "https://api.example.com/v1/resources" \
  --auth oauth2 \
  --oauth2-token-url "https://auth.example.com/oauth2/token" \
  --oauth2-client-id "myapp-client-id" \
  --oauth2-client-secret "myapp-client-secret" \
  --oauth2-scopes "read write" \
  --pretty
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --method GET `
  --url "https://api.example.com/v1/resources" `
  --auth oauth2 `
  --oauth2-token-url "https://auth.example.com/oauth2/token" `
  --oauth2-client-id "myapp-client-id" `
  --oauth2-client-secret "myapp-client-secret" `
  --oauth2-scopes "read write" `
  --pretty
```

##### cmd.exe

```bat
go-rest-api-cli.exe call ^
  --method GET ^
  --url "https://api.example.com/v1/resources" ^
  --auth oauth2 ^
  --oauth2-token-url "https://auth.example.com/oauth2/token" ^
  --oauth2-client-id "myapp-client-id" ^
  --oauth2-client-secret "myapp-client-secret" ^
  --oauth2-scopes "read write" ^
  --pretty
```

What happens:

1. The tool POSTs to the token URL:
   ```
   POST https://auth.example.com/oauth2/token
   Content-Type: application/x-www-form-urlencoded

   grant_type=client_credentials&client_id=myapp-client-id&client_secret=myapp-client-secret&scope=read+write
   ```
2. Parses the `access_token` from the JSON response.
3. Sets `Authorization: Bearer <access_token>` on the actual API request.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 8.2 Save OAuth2 credentials to a profile

Store all OAuth2 settings in a profile so you only need `--profile` on every call.

##### Bash

```bash
./go-rest-api-cli profile add \
  --name myapi \
  --base-url "https://api.example.com" \
  --auth oauth2 \
  --oauth2-token-url "https://auth.example.com/oauth2/token" \
  --oauth2-client-id "myapp-client-id" \
  --oauth2-client-secret "myapp-client-secret" \
  --oauth2-scopes "read write"
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe profile add `
  --name myapi `
  --base-url "https://api.example.com" `
  --auth oauth2 `
  --oauth2-token-url "https://auth.example.com/oauth2/token" `
  --oauth2-client-id "myapp-client-id" `
  --oauth2-client-secret "myapp-client-secret" `
  --oauth2-scopes "read write"
```

##### cmd.exe

```bat
go-rest-api-cli.exe profile add ^
  --name myapi ^
  --base-url "https://api.example.com" ^
  --auth oauth2 ^
  --oauth2-token-url "https://auth.example.com/oauth2/token" ^
  --oauth2-client-id "myapp-client-id" ^
  --oauth2-client-secret "myapp-client-secret" ^
  --oauth2-scopes "read write"
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 8.3 Use the OAuth2 profile for API calls

Once the profile is saved, just pass `--profile`:

##### Bash

```bash
./go-rest-api-cli call \
  --profile myapi \
  --method GET \
  --url "/v1/resources" \
  --pretty
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe call `
  --profile myapi `
  --method GET `
  --url "/v1/resources" `
  --pretty
```

##### cmd.exe

```bat
go-rest-api-cli.exe call ^
  --profile myapi ^
  --method GET ^
  --url "/v1/resources" ^
  --pretty
```

The tool will automatically:
1. Resolve `https://api.example.com/v1/resources` (base URL + relative path).
2. Fetch an OAuth2 token from the saved token endpoint.
3. Attach `Authorization: Bearer <token>` to the request.

> **Tip:** You can override individual OAuth2 values at call time by passing the corresponding flag (e.g. `--oauth2-scopes "admin"`) alongside `--profile`. CLI flags always take precedence over profile defaults.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

### 9. `token` command – fetch an OAuth2 token independently

Use `token` when you need the access token itself (e.g. to pass it to another tool, script, or environment variable) without making a full API call.

#### 9.1 Print access token to console (default)

##### Bash

```bash
./go-rest-api-cli token \
  --oauth2-token-url "https://auth.example.com/oauth2/token" \
  --oauth2-client-id "myapp-client-id" \
  --oauth2-client-secret "myapp-client-secret"
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe token `
  --oauth2-token-url "https://auth.example.com/oauth2/token" `
  --oauth2-client-id "myapp-client-id" `
  --oauth2-client-secret "myapp-client-secret"
```

##### cmd.exe

```bat
go-rest-api-cli.exe token ^
  --oauth2-token-url "https://auth.example.com/oauth2/token" ^
  --oauth2-client-id "myapp-client-id" ^
  --oauth2-client-secret "myapp-client-secret"
```

Output (just the raw token string):

```
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 9.2 Print full token response as JSON (`--json`)

Returns `access_token`, `token_type`, `expires_in`, and `scope`.

##### Bash

```bash
./go-rest-api-cli token \
  --oauth2-token-url "https://auth.example.com/oauth2/token" \
  --oauth2-client-id "myapp-client-id" \
  --oauth2-client-secret "myapp-client-secret" \
  --oauth2-scopes "read write" \
  --json
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe token `
  --oauth2-token-url "https://auth.example.com/oauth2/token" `
  --oauth2-client-id "myapp-client-id" `
  --oauth2-client-secret "myapp-client-secret" `
  --oauth2-scopes "read write" `
  --json
```

##### cmd.exe

```bat
go-rest-api-cli.exe token ^
  --oauth2-token-url "https://auth.example.com/oauth2/token" ^
  --oauth2-client-id "myapp-client-id" ^
  --oauth2-client-secret "myapp-client-secret" ^
  --oauth2-scopes "read write" ^
  --json
```

Output:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "read write"
}
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 9.3 Save token to a file, suppress console output (`--out` + `--quiet`)

##### Bash

```bash
./go-rest-api-cli token \
  --oauth2-token-url "https://auth.example.com/oauth2/token" \
  --oauth2-client-id "myapp-client-id" \
  --oauth2-client-secret "myapp-client-secret" \
  --out token.txt \
  --quiet
```

##### PowerShell

```powershell
.\go-rest-api-cli.exe token `
  --oauth2-token-url "https://auth.example.com/oauth2/token" `
  --oauth2-client-id "myapp-client-id" `
  --oauth2-client-secret "myapp-client-secret" `
  --out token.txt `
  --quiet
```

##### cmd.exe

```bat
go-rest-api-cli.exe token ^
  --oauth2-token-url "https://auth.example.com/oauth2/token" ^
  --oauth2-client-id "myapp-client-id" ^
  --oauth2-client-secret "myapp-client-secret" ^
  --out token.txt ^
  --quiet
```

`token.txt` will contain just the raw `access_token` string (file mode `0600`). Nothing is printed to the console.
Use `--json --out token.json --quiet` to save the full JSON response silently instead.

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 9.4 Use a saved profile (`--profile`)

```bash
# Save profile first (see section 8.2)
./go-rest-api-cli token --profile myapi --json
```

PowerShell:

```powershell
.\go-rest-api-cli.exe token --profile myapi --json
```

cmd.exe:

```bat
go-rest-api-cli.exe token --profile myapi --json
```

<hr style="border: 0; height: 1px; background-image: linear-gradient(to right, rgba(0,0,0,0), rgba(0,122,204,0.75), rgba(0,0,0,0));">

#### 9.5 Capture token into a shell variable

##### Bash

```bash
TOKEN=$(./go-rest-api-cli token \
  --oauth2-token-url "https://auth.example.com/oauth2/token" \
  --oauth2-client-id "myapp-client-id" \
  --oauth2-client-secret "myapp-client-secret")

echo "Got token: $TOKEN"

# Use it in a curl call
curl -H "Authorization: Bearer $TOKEN" https://api.example.com/v1/data
```

##### PowerShell

```powershell
$TOKEN = .\go-rest-api-cli.exe token `
  --oauth2-token-url "https://auth.example.com/oauth2/token" `
  --oauth2-client-id "myapp-client-id" `
  --oauth2-client-secret "myapp-client-secret"

Write-Host "Got token: $TOKEN"

# Use it in a subsequent call
.\go-rest-api-cli.exe call `
  --url "https://api.example.com/v1/data" `
  --auth bearer `
  --token $TOKEN `
  --pretty
```

##### cmd.exe

```bat
for /f "delims=" %%T in ('go-rest-api-cli.exe token ^
  --oauth2-token-url "https://auth.example.com/oauth2/token" ^
  --oauth2-client-id "myapp-client-id" ^
  --oauth2-client-secret "myapp-client-secret"') do set TOKEN=%%T

echo Got token: %TOKEN%
```

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
  - Subcommands: `CallCommand`, `ProfileCommand`, `InspectCommand`, `HelpCommand`, `VersionCommand`.

- **Strategy pattern (Auth)**
  - `internal/auth.Strategy` interface.
  - Concrete strategies: `NoAuth`, `Basic`, `Bearer`, `OAuth2ClientCredentials`.

- **Factory pattern (HTTP)**
  - `internal/httpclient.Factory` builds `*http.Request` and `*http.Client` from a config.

- **Configuration module**
  - `internal/config` isolates config file location & schema (profiles).

- **Versioning module**
  - `internal/version` exposes build-time metadata used by the banner and `version` command.

- **Separation of concerns**
  - CLI parsing & orchestration → `internal/command`
  - HTTP behavior → `internal/httpclient`
  - Auth logic → `internal/auth`
  - JSON & hashing → `internal/payload`
  - Persistence of profiles → `internal/config`

## License

go-rest-api-cli is released under the terms of the **MIT License**, a permissive open-source license that allows extensive reuse with minimal restrictions.

You are free to:

- Use the software for personal, academic, or commercial purposes.
- Modify the source code to fit your own requirements.
- Distribute original or modified versions.
- Include go-rest-api-cli as part of your own tools or products.

Conditions:

- You must retain the original copyright notice.
- You must include a copy of the MIT License in any substantial portions of the software you distribute.

For the full legal text, see the `LICENSE` file in this repository.
