module github.com/fossteams/fossteams-backend

go 1.18

require (
	github.com/ReneKroon/ttlcache/v2 v2.7.0
	github.com/alexflint/go-arg v1.4.2
	github.com/fossteams/teams-api v0.0.0-20210608193737-ead87df795c2
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.2
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
)

require (
	github.com/alexflint/go-scalar v1.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/fossteams/teams-api => ../teams-api
