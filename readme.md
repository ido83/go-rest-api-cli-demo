# go-rest-api-cli

A small, cross-platform Go CLI tool for making REST API calls.

- Works on **Windows** and **Linux**
- No external dependencies → good for **air-gapped** environments
- Supports **GET/POST/PUT/DELETE/etc.**
- Supports **JSON body** from:
    - inline `--data`
    - JSON file `--json-file`
    - **merged** (file + inline) with override
- Supports **headers** and simple **auth strategies** (none/basic/bearer)
- Uses Go OOP-style design: **Command**, **Factory**, **Strategy (Auth)**

---

## Features

- **Command pattern**: subcommands like `call`, `help`
- **Factory pattern** for constructing HTTP request + client
- **Auth strategy pattern**:
    - no auth
    - basic auth
    - bearer token
- **JSON payload handling**:
    - load from file
    - parse inline JSON string
    - merge file + inline (inline overrides duplicate keys)
- **Cross-platform**:
    - single static binary
    - no third-party libraries

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
    command/
      command.go       # Command interface & registry
      headers.go       # HeaderFlag for repeated --header
      call.go          # "call" command implementation
      help.go          # "help" command implementation



High-level architecture
main.go
Creates a command registry
Registers call and help commands
Dispatches based on os.Args[1]
internal/command
Command interface and registry
CallCommand implements the REST API logic
HelpCommand prints usage information
internal/httpclient
Factory builds *http.Request and *http.Client from a config struct
internal/auth
Auth strategies implementing a common Strategy interface
internal/payload
JSON helpers: load, parse, merge

Getting started
Requirements:
Go (any modern version, e.g. 1.20+)
Clone / copy
Copy the source files into a folder, for example:

mkdir go-rest-api-cli
cd go-rest-api-cli
# place main.go and internal/ here
```

Example go.mod
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
```
Main commands:
call – execute a REST API call
help – show help and examples
call command flags

--method
HTTP method, default: GET
(e.g. GET, POST, PUT, DELETE, PATCH …)

--url
Request URL (required).

--data
Inline JSON string.
Example: '{"name":"test"}' (see Windows notes below).

--json-file
Path to a JSON file. The payload from the file is merged with --data.
Inline JSON overrides duplicate keys from the file.

--header
Extra HTTP headers, can be used multiple times.
Format: "Key: Value"
Example: --header "X-Env: dev" --header "X-Client: cli"

--timeout
Timeout in seconds (default: 30).

--insecure
Skip TLS verification (InsecureSkipVerify=true).
Not recommended for production, but sometimes useful in labs.

--auth
Auth type: none (default), basic, bearer.

--user / --pass
Username/password for --auth basic.

--token
Bearer token for --auth bearer.
```

Examples
1. Simple GET (agify API)
Predict age from name via https://api.agify.io/?name=meelad.
Linux / macOS / Windows PowerShell
```./go-rest-api-cli call \
  --method GET \
  --url "https://api.agify.io/?name=meelad"
```

Windows cmd.exe
```go-rest-api-cli.exe call --method GET --url "https://api.agify.io/?name=meelad"```

2. POST with JSON file (restful-api.dev)
Target API: https://api.restful-api.dev/objects
Create payload.json:
```
{
  "name": "Base Object",
  "data": {
    "env": "prod",
    "version": 1
  }
}
```

Call:
PowerShell / Linux / macOS
```
./go-rest-api-cli call \
  --method POST \
  --url "https://api.restful-api.dev/objects" \
  --json-file "payload.json"
```

Windows cmd.exe
```go-rest-api-cli.exe call --method POST --url "https://api.restful-api.dev/objects" --json-file "payload.json"```


3. POST with JSON file + inline override

Same payload.json as above:
```{
  "name": "Base Object",
  "data": {
    "env": "prod",
    "version": 1
  }
}
```

PowerShell / Linux / macOS
```./go-rest-api-cli call \
  --method POST \
  --url "https://api.restful-api.dev/objects" \
  --json-file "payload.json" \
  --data '{"name":"Overridden Name","data":{"extra":"from-inline"}}'
```

File payload:
```{
  "name": "Base Object",
  "data": {
    "env": "prod",
    "version": 1
  }
}
```

Inline payload:
```{
  "name": "Overridden Name",
  "data": {
    "extra": "from-inline"
  }
}
```

Merged body actually sent:
```{
  "name": "Overridden Name",
  "data": {
    "env": "prod",
    "version": 1,
    "extra": "from-inline"
  }
}
```

Windows cmd.exe (note the escaping)
```
go-rest-api-cli.exe call --method POST --url "https://api.restful-api.dev/objects" --json-file "payload.json" --data "{\"name\":\"Overridden Name\",\"data\":{\"extra\":\"from-inline\"}}"

```

4. Adding headers
```
./go-rest-api-cli call \
  --method GET \
  --url "https://httpbin.org/headers" \
  --header "X-Env: dev" \
  --header "X-Client: go-rest-api-cli"
```

5. Basic auth example
```./go-rest-api-cli call \
  --method GET \
  --url "https://httpbin.org/basic-auth/user/pass" \
  --auth basic \
  --user user \
  --pass pass
```

6. Bearer token auth example
```./go-rest-api-cli call \
  --method GET \
  --url "https://httpbin.org/bearer" \
  --auth bearer \
  --token "YOUR_TOKEN_HERE"
```





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

