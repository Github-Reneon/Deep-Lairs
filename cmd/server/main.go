package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Placeholder for handling WebSocket connections
	upgrader := websocket.Upgrader{
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: false,
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		println("Error upgrading to websocket:", err.Error())
		return
	}

	defer conn.Close()

	// Handle incoming messages
	for {
		response := r.Host + ": Ping!"
		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
