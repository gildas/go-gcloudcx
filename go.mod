module github.com/gildas/go-purecloud

go 1.13

require (
	github.com/gildas/go-core v0.4.1
	github.com/gildas/go-errors v0.0.1
	github.com/gildas/go-logger v1.3.1
	github.com/gildas/go-request v0.2.2
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/pkg/errors v0.8.2-0.20191109102316-7f95ac13edff
	github.com/stretchr/testify v1.4.0
)

replace github.com/gildas/go-logger => ../go-logger

replace github.com/gildas/go-errors => ../go-errors

replace github.com/gildas/go-request => ../go-request
