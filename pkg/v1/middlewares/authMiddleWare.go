package v1

import (
	"crypto/rsa"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(privateKey *rsa.PrivateKey) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: privateKey.Public(),
			JWTAlg: jwtware.RS256},
	})
}
