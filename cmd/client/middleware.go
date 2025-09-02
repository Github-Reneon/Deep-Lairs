package main

import (
	"deep_lairs/internal/protocol"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(c *fiber.Ctx) error {
	// Check for the presence of the "Authorization" header
	/* add back later
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// redirect to login page
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Please log in to access this page")
	}
	*/
	return c.Next()
}

func AlreadyAuth(c *fiber.Ctx) error {
	// Check for the presence of the "Authorization" header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// redirect to game page
		return c.Redirect("/app/game")
	}
	return c.Next()
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	// Retrieve the custom status code if it's a *fiber.Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	// Send custom error page
	err = c.Status(code).Render("error", fiber.Map{
		"Error":      err.Error(),
		"StatusCode": code,
		"Version":    protocol.CLIENT_VERSION,
	})

	if err != nil {
		// In case the render fails
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	// Return from handler
	return nil
}
