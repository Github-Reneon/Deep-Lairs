package main

import (
	"deep_lairs/internal/dbo"
	"deep_lairs/internal/gameobjects"
	"deep_lairs/internal/protocol"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetGame(c *fiber.Ctx) error {
	// set content type to html
	c.Set("Content-Type", "text/html")
	WebSocketURL := protocol.DEV_WS_LINK
	if Prod {
		WebSocketURL = protocol.PROD_WS_LINK
	}

	if c.Cookies("user_id") == "" {
		return c.Redirect("/")
	}

	if c.Cookies("character_id") == "" {
		return c.Redirect("/app/character_select")
	}

	user, err := GetUserInMemFromId(c.Cookies("user_id"))
	if err != nil {
		log.Println(err.Error())
		return c.Redirect("/")
	}

	user.CurrentCharacterId = c.Cookies("character_id")

	return c.Render("game", fiber.Map{
		"Version":      protocol.CLIENT_VERSION,
		"WebSocketURL": WebSocketURL,
	})
}

func GetIndex(c *fiber.Ctx) error {

	LoggedIn := false

	if c.Locals("LoggedIn") != nil {
		LoggedIn = c.Locals("LoggedIn").(bool)
	}

	c.Set("Content-Type", "text/html")
	return c.Render("index", fiber.Map{
		"Version":  protocol.CLIENT_VERSION,
		"LoggedIn": LoggedIn,
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

func PostSignup(c *fiber.Ctx) error {

	log.Println("Signup request received")

	if err := c.BodyParser(&struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("signup", fiber.Map{
			"error": "Invalid request body",
		})
	}

	userName := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	if userName == "" || password == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).Render("signup", fiber.Map{
			"error": "All fields are required",
		})
	}

	hashedPassword := HashPassword(password)
	if hashedPassword == "" {
		return c.Status(fiber.StatusInternalServerError).Render("signup", fiber.Map{
			"error": "Failed to hash password",
		})
	}
	password = hashedPassword

	log.Println("Creating user:", userName)

	// check if user already exists in memory
	if FindUserMem(userName) {
		return c.Status(fiber.StatusBadRequest).Render("signup", fiber.Map{
			"error": "User already exists",
		})
	}

	// check if user already exists in database
	if _, err := dbo.LoadUser(userName); err == nil {
		return c.Status(fiber.StatusBadRequest).Render("signup", fiber.Map{
			"error": "User already exists",
		})
	}

	// create user in database
	if err := dbo.CreateUser(userName, password, email); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("signup", fiber.Map{
			"error": "Failed to create user",
		})
	}

	// load user into memory
	if !PutUserInMem(userName) {
		return c.Status(fiber.StatusInternalServerError).Render("signup", fiber.Map{
			"error": "Failed to load user into memory",
		})
	}

	// confirm user created
	user, err := GetUserInMemFromName(userName)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("User created:", userName)

	c.Cookie(&fiber.Cookie{
		Name:  "user_id",
		Value: user.ID,
	})

	// redirect to character select
	return c.Redirect("/app/character_select")
}

func PostLogin(c *fiber.Ctx) error {
	if err := c.BodyParser(&struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"error": "Invalid request body",
		})
	}

	userName := c.FormValue("username")
	password := c.FormValue("password")

	if userName == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"error": "Username and password are required",
		})
	}

	log.Println("Login request received for user:", userName)

	userMem, err := GetUserInMemFromName(userName)
	if err != nil {
		user, err := dbo.LoadUser(userName)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).Render("login", fiber.Map{
				"error": "Invalid username or password",
			})
		}
		PutUserInMemFromUserObj(user)
		userMem, err = GetUserInMemFromName(user.Username)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).Render("login", fiber.Map{
				"error": "Invalid username or password",
			})
		}
	}

	if !CheckPasswordHash(password, userMem.Password) {
		return c.Status(fiber.StatusUnauthorized).Render("login", fiber.Map{
			"error": "Invalid username or password",
		})
	}

	log.Println("User logged in:", userName)

	c.Cookie(&fiber.Cookie{
		Name:  protocol.COOKIE_USER_ID,
		Value: userMem.ID,
	})

	return c.Redirect("/app/character_select")
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
		"Version":    protocol.CLIENT_VERSION,
		"Characters": []gameobjects.Character{},
	})
}

func PostLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    protocol.COOKIE_USER_ID,
		Expires: time.Now(),
		Value:   "",
	})
	return c.Redirect("/")
}
