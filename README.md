# irisx
extends iris



## Installation
```go
go get -u github.com/daqiancode/irisx@main

```

## Example

```go
func setupDependencies(app *iris.Application) {
	app.RegisterDependency(func(ctx iris.Context) irisx.Context {
		return irisx.Context{Context: ctx}
	})
	app.RegisterDependency(new(service.Users))
}



type UserController struct {
	Ctx              irisx.Context
	Users            *service.Users
	SignupEmailCodes *service.SignupEmailCodes
}
func (c *UserController) Get() {
	uid := c.Ctx.GetUID()
	c.Ctx.OK(c.Users.Get(uid))
}

```