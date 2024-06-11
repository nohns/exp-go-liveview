module github.com/nohns/go-liveview

go 1.22.0

replace github.com/a-h/templ => ../exp-templ-live

require (
	github.com/a-h/templ v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.2
)

require golang.org/x/net v0.24.0 // indirect
