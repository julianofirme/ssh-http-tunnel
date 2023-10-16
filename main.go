package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gliderlabs/ssh"
)

type Tunnel struct {
	writer io.Writer
	donech chan struct{}
}

var tunnels = map[int]chan Tunnel{}

func main() {

	fmt.Println("Starting tunnel, waiting for connection...")

	go func() {
		http.HandleFunc("/", handleRequest)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	ssh.Handle(func(s ssh.Session) {
		id := rand.Intn(math.MaxInt)
		tunnels[id] = make(chan Tunnel)

		fmt.Println("tunnel ID ->", id)

		tunnel := <-tunnels[id]
		fmt.Println("tunnel connection established")

		_, err := io.Copy(tunnel.writer, s)
		if err != nil {
			log.Fatal(err)
		}
		close(tunnel.donech)
		s.Write([]byte("we are done"))
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func handleRequest(writer http.ResponseWriter, reader *http.Request) {
	idstr := reader.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)

	tunnel, ok := tunnels[id]
	if !ok {
		writer.Write([]byte("tunnel not found"))
		return
	}

	donech := make(chan struct{})

	tunnel <- Tunnel{
		writer,
		donech,
	}

	<-donech
}
