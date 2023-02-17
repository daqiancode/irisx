# irisx
extends iris



## Installation
```go
go get -u github.com/daqiancode/irisx@main

```

## Example

```go
import (
	"github.com/daqiancode/irisx"
	"github.com/daqiancode/irisx/jwts"
)

func setupDependencies(app *iris.Application) {
	app.RegisterDependency(func(ctx iris.Context) irisx.Context {
		return irisx.Context{Context: ctx}
	})
	app.RegisterDependency(new(service.Users))
}

func setupControllers(app *iris.Application) {
	accessTokenSetter := jwts.AccessTokenSetter(jwts.AccessTokenSetterConfig{PublicKey: config.Getenv("JWT_PUBLIC_KEY")})
	api := mvcApp.Party(config.PREFIX+"/api", accessTokenSetter)
	api.Party("/major", jwts.Require(),jwts.RBAC([]string{"USER"})).Handle(new(UserController))
}



type UserController struct {
	Ctx              irisx.Context
	Users            *service.Users
}
func (c *UserController) Get() {
	uid := c.Ctx.GetUID() // get user id from access_token
	c.Ctx.OK(c.Users.Get(uid))
}

```