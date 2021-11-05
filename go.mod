module plantumlive-ws-server

go 1.17

require (
	github.com/daixiang0/gci v0.2.9
	github.com/golang/mock v1.6.0
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	mvdan.cc/gofumpt v0.1.1
)

require (
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/subcommands v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/tools v0.1.1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace plantumlive-ws-server/wsmsg => ./wsmsg
