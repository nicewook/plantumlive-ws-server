# websocket with Gorilla


Posting link: https://golangdocs.com/golang-gorilla-websockets

GitHub: https://github.com/gorilla/websocket

chatting example
https://github.com/gorilla/websocket/tree/master/examples/chat

## General design

- server and ws

### server

- bind a port and start listening for connection
- connection act over a raw HTTP connection

### ws

- try to connect with the server using websocket URL
- BTW, you don't necessarilly need to implement with Golang

