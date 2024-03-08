package main

import (
	"fmt"
	"io"
	"net/http"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool // connections we are maintaining
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true // consider using a mutex for race conditions

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil { // how do we want to handle this (work with malformed message?)
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))

		ws.Write([]byte("thank you for the message"))
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":3000", nil)
}