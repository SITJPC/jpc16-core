package playEndpoint

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"jpc16-core/common/mng"
	"jpc16-core/type/collection"
	"jpc16-core/type/misc"
	"jpc16-core/type/payload"
	"jpc16-core/type/response"
	"jpc16-core/util/crypto"
	"jpc16-core/util/text"
)

// HandleEnterPin
// @ID pinEnter
// @Summary Pin Enter
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {object} response.InfoResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /play/pin [post]
func HandleEnterPin(c *fiber.Ctx) error {
	// * Parse context
	ct := c.Locals("ct").(context.Context)

	// * Parse body
	body := new(payload.PinBody)
	if err := c.BodyParser(body); err != nil {
		return response.Error(true, "Unable to parse body", err)
	}

	// * Validate body
	if err := text.Validate(body); err != nil {
		return err
	}

	// * Find player by pin
	player := new(collection.Player)
	if err := mng.PlayerCollection.First(bson.M{"pin": body.Pin}, player); err != nil {
		return response.Error(true, "Unable to find player", err)
	}
	ct = context.WithValue(ct, "player", player)

	// * Construct JWT claim
	claim := &misc.PlayerClaim{
		Id: player.ID,
	}

	// * Sign claim
	token, err := crypto.SignJwt(claim)
	if err != nil {
		return err
	}

	// * Set cookie
	c.Cookie(&fiber.Cookie{
		Name:  "playerToken",
		Value: token,
	})

	// * Return response
	return c.JSON(response.Info(&payload.PinTokenResponse{
		Token: &token,
	}))
}
