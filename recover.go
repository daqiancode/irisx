package irisx

import (
	"fmt"
	"runtime/debug"

	"github.com/kataras/iris/v12/context"
)

func NewRecover() context.Handler {

	return func(ctx *context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() { // handled by other middleware.
					return
				}
				switch e := err.(type) {

				case error:
					ctx.StatusCode(500)
					stack := string(debug.Stack())
					fmt.Println(stack)
					ctx.JSON(Result{State: 1, Error: e.Error()})
				default:
					stack := string(debug.Stack())
					fmt.Println(stack)
					ctx.StatusCode(500)
					ctx.JSON(Result{State: 1})
				}
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}
