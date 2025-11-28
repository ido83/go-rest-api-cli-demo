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


## 🌐 Usage Examples

The CLI tool uses a subcommand structure, typically starting with `call`.

### 1. HTTP GET Request

This example performs a simple GET request to the Agify API.

```bat
go-rest-api-cli.exe call --method GET --url "[https://api.agify.io/?name=meelad](https://api.agify.io/?name=meelad)"


### 2. HTTP POST Request with JSON Payload

To send a POST request with a body, use the `--json-file` flag pointing to a local file (e.g., `payload.json`).

**Windows Batch Example:**
```bat
go-rest-api-cli.exe call --method POST --url "[https://api.restful-api.dev/objects](https://api.restful-api.dev/objects)" --json-file "payload.json"


**PowerShell Example:**
Note the use of `.\` to execute the file in the current directory in PowerShell.
```ps1
.\go-rest-api-cli.exe call --method POST --url "[https://api.restful-api.dev/objects](https://api.restful-api.dev/objects)" --json-file "payload.json"


### 3. Comparison with Other Tools

This table shows how the equivalent POST request might look using other common tools:

| Tool | Command |
| :--- | :--- |
| **Current CLI** | `go-rest-api-cli.exe call --method POST --url "..." --json-file "payload.json"` |
| **cURL** | `curl -X POST -H "Content-Type: application/json" -d @payload.json "https://api.restful-api.dev/objects"` |
| **restcli (Example GET)** | `restcli call --method GET --url "https://api.agify.io/?name=meelad"` |

