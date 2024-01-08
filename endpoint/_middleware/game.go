package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"jpc16-core/type/collection"
	"jpc16-core/type/response"
)

func GameMiddleware(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	headerSegment := strings.Split(header, " ")
	if len(headerSegment) != 2 {
		return response.Error(false, "Invalid authorization header format")
	}
	if headerSegment[0] != "Bearer" {
		return response.Error(false, "Invalid authorization type other than bearer")
	}
	bearer := headerSegment[1]
	bearerSegment := strings.Split(bearer, ".")
	if len(bearerSegment) != 2 {
		return response.Error(false, "Invalid Bearer format")
	}

	// * Parse bearer
	gameId, err := primitive.ObjectIDFromHex(bearerSegment[0])
	if err != nil {
		return response.Error(false, "Unable to parse game id", err)
	}
	token := bearerSegment[1]

	// * Find game
	game := new(collection.Game)
	if err := game.Collection().FindByID(gameId, game); err != nil {
		return response.Error(false, "Unable to find game", err)
	}

	// * Validate token
	if *game.Token != token {
		return response.Error(false, "Invalid game token")
	}

	// * Set game to locals
	c.Locals("game", game)

	// * Next
	return c.Next()
}
