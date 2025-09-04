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

	"deep_lairs/internal/dbo"
	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

var world = gameobjects.World{
	Places:            make(map[string]*gameobjects.Place),
	CurrentCharacters: 0,
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
	// Initialize world with default characters

	id := uuid.New().String()

	world.Places["hall_of_heroes"].AddCharacter(&gameobjects.Character{
		ID:       uuid.MustParse(id),
		Name:     "Adventurer",
		Location: world.Places["hall_of_heroes"],
	})

	character := gameobjects.GetCharacter(&world, id)

	character.AddKnownLocation(world.Places["hall_of_heroes"])
	// comment stats
	character.Init(
		100, // health
		10,  // attack
		5,   // defense
		20,  // mana
		15,  // stamina
		1,   // speed
		1,   // intelligence
	)

	character.Save()

	character.AddMessage(fmt.Sprintf(protocol.IMAGE, "logo.webp"))
	character.AddMessage(fmt.Sprintf("Character %s connected", character.GetName()))

	g, ctx := errgroup.WithContext(context.Background())

	CharacterJoin(character)

	go character.StartCalcStatsHandler()
	SendCharacterState(character)

	defer conn.Close()

	// Handle outgoing messages
	g.Go(func() error {
		return handleOutgoingMessages(ctx, conn, character)
	})
	// handle incoming messages
	g.Go(func() error {
		return handleIncomingMessages(ctx, conn, character)
	})

	if err := g.Wait(); err != nil {
		log.Println("Error occurred:", err)
	}

}

func handleOutgoingMessages(ctx context.Context, conn *websocket.Conn, character *gameobjects.Character) error {
	for {
		message := ""
		select {
		case <-ctx.Done():
			return nil
		default:
			if len(character.MessageQueue) > 0 {
				message = character.MessageQueue[0]
				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("Error writing message:", err)
					character.Location.AddMessage(fmt.Sprintf("Character %s disconnected", character.GetName()))
					character.Location.RemoveCharacter(character, "poof")
					return err
				} else {
					character.ClearLastMessage()
				}
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func handleIncomingMessages(ctx context.Context, conn *websocket.Conn, character *gameobjects.Character) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				character.Location.AddMessage(fmt.Sprintf("Character %s disconnected", character.GetName()))
				character.Location.RemoveCharacter(character, "poof")
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
				CharacterSay(splitMsg, character)
			case "shout", "sh":
				CharacterShout(splitMsg, character)
			case "look", "l":
				if !character.Looked {
					CharacterLook(splitMsg, character)
					character.Looked = true
				} else {
					CharacterQuickLook(splitMsg, character)
				}
			case "quick_look", "ql":
				CharacterQuickLook(splitMsg, character)
			case "set_name":
				if len(splitMsg) == 2 {
					character.Name = splitMsg[1]
					character.AddMessage(fmt.Sprintf("Name set to: %s", character.GetName()))
				} else {
					character.AddMessage("Usage: set_name <name>")
				}
				SendCharacterState(character)
			case "go", "g":
				knownLocation, err := CharacterGo(splitMsg, character)
				if err != nil {
					log.Println("Error in CharacterGo:", err)
					break
				}
				if !knownLocation {
					CharacterJoin(character)
				} else {
					CharacterWhere(splitMsg, character)
				}
			case "givexp":
				if len(splitMsg) == 2 {
					amount, err := strconv.Atoi(splitMsg[1])
					if err != nil {
						character.AddMessage("Invalid XP amount.")
					} else {
						character.XP += amount
					}
				} else {
					character.AddMessage("Usage: givexp <amount>")
				}
				SendCharacterState(character)
			case "flipcombat":
				character.InCombat = !character.InCombat
				if character.InCombat {
					character.AddMessage("You are now in combat.")
				} else {
					character.AddMessage("You are no longer in combat.")
				}
				SendCharacterState(character)
			case "where", "w":
				CharacterWhere(splitMsg, character)
			case "time", "t":
				character.AddMessage(fmt.Sprintf("Current server time: %s", time.Now().Format(time.RFC1123)))
			case "questboard", "qb":
				CharacterQuestBoard(splitMsg, character)
			case "help":
				character.AddMessage("Available commands: say, look, set_name, help")
			case "lol", "lmao":
				CharacterLaugh(character)
			case "se", "search":
				CharacterSearch(splitMsg, character)
			case "i", "inv", "inventory":
				CharacterInventory(character)
			case "equip", "eq":
				CharacterEquip(splitMsg, character)
			case "unequip", "ue":
				CharacterUnequip(splitMsg, character)
			case "do", "d":
				CharacterDo(splitMsg, character)
			default:
				character.AddMessage(fmt.Sprintf(protocol.I_DONT_KNOW_HOW_TO, firstWord))
			}
		}
	}
}

// entry point for the application
// ad astra!!
func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on", protocol.SERVER_PORT)

	dbo.InitDBO()

	// cors allow all origins
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	loadWorld()

	// Initialise the worlds and set message threads for each place
	for _, place := range world.Places {
		// init currently sets the paths between places
		// but it will also initially create the enemies in the world
		// and NPCs
		place.Init(&world)
		go place.StartMessageHandler()
		go place.StartCheckCharactersHandler()
	}
	go world.StartJingleHandler()

	if err := http.ListenAndServe(protocol.SERVER_PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// loadWorld initializes the game world by loading place data from JSON files.
func loadWorld() {
	// for each json file in ./json folder consume and deserialise to a place
	files, err := os.ReadDir("./json/places/")
	if err != nil {
		log.Println("Error reading json directory:", err)
		panic(err)
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile("./json/places/" + file.Name())
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
}
