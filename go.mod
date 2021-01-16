module go-bookstore-users-api

go 1.15

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0

require (
	github.com/Komdosh/go-bookstore-oauth-client v0.0.0-20210115060325-2a1a91078ff2
	github.com/Komdosh/go-bookstore-users-api v0.0.0-20210115060207-219326d2a46f
	github.com/Komdosh/golang-microservices v0.0.0-20201122120135-bcf27ea29be1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-resty/resty v0.0.0-00010101000000-000000000000 // indirect
	github.com/lib/pq v1.9.0
	github.com/stretchr/testify v1.4.0
	go.uber.org/zap v1.16.0
)
