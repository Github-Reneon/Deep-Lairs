package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	EnableCompression: true,
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

	world.Places["tavern"].AddUser(&gameobjects.User{
		ID:       id,
		Name:     "Reneon",
		Location: world.Places["tavern"],
	})

	user := gameobjects.GetUser(&world, id)
	// health int, attack int, defense int, mana int, stamina int
	user.InitFightable(
		7, // health
		5, // attack
		5, // defense
		1, // mana
		1, // stamina
		3, // speed
		1, // intelligence
	)

	user.AddKnownLocation(world.Places["tavern"])

	user.AddMessage(fmt.Sprintf("User %s connected", user.GetName()))
	g, ctx := errgroup.WithContext(context.Background())

	UserJoin(user)

	go user.StartCalcStatsHandler()
	SendUserState(user)

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
					user.Location.RemoveUser(user, "poof")
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
				user.Location.RemoveUser(user, "poof")
				return err
			}
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
				if !user.Looked {
					UserLook(splitMsg, user)
					user.Looked = true
				} else {
					UserQuickLook(splitMsg, user)
				}
			case "quick_look", "ql":
				UserQuickLook(splitMsg, user)
			case "set_name":
				if len(splitMsg) == 2 {
					user.Name = splitMsg[1]
					user.AddMessage(fmt.Sprintf("Name set to: %s", user.GetName()))
				} else {
					user.AddMessage("Usage: set_name <name>")
				}
				SendUserState(user)
			case "go", "g":
				knownLocation, err := UserGo(splitMsg, user)
				if err != nil {
					log.Println("Error in UserGo:", err)
					break
				}
				if !knownLocation {
					UserJoin(user)
				} else {
					UserWhere(splitMsg, user)
				}
			case "givexp":
				if len(splitMsg) == 2 {
					amount, err := strconv.Atoi(splitMsg[1])
					if err != nil {
						user.AddMessage("Invalid XP amount.")
					} else {
						user.XP += amount
					}
				} else {
					user.AddMessage("Usage: givexp <amount>")
				}
				SendUserState(user)
			case "flipcombat":
				user.InCombat = !user.InCombat
				if user.InCombat {
					user.AddMessage("You are now in combat.")
				} else {
					user.AddMessage("You are no longer in combat.")
				}
				SendUserState(user)
			case "where", "w":
				UserWhere(splitMsg, user)
			case "time", "t":
				user.AddMessage(fmt.Sprintf("Current server time: %s", time.Now().Format(time.RFC1123)))
			case "questboard", "qb":
				UserQuestBoard(splitMsg, user)
			case "help":
				user.AddMessage("Available commands: say, look, set_name, help")
			case "lol", "lmao":
				UserLaugh(user)
			case "se", "search":
				UserSearch(splitMsg, user)
			case "i", "inv", "inventory":
				UserInventory(user)
			case "equip", "eq":
				UserEquip(splitMsg, user)
			case "unequip", "ue":
				UserUnequip(splitMsg, user)
			default:
				user.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, firstWord))
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on", protocol.PORT)

	// for each json file in ./json folder consume and deserialise to a place
	files, err := os.ReadDir("./json")
	if err != nil {
		log.Println("Error reading json directory:", err)
		return
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile("./json/" + file.Name())
			if err != nil {
				log.Println("Error reading json file:", err)
				continue
			}
			var place gameobjects.Place
			if err := json.Unmarshal(data, &place); err != nil {
				log.Println("Error unmarshalling json:", err)
				continue
			}
			log.Println(place.ID, place.Name, place.JoiningLocationIds)
			world.Places[place.ID] = &place
		}
	}

	// Set message threads for each place
	for _, place := range world.Places {
		place.Init(&world)
		go place.StartMessageHandler()
		go place.StartCheckUsersHandler()
	}
	go world.StartJingleHandler()

	if err := http.ListenAndServe(protocol.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
