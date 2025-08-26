package main

import (
	"deep_lairs/internal/protocol"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {

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
