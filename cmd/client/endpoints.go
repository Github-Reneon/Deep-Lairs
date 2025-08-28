package main

import (
	"deep_lairs/internal/protocol"

	"github.com/gofiber/fiber/v2"
)

func GetGame(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	WebSocketURL := protocol.DEV_WS_LINK
	if Prod {
		WebSocketURL = protocol.PROD_WS_LINK
	}
	return c.Render("index", fiber.Map{
		"Version":      "0.1.0",
		"WebSocketURL": WebSocketURL,
	})
}
