package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"deep_lairs/internal/protocol"
	"deep_lairs/internal/user"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin:       func(r *http.Request) bool { return true },
	EnableCompression: false,
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Placeholder for handling WebSocket connections
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading to websocket:", err.Error())
		return
	}

	user := user.User{
		ID:           uuid.New().String(),
		MessageQueue: make([]string, 0),
	}

	user.MessageQueue = append(user.MessageQueue, fmt.Sprintf("User %s connected", user.ID))
	debugIncomingMessage := fmt.Sprintf(protocol.LOOK_MESSAGE, "You see a tavern wench", "drinks.webp")
	user.MessageQueue = append(user.MessageQueue, debugIncomingMessage)

	wg := sync.WaitGroup{}
	defer conn.Close()
	wg.Add(1)
	// Handle outgoing messages
	go handleOutgoingMessages(conn, &wg, &user)
	// handle incoming messages
	go handleIncomingMessages(conn, &wg, &user)
	wg.Wait()
}

func handleOutgoingMessages(conn *websocket.Conn, wg *sync.WaitGroup, user *user.User) {
	defer wg.Done()
	for {
		message := ""
		if len(user.MessageQueue) > 0 {
			message = user.MessageQueue[0]
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Error writing message:", err)
				break
			} else {
				log.Printf("%s: Message sent successfully\n", user.ID)
				user.MessageQueue = user.MessageQueue[1:]
			}
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func handleIncomingMessages(conn *websocket.Conn, wg *sync.WaitGroup, user *user.User) {
	defer wg.Done()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message: %s\n", msg)

		// basic parse message replace soon
		// split on spaces
		if string(msg) == "" {
			continue
		}
		splitMsg := strings.Split(string(msg), " ")

		firstWord := strings.ToLower(splitMsg[0])

		switch firstWord {
		case "say", "s":
			user.MessageQueue = append(user.MessageQueue, fmt.Sprintf(protocol.SAY, user.GetName(), strings.Join(splitMsg[1:], " ")))
		case "look", "l":
			user.MessageQueue = append(user.MessageQueue, fmt.Sprintf(protocol.LOOK_MESSAGE, "You see a tavern wench", "drinks.webp"))
		case "set_name":
			if len(splitMsg) == 2 {
				user.Name = splitMsg[1]
				user.MessageQueue = append(user.MessageQueue, fmt.Sprintf("Name set to: %s", user.GetName()))
			} else {
				user.MessageQueue = append(user.MessageQueue, "Usage: set_name <name>")
			}
		case "help":
			user.MessageQueue = append(user.MessageQueue, "Available commands: say, look, set_name, help")
		default:
			user.MessageQueue = append(user.MessageQueue, fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, firstWord))
		}

		// add to user message queue message
		log.Printf("%s: Message added to queue\n", user.ID)
		log.Println(user.MessageQueue)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
