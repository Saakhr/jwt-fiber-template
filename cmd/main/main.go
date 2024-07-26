package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"

	v1routes "github.com/Saakhr/jwt-fiber-template/pkg/v1/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	privateKey *rsa.PrivateKey
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()

	// Login route
	err = getKey()
	if err != nil {
		log.Fatal("Couldn't Load JWT RSA key" + err.Error())
	}
	app.Use(logger.New())
	v1 := v1routes.GetRoutes(privateKey)
	app.Mount("/v1", v1)

	log.Fatal(app.Listen(":8080"))
}
func getKey() error {
	// Replace with your actual key file path
	// keyFilePath := "key.pem"
	//
	// // Read the key file
	// file, err := os.Open(keyFilePath)
	// if err != nil {
	// 	return errors.New("Error opening file:" + keyFilePath)
	// }
	// defer file.Close()
	//
	// // Read the file contents
	// _, err = io.ReadAll(file)
	// if err != nil {
	// 	return errors.New("Error reading file:" + keyFilePath)
	// }

	xs := os.Getenv("JWT_RS_SECRET")

	// Decode the PEM block
	block, _ := pem.Decode([]byte(xs))
	if block == nil {
		return errors.New("Error decoding PEM block")
	}

	// Parse the private key
	privateKeys, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("Error parsing private key")
	}
	var ok bool
	privateKey, ok = privateKeys.(*rsa.PrivateKey)
	if !ok {
		return errors.New("Error: parsed key is not an RSA private key")
	}
	return nil
}
