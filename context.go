package irisx

import (
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
var IgnoreErrPath = true
var DecapitalizeErrorField = true
var DefaultFieldErrorMsg = "field format error"

func (s Context) HandleValidationErrors(err error) error {
	if err == nil || IgnoreErrPath && iris.IsErrPath(err) {
		return nil
	}
	var field string
	if es, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]string, len(es))
		for _, v := range es {
			field = v.Field()
			if DecapitalizeErrorField {
				field = decapitalize(field)
			}
			fieldErrors[field] = s.Tr("error."+v.ActualTag(), v.Param())
			if fieldErrors[field] == "" {
				fieldErrors[field] = DefaultFieldErrorMsg
			}
		}
		return &ValidationErrors{FieldErrors: fieldErrors, Err: "request parameter error"}
	}
	return err
}

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
	err = c.Application().Validate(outPtr)
	return c.HandleValidationErrors(err)
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

	return c.HandleValidationErrors(c.Application().Validate(ptr))
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
	return c.HandleValidationErrors(c.Application().Validate(formObject))
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
	c.StatusCode(422)
	return c.JSON(Result{State: 1, ErrorInfo: "request parameter error", FieldErrors: fieldErrors})
}

// type stackTracer interface {
// 	StackTrace() errors.StackTrace
// }

func (c *Context) Error(err error, statusCode int) error {
	if err == nil {
		return c.OK(nil)
	}
	r := ParseError(err)
	if r.HttpStatusCode != 0 {
		c.StatusCode(r.HttpStatusCode)
	} else {
		c.StatusCode(statusCode)
	}
	return c.JSON(r)
}

// request parameter error
func (c *Context) ErrorParam(err error) error {
	return c.Error(c.HandleValidationErrors(err), 422)
}

func (c *Context) GetIP() string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	ip = c.GetHeader("X-Real-Ip")
	if ip != "" {
		return ip
	}
	return c.Context.RemoteAddr()
}

func (c *Context) TranslateError(err error) Result {
	var r Result
	if err == nil {
		return r
	}
	r.State = 1
	var msgKey string
	var msgParams []any
	if e, ok := err.(ErrorKeyGetter); ok {
		msgKey = e.GetErrorKey()
	}
	if e, ok := err.(ErrorParamsGetter); ok {
		msgParams = e.GetErrorParams()
	}
	msg := c.Tr(msgKey, msgParams...)
	r.ErrorInfo = msg
	e := c.HandleValidationErrors(err)
	if ve, ok := e.(*ValidationErrors); ok {
		r.FieldErrors = ve.FieldErrors
	}
	return r
}
