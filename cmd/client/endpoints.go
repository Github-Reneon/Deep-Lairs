package main

import (
	"deep_lairs/internal/gameobjects"
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

	if c.Cookies("character_id", protocol.USER_NOT_FOUND) == protocol.USER_NOT_FOUND {
		return c.Redirect("/app/character_select")
	}

	return c.Render("game", fiber.Map{
		"Version":      protocol.CLIENT_VERSION,
		"WebSocketURL": WebSocketURL,
	})
}

func GetIndex(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("index", fiber.Map{
		"Version": protocol.CLIENT_VERSION,
	})
}

func GetLogin(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("login", fiber.Map{
		"Version": protocol.CLIENT_VERSION,
	})
}

func GetSignup(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("signup", fiber.Map{
		"Version": protocol.CLIENT_VERSION,
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
		"Version": protocol.CLIENT_VERSION,
	})
}

func GetCharacterSelect(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.Render("character_select", fiber.Map{
		"Version": protocol.CLIENT_VERSION,
		"Characters": []gameobjects.Character{
			{
				ID:   "char1",
				Name: "Hero123",
				UserFightable: gameobjects.UserFightable{
					Level: 1,
					Class: "Mage",
					Fightable: gameobjects.Fightable{
						Image: "portrait_elf_8.webp",
					},
				},
			},
			{
				ID:   "char2",
				Name: "Warrior456",
				UserFightable: gameobjects.UserFightable{
					Class: "Warrior",
					Level: 5,
					Fightable: gameobjects.Fightable{
						Image: "portrait_dwarf_1.webp",
					},
				},
			},
		},
	})
}
