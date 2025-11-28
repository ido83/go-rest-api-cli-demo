# go-rest-api-cli

A small, cross-platform Go CLI tool for making REST API calls.

- Works on **Windows** and **Linux**
- No external dependencies → good for **air-gapped** environments
- Supports **GET/POST/PUT/DELETE/etc.**
- Supports **JSON body** from:
    - inline `--data`
    - JSON file `--json-file`
    - **merged** (file + inline) with override
- Profiles for **base URL + default headers/auth**
- Output strategies: `--pretty`, `--raw`, `--json-only`
- Save response to file: `--out`
- Retry logic: `--retries`, `--retry-delay`
- Uses Go “OOP-style” design: **Command**, **Factory**, **Strategy (Auth)**, config module

## Features (current)
### Commands

- `call` – execute a REST API call
- `profile` – manage saved profiles:
    - `profile add`
    - `profile list`
    - `profile remove`
- `inspect` – inspect stored profiles:
    - `inspect profiles`
    - `inspect profile --name NAME`
- `help` – show help and examples

### Profiles

Profiles let you save:

- Base URL (e.g. `https://api.example.com`)
- Default headers (`X-Env`, auth headers, etc.)
- Default auth:
    - `none`
    - `basic` (user/pass)
    - `bearer` (token)

Then you call APIs with `--profile` so you don’t repeat all parameters each time.

### Output strategies

- `--pretty`  
  Pretty-print JSON responses (if `Content-Type` is JSON).
- `--json-only`  
  Only print the JSON body (no status/headers).
- `--raw`  
  Print **only** the response body (any content).

Order of precedence:

- If `--json-only` is set → JSON body only (pretty if `--pretty`).
- Else if `--raw` is set → body only as-is (or pretty JSON).
- Else (default) → status, headers, then body.

### Save response to a file

- `--out path/to/file.json`  
  Writes the final printed body (raw or pretty JSON) to a file.

### Retry logic

- `--retries N` – number of retries on:
    - network errors
    - HTTP `5xx` responses
- `--retry-delay SECONDS` – delay between retries

---

## Project structure

```
go-rest-api-cli/
  go.mod
  main.go
  internal/
    auth/
      auth.go          # Auth strategies (none, basic, bearer)
    httpclient/
      factory.go       # HTTP request/client factory
    payload/
      json.go          # JSON helpers (file, inline, merge)
    config/
      config.go        # Profiles + config file load/save
    command/
      command.go       # Command interface & registry
      headers.go       # HeaderFlag for repeated --header
      call.go          # "call" command implementation
      profile.go       # "profile" command (add/list/remove)
      inspect.go       # "inspect" command (view profiles)
      help.go          # "help" command


```
Code is structure:
main.go
*   Creates a command.Registry
*   Registers:
*
  * CallCommand
  * ProfileCommand
  * InspectCommand
  * HelpCommand
* Reads os.Args[1] to decide which command to run
  internal/command
* 
  * Command interface:
```
type Command interface {
    Name() string
    Description() string
    Run(args []string) error
}
```

* Registry holds and dispatches commands.
* CallCommand:
    * Parses CLI flags (--method, --url, --profile, --data, --json-file, --pretty, --raw, --json-only, --out, --retries, etc.)
    * Loads profile (if --profile is used)
    * Merges:
        * profile base URL + relative --url
        * profile headers + --header flags
        * profile auth + CLI auth flags
    * Builds httpclient.Config and uses httpclient.Factory to create an HTTP request & client
    * Handles retries
    * Handles formatting of the response and writing to file
* ProfileCommand:
    * profile add ... -> loads config, adds/updates profile, saves
    * profile list -> prints all profiles 
    * profile remove --name NAME -> deletes a profile
* InspectCommand:
  * inspect profiles -> shows a detailed list of all profiles
  * inspect profile --name NAME -> shows details for a single profile
* HeaderFlag in headers.go:
  * Implements flag.Value so you can pass --header "Key: Value" multiple times.
* internal/config
  * Manages a Config struct that contains a map[string]Profile.
  * Profile includes:
    * Name, BaseURL, Headers
    * AuthType, User, Pass, Token
  * Knows where to store file:
    * Uses os.UserConfigDir() (fallback to ~/.go-rest-api-cli) and writes config.json.
  * Load() / Save() handle reading & writing JSON configuration.
  internal/httpclient
  * Config struct holds all request details.
  * Factory.Build(cfg):
    * Builds *http.Request from method, URL, headers, body.
    * Creates *http.Client with:
      * timeout
      * optional SkipTLSVerify
    * Applies the selected auth strategy.
  internal/auth
    * Strategy interface with Apply(req *http.Request).
    * NoAuth, Basic, Bearer structs implement it.
    * Used by CallCommand to apply auth in a pluggable way (Strategy pattern).
  * internal/payload
    * LoadJSONFile(path) → map[string]interface{}.
    * ParseJSONInline(string) → map[string]interface{}.
    * Merge(fileMap, inlineMap) → map[string]interface{} where inline overrides file keys.
    * Used by CallCommand to combine --json-file and --data.


