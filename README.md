# fossteams-backend

A Golang based web server that exposes some endpoints to interact with a
Microsoft Teams session. This allows different clients (such as
[fossteams-frontend](https://github.com/fossteams/fossteams-frontend)) to
interact with Microsoft Teams backends without having to deal with
changing APIs and authentication.

## Requirements

- A Microsoft Teams account
- Go
- [teams-token](https://github.com/fossteams/teams-token)

## Running the server

1. Get a token with [teams-token](https://github.com/fossteams/teams-token)
1. `go run ./cmd/fossteams-backend/`

## Check if the server is working properly

Visit: http://127.0.0.1:8050/api/v1/conversations, if you get a non-empty JSON
and no errors on the console (or a 200 OK status code), all good!
