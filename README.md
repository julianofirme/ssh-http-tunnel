## SSH/HTTP Tunnel
Just a simple ssh/http tunnel built with go

## If you want to test

- Run the server with `go run main.go`

- Open a new tab in the terminal and execute `ssh -tt localhost -p 2222 <<< "some message here"`

- In the tab where `main.go` is running, an ID will appear. Access it in your browser by going to `localhost:8080/?id=paste_the_id`
