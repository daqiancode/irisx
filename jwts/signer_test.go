package jwts_test

import (
	"fmt"
	"testing"

	"github.com/daqiancode/irisx/jwts"

	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/stretchr/testify/assert"
)

func TestHS256(t *testing.T) {
	pk := "123456"
	uid := "123"
	r, err := jwts.GenerateTokenPair(jwts.SignAlg("HS256"), pk, uid, "ADMIN", "xxxxxxxx", 30*60, 24*3600)
	assert.Nil(t, err)
	assert.True(t, len(r.AccessToken) > 0)
	assert.True(t, len(r.RefreshToken) > 0)
	verifier := jwt.NewVerifier(jwt.HS256, pk)
	a, err := verifier.VerifyToken([]byte(r.AccessToken))
	assert.Nil(t, err)
	var token jwts.AccessToken
	a.Claims(&token)
	assert.Equal(t, token.Subject, uid)
}

func TestEdDSA(t *testing.T) {
	publicKey, privateKey, err := jwts.GenerateEdDSAKeyPair()
	assert.Nil(t, err)
	r, err := jwts.GenerateTokenPair(jwts.SignAlg("EdDSA"), privateKey, "123", "ADMIN", "xxxxxxxx", 30*60, 24*3600)
	assert.Nil(t, err)
	fmt.Println(r)
	at, err := jwts.JwtVerify([]byte(r.AccessToken), publicKey)
	assert.Nil(t, err)
	assert.Equal(t, at.StandardClaims.Subject, "123")
}
