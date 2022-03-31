module github.com/gildas/go-gcloudcx

go 1.16

require (
	github.com/gildas/go-core v0.4.10
	github.com/gildas/go-errors v0.3.1
	github.com/gildas/go-logger v1.5.2
	github.com/gildas/go-request v0.7.6
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.5.0
	github.com/joho/godotenv v1.4.0
	github.com/matoous/go-nanoid/v2 v2.0.0
	github.com/stretchr/testify v1.7.1
	golang.org/x/sys v0.0.0-20220330033206-e17cdc41300f // indirect
	google.golang.org/genproto v0.0.0-20220329172620-7be39ac1afc7 // indirect
)

replace github.com/gildas/go-errors => ../go-errors
