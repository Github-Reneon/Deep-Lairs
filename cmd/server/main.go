package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

var world = gameobjects.World{
	Places:       make(map[string]*gameobjects.Place),
	CurrentUsers: 0,
}

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
	// Initialize world with default users

	id := uuid.New().String()

	world.Places["tavern"].Users[id] = &gameobjects.User{
		ID:       id,
		Name:     "Reneon",
		Location: world.Places["tavern"],
	}

	user := gameobjects.GetUser(&world, id)

	user.AddMessage(fmt.Sprintf("User %s connected", user.GetName()))
	g, ctx := errgroup.WithContext(context.Background())

	user.AddMessage(fmt.Sprintf(protocol.YOU_ARE_IN, user.Location.Name, user.Location.LocationImage))

	defer conn.Close()

	// Handle outgoing messages
	g.Go(func() error {
		return handleOutgoingMessages(ctx, conn, user)
	})
	// handle incoming messages
	g.Go(func() error {
		return handleIncomingMessages(ctx, conn, user)
	})

	if err := g.Wait(); err != nil {
		log.Println("Error occurred:", err)
	}

}

func handleOutgoingMessages(ctx context.Context, conn *websocket.Conn, user *gameobjects.User) error {
	for {
		message := ""
		select {
		case <-ctx.Done():
			return nil
		default:
			if len(user.MessageQueue) > 0 {
				message = user.MessageQueue[0]
				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("Error writing message:", err)
					user.Location.AddMessage(fmt.Sprintf("User %s disconnected", user.GetName()))
					user.Location.RemoveUser(user)
					return err
				} else {
					user.ClearLastMessage()
				}
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func handleIncomingMessages(ctx context.Context, conn *websocket.Conn, user *gameobjects.User) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				user.Location.AddMessage(fmt.Sprintf("User %s disconnected", user.GetName()))
				user.Location.RemoveUser(user)
				return err
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
				UserSay(splitMsg, user)
			case "shout", "sh":
				UserShout(splitMsg, user)
			case "look", "l":
				UserLook(splitMsg, user)
			case "quick_look", "ql":
				user.AddMessage(fmt.Sprintf(protocol.LOOK_NO_IMAGE, user.Location.QuickLook))
			case "set_name":
				if len(splitMsg) == 2 {
					user.Name = splitMsg[1]
					user.AddMessage(fmt.Sprintf("Name set to: %s", user.GetName()))
				} else {
					user.AddMessage("Usage: set_name <name>")
				}
			case "where", "w":
				user.AddMessage(fmt.Sprintf("You are in %s<br>%s", user.Location.Name, user.Location.Description))
			case "time", "t":
				user.AddMessage(fmt.Sprintf("Current server time: %s", time.Now().Format(time.RFC1123)))
			case "help":
				user.AddMessage("Available commands: say, look, set_name, help")
			case "lol", "lmao":
				UserLaugh(user)
			default:
				user.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, firstWord))
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on :3000")

	// Initialize world with default places
	world.Places["tavern"] = InitPlace()
	go world.Places["tavern"].StartMessageHandler()
	go world.StartJingleHandler()

	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
