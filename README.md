# irisx
extends iris



## Installation
```go
go get github.com/daqiancode/irisx

```

## Example

```go
v := irisx.NewVerifier(env.Getenv("jwt_public_key"), jwt.EdDSA)
mvcApp := mvc.New(app).Party(env.Getenv("app_prefix", "/"))
mvcApp.Party("/").Handle(new(controller.PublicController))
mvcApp.Party("/user", v.Middleware(), irisx.NewRbacMiddleware("USER").Middleware).Handle(new(controller.UserController))
mvcApp.Party("/admin/user", v.Middleware(), irisx.NewRbacMiddleware("ADMIN").Middleware).Handle(new(controller.AdminUserController))

```