Getting started:

Requirements: 
```Go 1.20+ (or adjust the go.mod version to your Go version)```

Example - go.mod: 
```
module go-rest-api-cli
go 1.22
```

Initialize dependencies:
```
go mod tidy
```
Building
Linux binary
```GOOS=linux GOARCH=amd64 go build -o go-rest-api-cli .```

Windows binary (from any OS)
```GOOS=windows GOARCH=amd64 go build -o go-rest-api-cli.exe .```

Copy the resulting binary (go-rest-api-cli or go-rest-api-cli.exe) to your target machine (including air-gapped environments).

Usage

General syntax:
```go-rest-api-cli <command> [flags]```

Commands overview:
```
call – perform HTTP request
profile – manage profiles
inspect – show saved profiles
help – show help
```

Call command:
```
Flags (main ones)

--method
HTTP method (default: GET).

--url
Request URL:

absolute: https://api.example.com/v1/users

or relative: /v1/users when using --profile.

--profile
Profile name to use (base URL, headers, auth).

--data
Inline JSON string.

--json-file
JSON file; merged with --data (inline overrides file).

--header
Extra header Key: Value (can be repeated).

--timeout
Timeout (seconds), default 30.

--insecure
Skip TLS verification (use only in dev/lab).

--auth / --user / --pass / --token
CLI-level auth. If --profile is used, profile auth is default and CLI flags override it.

--pretty
Pretty-print JSON response.

--raw
Print only the body.

--json-only
Print only JSON body (no status/headers).

--out
Save response body (after any pretty-print) to file.

--retries / --retry-delay
Number of retries on network/5xx errors and delay (seconds).
```



Examples

1. Simple GET (agify API)
```
Predict age from name via https://api.agify.io/?name=meelad.
Linux / macOS / Windows PowerShell
./go-rest-api-cli call \
  --method GET \
  --url "https://api.agify.io/?name=meelad"
```

Windows cmd.exe
```go-rest-api-cli.exe call --method GET --url "https://api.agify.io/?name=meelad"```

2. Create a profile for restful-api.dev
```aiignore
go-rest-api-cli.exe profile add ^
  --name restful ^
  --base-url https://api.restful-api.dev ^
  --auth none ^
  --header "X-Env: dev"

```

Check profiles:
```aiignore
go-rest-api-cli.exe profile list
go-rest-api-cli.exe inspect profiles
```

Inspect one:
```aiignore
go-rest-api-cli.exe inspect profile --name restful
```

3. GET with profile + pretty JSON
```aiignore
go-rest-api-cli.exe call ^
  --profile restful ^
  --method GET ^
  --url "/objects/1" ^
  --pretty

```

Because --url is relative and --profile has base-url = https://api.restful-api.dev, the final URL becomes:
```aiignore
https://api.restful-api.dev/objects/1

```

4. POST with JSON file + inline override

payload.json:
```aiignore
{
  "name": "Base Object",
  "data": {
    "env": "prod",
    "version": 1
  }
}

```

PowerShell / Linux / macOS
```aiignore
./go-rest-api-cli call \
  --profile restful \
  --method POST \
  --url "/objects" \
  --json-file "payload.json" \
  --data '{"name":"Overridden Name","data":{"extra":"from-inline"}}' \
  --pretty \
  --out "response.json"

```

Windows cmd.exe (note escaping):
```aiignore
go-rest-api-cli.exe call ^
  --profile restful ^
  --method POST ^
  --url "/objects" ^
  --json-file "payload.json" ^
  --data "{\"name\":\"Overridden Name\",\"data\":{\"extra\":\"from-inline\"}}" ^
  --pretty ^
  --out "response.json"

```

File payload:
```aiignore
{
  "name": "Base Object",
  "data": {
    "env": "prod",
    "version": 1
  }
}

```

Inline payload:
```aiignore
{
  "name": "Overridden Name",
  "data": {
    "extra": "from-inline"
  }
}

```

Final merged body (sent to server):
```aiignore
{
  "name": "Overridden Name",
  "data": {
    "env": "prod",
    "version": 1,
    "extra": "from-inline"
  }
}

```

5. Bearer auth with profile override
Create a profile with default bearer token:
```aiignore
go-rest-api-cli.exe profile add ^
  --name secureapi ^
  --base-url https://api.example.com ^
  --auth bearer ^
  --token "DEFAULT_TOKEN"

```
Call with profile, but override token:
```aiignore
go-rest-api-cli.exe call ^
  --profile secureapi ^
  --method GET ^
  --url "/v1/me" ^
  --auth bearer ^
  --token "OVERRIDE_TOKEN" ^
  --json-only ^
  --pretty

```

Profile provides base-url and default auth.
CLI --auth and --token override profile auth.

6. Retry logic
```aiignore
go-rest-api-cli.exe call ^
  --method GET ^
  --url "https://flaky.api.example.com/data" ^
  --retries 3 ^
  --retry-delay 2 ^
  --pretty

```

