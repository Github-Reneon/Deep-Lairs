package main

import (
	"deep_lairs/internal/dbo"
	"deep_lairs/internal/protocol"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var Prod = false

func main() {

	// if --prod then prod = true
	if len(os.Args) > 1 && os.Args[1] == "--prod" {
		Prod = true
		fmt.Println("Running in production mode")
	} else {
		fmt.Println("Running in development mode")
	}

	engine := html.New("./views", ".html")

	engine.Reload(true)

	app := fiber.New(
		fiber.Config{
			Views:        engine,
			ErrorHandler: ErrorHandler,
		},
	)

	app.Static("/", "./static", fiber.Static{
		Compress: true,
		Browse:   false,
		Download: false,
	})

	dbo.InitDBO()
	setCORS(app)
	setUpRoutes(app)

	if err := app.Listen(protocol.CLIENT_PORT); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func setCORS(app *fiber.App) {
	// cors allow all origins
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})
}

func setUpRoutes(app *fiber.App) {

	app.Get("/", GetIndex)
	app.Get("/index", GetIndex)
	app.Get("/login", func(c *fiber.Ctx) error { return c.Redirect("/auth/login") })
	app.Get("/signup", func(c *fiber.Ctx) error { return c.Redirect("/auth/signup") })
	app.Post("/login", func(c *fiber.Ctx) error { return c.Redirect("/auth/login") })
	app.Post("/signup", func(c *fiber.Ctx) error { return c.Redirect("/auth/signup") })

	auth := app.Group("/auth")
	auth.Get("/login", GetLogin)
	auth.Post("/login", PostLogin)
	auth.Post("/signup", PostSignup)
	auth.Get("/signup", GetSignup)
	auth.Use(GetLoggedIn)

	game := app.Group("/app")
	game.Use(AuthRequired)
	game.Get("/game", GetGame)
	game.Get("/character_creation", GetCharacterCreation)
	game.Get("/character_select", GetCharacterSelect)
	game.Use(GetLoggedIn)
}
