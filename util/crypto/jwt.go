package crypto

import (
	"github.com/golang-jwt/jwt/v4"

	cc "jpc16-core/common"
	"jpc16-core/type/response"
)

func SignJwt(claim jwt.Claims) (string, *response.ErrorInstance) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	if signedToken, err := token.SignedString([]byte(*cc.Config.Secret)); err != nil {
		return "", response.Error(true, "Unable to sign JWT token", err)
	} else {
		return signedToken, nil
	}
}
