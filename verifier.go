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
	pk, err := jwt.ParsePublicKeyEdDSA([]byte(processPublicKey(publicKey)))
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

func (v *Verifier) Middleware(validators ...jwt.TokenValidator) context.Handler {
	return v.verifier.Verify(func() interface{} {
		return new(RbacClaims)
	}, validators...)
}

//findRole: ([user,admin],[user])->OK ,([seller*],[seller:sales]) -> OK
func findRole(giveRoles, myRoles []string) bool {
	if len(giveRoles) == 0 {
		return true
	}
	var prefixRoles []string
	giveRolesMap := make(map[string]bool, len(giveRoles))
	for _, v := range giveRoles {
		giveRolesMap[v] = true
		if v[len(v)-1] == '*' {
			prefixRoles = append(prefixRoles, v[0:len(v)-1])
		}
	}
	for _, v := range myRoles {
		if giveRolesMap[v] {
			return true
		}
	}
	if len(prefixRoles) > 0 {
		for _, x := range myRoles {
			for _, v := range prefixRoles {
				if strings.HasPrefix(x, v) {
					return true
				}
			}
		}
	}
	return false
}

type RbacMiddleware struct {
	roles []string
}

func NewRbacMiddleware(roles ...string) *RbacMiddleware {
	return &RbacMiddleware{
		roles: roles,
	}
}

func (s *RbacMiddleware) Middleware(ctx iris.Context) {
	claims := jwt.Get(ctx).(*RbacClaims)
	roles := strings.Split(claims.Roles, " ")
	if !findRole(s.roles, roles) {
		ctx.StopWithText(iris.StatusForbidden, "Forbidden")
		return
	}
	ctx.Next()
}
