package main

import (
	"github.com/gofiber/fiber/v2"
)

func GetIndex(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("index", fiber.Map{
		"Version": "0.1.0",
	})
}
