package irisx

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/daqiancode/gocommons/commons"
	"github.com/daqiancode/gocommons/commons/states"
	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog"
)

func NewRecover(log zerolog.Logger) context.Handler {

	return func(ctx *context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() { // handled by other middleware.
					return
				}
				switch e := err.(type) {
				case *commons.ServiceError:
					ctx.StatusCode(406)
					log.Error().Err(e).Msg(e.Error())
					ctx.JSON(commons.Result{State: states.ServiceError, Message: e.Message})
				case error:
					ctx.StatusCode(500)
					stack := string(debug.Stack())
					fmt.Println(stack)
					log.Error().Err(err.(error)).Msg(e.Error())
					log.Error().Stack().Err(err.(error)).Msg("")
					ctx.JSON(commons.Result{State: 1, Message: e.Error()})
				default:
					stack := string(debug.Stack())
					fmt.Println(stack)
					log.Error().Stack().Err(errors.New("not error type")).Interface("error", e).Send()
					ctx.StatusCode(500)
					ctx.JSON(commons.Result{State: 1})
				}
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}
