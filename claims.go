package irisx

import "github.com/kataras/jwt"

// subject(sub) is user id
type RbacClaims struct {
	jwt.Claims
	// roles seperate by space. e.g., "USER ADMIN"
	Roles string `json:"roles,omitempty"`
}
