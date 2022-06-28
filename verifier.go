package irisx

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/middleware/jwt/blocklist/redis"
)

type Verifier struct {
	verifier *jwt.Verifier
}

func processPublicKey(publicKey string) string {
	publicKey = strings.TrimSpace(publicKey)
	if !strings.HasPrefix(publicKey, "---") {
		return "-----BEGIN PUBLIC KEY-----\n" + publicKey + "\n-----END PUBLIC KEY-----"
	}
	return publicKey
}

//NewVerifier please use ECDSA alg
func NewVerifier(publicKey string, alg jwt.Alg) *Verifier {
	pk, err := jwt.ParsePublicKeyECDSA([]byte(processPublicKey(publicKey)))
	if err != nil {
		panic(err)
	}
	return &Verifier{
		verifier: jwt.NewVerifier(alg, pk),
	}
}

func (v *Verifier) GetVerifier() *jwt.Verifier {
	return v.verifier
}

func (v *Verifier) SetBlockList(blocklist *redis.Blocklist) {
	v.verifier.Blocklist = blocklist
}

func (v *Verifier) GetMiddleware(validators ...jwt.TokenValidator) context.Handler {
	return v.verifier.Verify(func() interface{} {
		return new(RbacClaims)
	}, validators...)
}

func findRole(giveRoles, myRoles []string) bool {
	if len(giveRoles) == 0 {
		return true
	}
	giveRolesMap := make(map[string]bool, len(giveRoles))
	for _, v := range giveRoles {
		giveRolesMap[v] = true
	}
	for _, v := range myRoles {
		if giveRolesMap[v] {
			return true
		}
	}
	return false
}

type RbacMiddleware struct {
	roles []string
}

func NewRbacMiddleware(roles []string) *RbacMiddleware {
	return &RbacMiddleware{
		roles: roles,
	}
}

func (s *RbacMiddleware) Middleware(ctx iris.Context) {
	claims := jwt.Get(ctx).(*RbacClaims)
	roles := strings.Split(claims.Roles, " ")
	if !findRole(s.roles, roles) {
		ctx.StopWithJSON(iris.StatusForbidden, "Forbidden")
	}
	ctx.Next()
}
