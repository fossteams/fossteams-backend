module github.com/fossteams/fossteams-backend

go 1.16

require (
	github.com/ReneKroon/ttlcache/v2 v2.7.0
	github.com/alexflint/go-arg v1.4.2
	github.com/fossteams/teams-api v0.0.0-20210608193737-ead87df795c2
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
)

replace github.com/fossteams/teams-api => ../teams-api
