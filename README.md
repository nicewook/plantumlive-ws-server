# websocket with Gorilla


Posting link: https://golangdocs.com/golang-gorilla-websockets

GitHub: https://github.com/gorilla/websocket

## General design

- server and client

### server

- bind a port and start listening for connection
- connection act over a raw HTTP connection

### client

- try to connect with the server using websocket URL
- BTW, you don't necessarilly need to implement with Golang

