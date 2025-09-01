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
	return c.Render("game", fiber.Map{
		"Version":      "0.1.0",
		"WebSocketURL": WebSocketURL,
	})
}

func GetIndex(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("index", fiber.Map{
		"Version": "0.1.0",
	})
}

func GetLogin(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("login", fiber.Map{
		"Version": "0.1.0",
	})
}

func GetSignup(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("signup", fiber.Map{
		"Version": "0.1.0",
	})
}

func PostLogin(c *fiber.Ctx) error {
	if err := c.BodyParser(&struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	c.Set("Content-Type", "application/json")
	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func GetCharacterCreation(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("character_creation", fiber.Map{
		"Version": "0.1.0",
	})
}
