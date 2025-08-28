package main

import (
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
			Views: engine,
		},
	)

	app.Static("/", "./static", fiber.Static{
		Compress: true,
		Browse:   false,
		Download: false,
	})

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

	app.Get("/", GetIndex)

	if err := app.Listen(protocol.CLIENT_PORT); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
