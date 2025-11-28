bash

go mod tidy
GOOS=windows GOARCH=amd64 go build -o go-rest-api-cli.exe

bat

go-rest-api-cli.exe call --method GET --url "https://api.agify.io/?name=meelad"
go-rest-api-cli.exe call --method POST --url "https://api.restful-api.dev/objects" --json-file "payload.json"


restcli call \
  --method GET \
  --url "https://api.agify.io/?name=meelad"



