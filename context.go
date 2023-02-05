package irisx

import (
	"encoding/json"
	"strings"

	"github.com/daqiancode/jsoniter"
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Contextx struct {
	iris.Context
}

var JSON = jsoniter.Decapitalized

func (c *Contextx) ReadJSON(outPtr interface{}) error {
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
func (c *Contextx) ReadQuery(ptr interface{}) error {
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

func (c *Contextx) ReadForm(formObject interface{}) error {
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

func (c *Contextx) JSON(v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.ContentType(context.ContentJSONHeaderValue)
	_, err = c.Write(bs)
	return err
}

func (c *Contextx) Finish(data interface{}, err error) error {
	if err != nil {
		c.Error(err, 500)
		return err
	}
	return c.OK(data)
}

func (c *Contextx) Page(items interface{}, pageIndex, pageSize int, total int64, err error) error {
	if err != nil {
		c.Error(err, 500)
		return err
	}
	return c.OK(Page{PageIndex: pageIndex, PageSize: pageSize, Total: int(total), Items: items})
}

func (c *Contextx) OK(data interface{}) error {
	return c.JSON(Result{Data: data})
}

func (c *Contextx) Fail(message string, state, httpStatus int) error {
	c.StatusCode(httpStatus)
	return c.JSON(Result{State: state, ErrorCode: message})
}

// server internal error
func (c *Contextx) FailInternal(message string, state int) error {
	return c.Fail(message, state, 500)
}

// bussiness logic error
func (c *Contextx) FailService(message string, state int) error {
	return c.Fail(message, state, 406)
}

// request parameter error
func (c *Contextx) FailParams(fieldErrors map[string]string) error {
	c.StatusCode(400)
	return c.JSON(Result{State: 1, ErrorCode: "request parameter error", FieldErrors: fieldErrors})
}

func (c *Contextx) Error(err error, statusCode int) error {
	if err == nil {
		return c.OK(nil)
	}
	c.StatusCode(statusCode)
	return c.JSON(Result{State: 1, Error: err.Error()})
}

// request parameter error
func (c *Contextx) ErrorParam(err error) error {
	if es, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]string, len(es))
		for _, v := range es {
			fieldErrors[v.Field()] = v.ActualTag()
			// fieldErrors[v.Field()] = v.Error()
		}
		return c.FailParams(fieldErrors)
	}
	c.StatusCode(400)
	return c.JSON(Result{State: 1, Error: err.Error()})
}

func (c *Contextx) GetIP() string {
	ip := c.GetHeader("X-Real-Ip")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
		if ip != "" {
			return strings.TrimSpace(strings.Split(ip, ",")[0])
		}
	}
	return c.Context.RemoteAddr()
}
