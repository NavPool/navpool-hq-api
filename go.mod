module github.com/NavPool/navpool-hq-api

go 1.12

require (
	github.com/NavExplorer/navexplorer-sdk-go v0.0.0-20190524154412-6693c72d5c6a
	github.com/appleboy/gin-jwt v2.6.0+incompatible
	github.com/certifi/gocertifi v0.0.0-20190506164543-d2eda7129713 // indirect
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/getsentry/raven-go v0.2.0
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-gonic/gin v1.4.0
	github.com/jinzhu/gorm v1.9.8
	github.com/satori/go.uuid v1.2.0
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	gopkg.in/dgrijalva/jwt-go.v3 v3.2.0 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
