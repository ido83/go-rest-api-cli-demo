bash
go mod init 
go mod tidy
GOOS=windows GOARCH=amd64 go build -o go-rest-api-cli.exe

ps1
# Set the environment variables
$env:GOOS='windows'
$env:GOARCH='amd64'

# Run the build command
go build -o go-rest-api-cli.exe

# Optional: Clear the environment variables (good practice)
$env:GOOS=''
$env:GOARCH=''

bat

:: Set the environment variables
set GOOS=windows
set GOARCH=amd64

:: Run the build command
go build -o go-rest-api-cli.exe

:: Optional: Clear the environment variables (good practice)
set GOOS=
set GOARCH=




bat

go-rest-api-cli.exe call --method GET --url "https://api.agify.io/?name=meelad"
go-rest-api-cli.exe call --method POST --url "https://api.restful-api.dev/objects" --json-file "payload.json"


ps1
.\go-rest-api-cli.exe call --method POST --url "https://api.restful-api.dev/objects" --json-file "payload.json"

curl -X POST -H "Content-Type: application/json" -d @payload.json "https://api.restful-api.dev/objects"


restcli call \
  --method GET \
  --url "https://api.agify.io/?name=meelad"



