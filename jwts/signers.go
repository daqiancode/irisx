package jwts

import (
	"errors"
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"
)

// type Signer struct {
// 	alg             jwt.Alg
// 	privateKey      ed25519.PrivateKey
// 	signer          *jwt.Signer
// 	tokenTtl        time.Duration
// 	refreshTokenTtl time.Duration
// }

func Sign(alg SignAlg, claim interface{}, privateKey string) (string, error) {
	a := alg.Alg()
	if a == nil {
		return "", errors.New("not support such signature algorithm: " + string(alg))
	}
	pk, err := alg.ParsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	bs, err := jwt.NewSigner(a, pk, 0).Sign(claim)
	return string(bs), err
}

// func SignHS256(claim interface{}, privateKey string) (string, error) {
// 	bs, err := jwt.NewSigner(jwt.HS256, privateKey, 0).Sign(claim)
// 	return string(bs), err
// }

// // SignEdDSA ,Please use GenerateEdDSAKeyPair to generate key pair
// func SignEdDSA(claim interface{}, privateKey string) (string, error) {
// 	pk, err := jwt.ParsePrivateKeyEdDSA([]byte(processPrivateKey(privateKey)))
// 	if err != nil {
// 		return "", err
// 	}
// 	bs, err := jwt.NewSigner(jwt.EdDSA, pk, 0).Sign(claim)
// 	return string(bs), err
// }

// func processPrivateKey(privateKey string) string {
// 	privateKey = strings.TrimSpace(privateKey)
// 	if !strings.HasPrefix(privateKey, "---") {
// 		return "-----BEGIN PRIVATE KEY-----\n" + privateKey + "\n-----END PRIVATE KEY-----"
// 	}
// 	return privateKey
// }

type TokenPair struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

var UUIDLen = 20

// GenerateTokenPair generate access_token & refresh_token, roles are seperate by space .eg, "ADMIN SYSTEM"
func GenerateTokenPair(alg SignAlg, privateKey, uid, roles, encryptedPassword string, accessTokenMaxAge, refreshTokenMaxAge int64) (TokenPair, error) {
	refreshTokenID := UUID(UUIDLen)
	now := time.Now().Unix()

	accessToken := AccessToken{Claims: jwt.Claims{ID: UUID(UUIDLen), Subject: uid, Expiry: now + accessTokenMaxAge}, Roles: roles, Rid: refreshTokenID}
	refreshToken := &RefreshToken{Claims: jwt.Claims{ID: refreshTokenID, Subject: uid, Expiry: now + accessTokenMaxAge}}
	refreshToken.V = CreateRefreshTokenVerfication(encryptedPassword, uid, roles, refreshToken.Expiry)
	jwtAt, err := Sign(alg, accessToken, privateKey)
	if err != nil {
		return TokenPair{}, err
	}
	jwtRt, err := Sign(alg, refreshToken, privateKey)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: jwtAt, RefreshToken: jwtRt}, nil

}

func GenerateAccessToken(alg SignAlg, privateKey, uid, roles string, accessTokenMaxAge int64, refreshTokenID string) (string, error) {
	accessToken := AccessToken{Claims: jwt.Claims{ID: UUID(UUIDLen), Subject: uid, Expiry: time.Now().Unix() + accessTokenMaxAge}, Roles: roles, Rid: refreshTokenID}
	return Sign(alg, accessToken, privateKey)
}
