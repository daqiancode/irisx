package irisx

import (
	"runtime/debug"
	"strings"

	"github.com/daqiancode/jsoniter"
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Context struct {
	iris.Context
}

var JSON = jsoniter.Decapitalized

func (c *Context) ReadJSON(outPtr interface{}) error {
	body, restoreBody, err := context.GetBody(c.Request(), true)
	if err != nil {
		return err
	}
	restoreBody()
	err = JSON.Unmarshal(body, outPtr)
	if err != nil {
		return err
	}
	return c.Application().Validate(outPtr)
}
func (c *Context) ReadQuery(ptr interface{}) error {
	values := c.Request().URL.Query()
	if len(values) == 0 {
		if c.Application().ConfigurationReadOnly().GetFireEmptyFormError() {
			return context.ErrEmptyForm
		}
		return nil
	}

	err := schema.DecodeQuery(values, ptr)
	if err != nil {
		return err
	}

	return c.Application().Validate(ptr)
}

func (c *Context) ReadForm(formObject interface{}) error {
	values := c.FormValues()
	if len(values) == 0 {
		if c.Application().ConfigurationReadOnly().GetFireEmptyFormError() {
			return context.ErrEmptyForm
		}
		return nil
	}
	for k, v := range values {
		values[k] = v
	}

	err := schema.DecodeForm(values, formObject)
	if err != nil {
		return err
	}
	return c.Application().Validate(formObject)
}

func (c *Context) JSON(v interface{}) error {
	bs, err := JSON.Marshal(v)
	if err != nil {
		return err
	}
	c.ContentType(context.ContentJSONHeaderValue)
	_, err = c.Write(bs)
	return err
}

func (c *Context) Finish(data interface{}, err error) error {
	if err != nil {
		c.Error(err, 500)
		return err
	}
	return c.OK(data)
}

func (c *Context) Page(items interface{}, pageIndex, pageSize, total int, err error) error {
	if err != nil {
		c.Error(err, 500)
		return err
	}
	return c.OK(NewPage(items, pageIndex, pageSize, total))
}

func (c *Context) OK(data interface{}) error {
	return c.JSON(Result{Data: data})
}

func (c *Context) Fail(message string, state, httpStatus int) error {
	c.StatusCode(httpStatus)
	return c.JSON(Result{State: state, ErrorCode: message})
}

// server internal error
func (c *Context) FailInternal(message string, state int) error {
	return c.Fail(message, state, 500)
}

// bussiness logic error
func (c *Context) FailService(message string, state int) error {
	return c.Fail(message, state, 406)
}

// request parameter error
func (c *Context) FailParams(fieldErrors map[string]string) error {
	c.StatusCode(400)
	return c.JSON(Result{State: 1, Error: "request parameter error", FieldErrors: fieldErrors})
}

// type stackTracer interface {
// 	StackTrace() errors.StackTrace
// }

func (c *Context) Error(err error, statusCode int) error {
	if err == nil {
		return c.OK(nil)
	}
	r := Result{State: 1, Error: err.Error()}
	if v, ok := err.(FieldErrorsGetter); ok {
		statusCode = 400
		r.FieldErrors = v.GetFieldErrors()
	}
	if v, ok := err.(StateGetter); ok {
		r.State = v.GetState()
	}
	if v, ok := err.(ErrorCodeGetter); ok {
		r.ErrorCode = v.GetErrorCode()
	}
	if v, ok := err.(HttpStatusCodeGetter); ok {
		statusCode = v.GetHttpStatusCode()
	}
	c.Application().Logger().Error(err)
	if statusCode >= 500 {
		c.Application().Logger().Error(string(debug.Stack()))
	}
	c.StatusCode(statusCode)
	return c.JSON(r)
}

// request parameter error
func (c *Context) ErrorParam(err error) error {
	if es, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]string, len(es))
		for _, v := range es {
			fieldErrors[v.Field()] = v.ActualTag()
			// fieldErrors[v.Field()] = v.Error()
		}
		return c.FailParams(fieldErrors)
	}
	return c.Error(err, 406)
}

func (c *Context) GetIP() string {
	ip := c.GetHeader("X-Real-Ip")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
		if ip != "" {
			return strings.TrimSpace(strings.Split(ip, ",")[0])
		}
	}
	return c.Context.RemoteAddr()
}
