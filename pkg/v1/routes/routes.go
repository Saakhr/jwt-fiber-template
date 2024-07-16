package v1

import (
	"crypto/rsa"
	"time"

	v1middlewares "github.com/Saakhr/jwt-fiber-template/pkg/v1/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
)

func GetRoutes(Key *rsa.PrivateKey) *fiber.App {
	v1 := fiber.New()
	privateKey = Key

	// Unauthenticated route
	v1.Post("/login", login)
	v1.Get("/", accessible)

	// Restricted Routes
	v1.Get("/2", accessible)
	v1.Get("/restricted", v1middlewares.NewAuthMiddleware(privateKey), restricted)

	return v1
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(privateKey)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