Behavior:

Up to 3 + 1 = 4 attempts total:
First try
+3 retries on network error or HTTP 5xx
Wait 2 seconds between retries.


### Windows quoting notes (very important)
Most JSON errors on Windows come from quoting.

```
.\go-rest-api-cli.exe call `
  --method POST `
--url "https://api.restful-api.dev/objects" `
--data '{"name":"Test Object","data":{"env":"dev","owner":"ido"}}'
```

* Outer quotes: '...'
* Inner JSON quotes: "..."
* No escaping needed.

cmd.exe
In cmd.exe, single quotes are not special, so use double quotes and escape inner quotes:

```aiignore
go-rest-api-cli.exe call --method POST --url "https://api.restful-api.dev/objects" --data "{\"name\":\"Test Object\",\"data\":{\"env\":\"dev\",\"owner\":\"ido\"}}"

```

If you see:
```aiignore
invalid character '\'' looking for beginning of value

```

it means the JSON string started with a literal '.
Use the double-quote + escaping form or move JSON to a file and use --json-file.


### Design patterns used (OOP)
Command pattern
* Command interface and Registry keep commands pluggable.
* Implementations: CallCommand, ProfileCommand, InspectCommand, HelpCommand.

Factory pattern
* httpclient.Factory builds HTTP requests/clients from a Config.
* Command code doesn’t deal with low-level HTTP details.

Strategy pattern (Auth)
* auth.Strategy with implementations:
    * NoAuth
    * Basic
    * Bearer
* CallCommand just picks a strategy and passes it to the factory.

Config module
* internal/config encapsulates:
    * where the config file lives
    * how to load/save
    * the Profile structure
* Commands just call Load() and Save().

Separation of concerns
* CLI parsing & orchestration in command package
* HTTP details in httpclient
* Auth in auth
* JSON payload operations in payload
* Persistent profiles in config


Extra Environment variables and build option
GOOS=windows GOARCH=amd64 go build -o go-rest-api-cli.exe


#### B. PowerShell (`.ps1`)

Set and then clear the environment variables.

```ps1
# Set the environment variables
$env:GOOS='windows'
$env:GOARCH='amd64'

# Run the build command
go build -o go-rest-api-cli.exe

# Optional: Clear the environment variables (good practice)
$env:GOOS=''
$env:GOARCH=''
```

#### C. Windows Batch (`.bat`)

Set and then clear the environment variables using the `set` command.

```bat
:: Set the environment variables
set GOOS=windows
set GOARCH=amd64

:: Run the build command
go build -o go-rest-api-cli.exe

:: Optional: Clear the environment variables (good practice)
set GOOS=
set GOARCH=
```

## 🌐 Usage Examples

The CLI tool uses a subcommand structure, typically starting with `call`.

### 1. HTTP GET Request

This example performs a simple GET request to the Agify API.

```bat
go-rest-api-cli.exe call --method GET --url "[https://api.agify.io/?name=meelad](https://api.agify.io/?name=meelad)"
```

### 2. HTTP POST Request with JSON Payload

To send a POST request with a body, use the `--json-file` flag pointing to a local file (e.g., `payload.json`).

**Windows Batch Example:**
```bat
go-rest-api-cli.exe call --method POST --url "[https://api.restful-api.dev/objects](https://api.restful-api.dev/objects)" --json-file "payload.json"
```

**PowerShell Example:**
Note the use of `.\` to execute the file in the current directory in PowerShell.
```ps1
.\go-rest-api-cli.exe call --method POST --url "[https://api.restful-api.dev/objects](https://api.restful-api.dev/objects)" --json-file "payload.json"
```

### 3. Comparison with Other Tools

This table shows how the equivalent POST request might look using other common tools:

| Tool | Command |
| :--- | :--- |
| **Current CLI** | `go-rest-api-cli.exe call --method POST --url "..." --json-file "payload.json"` |
| **cURL** | `curl -X POST -H "Content-Type: application/json" -d @payload.json "https://api.restful-api.dev/objects"` |
| **restcli (Example GET)** | `restcli call --method GET --url "https://api.agify.io/?name=meelad"` |

Quick sanity check new features (profile, inspect and etc...):
Build (Windows exe):
```GOOS=windows GOARCH=amd64 go build -o go-rest-api-cli.exe .```

Examples:
```
:: Create profile
go-rest-api-cli.exe profile add --name myapi --base-url https://api.restful-api.dev --auth none --header "X-Env: dev"

:: Inspect profiles
go-rest-api-cli.exe inspect profiles

:: Call with profile + pretty JSON
go-rest-api-cli.exe call --profile myapi --method GET --url "/objects/1" --pretty

:: POST with payload.json, retry 2 times, save to file
go-rest-api-cli.exe call ^
  --profile myapi ^
  --method POST ^
  --url "/objects" ^
  --json-file payload.json ^
  --retries 2 ^
  --retry-delay 2 ^
  --pretty ^
  --out "response.json"

```