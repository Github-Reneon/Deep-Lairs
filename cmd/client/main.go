package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {

	engine := html.New("./views", ".html")

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

	app.Get("/", GetIndex)

	if err := app.Listen(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
