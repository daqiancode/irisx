package irisx

import (
	"strings"

	"github.com/daqiancode/gocommons/commons"
	"github.com/daqiancode/gocommons/commons/states"
	"github.com/daqiancode/gocommons/logger"
	"github.com/daqiancode/jsoniter"
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
)

type Contextx struct {
	iris.Context
}

var IrisxLog = logger.NewLogger(map[string]string{"src": "irisx"}, true, false)

var (
	json = jsoniter.Decapitalized
)

func GetJSONSerializer() jsoniter.API {
	return json
}

func (c *Contextx) ReadJSON(outPtr interface{}) error {
	body, restoreBody, err := context.GetBody(c.Request(), true)
	if err != nil {
		IrisxLog.Error().Err(err).Msg("Contextx.ReadJSON failed")
		return err
	}
	restoreBody()
	err = json.Unmarshal(body, outPtr)
	if err != nil {
		IrisxLog.Error().Err(err).Msg("Contextx.ReadJSON - json.Unmarshal failed")
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
		IrisxLog.Error().Err(err).Msg("Contextx.ReadQuery - DecodeQuery failed")
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
		IrisxLog.Error().Err(err).Msg("Contextx.ReadForm - DecodeForm failed")
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

func (c *Contextx) OK(data interface{}) error {
	r := OK(data)
	return c.JSON(r)
}

func (c *Contextx) Fail(message string, state, httpStatus int) error {
	c.StatusCode(httpStatus)
	return c.JSON(commons.Result{State: state, Message: message})
}

func (c *Contextx) Error(err error) error {
	switch v := err.(type) {
	case *commons.ServiceError:
		IrisxLog.Error().Err(err).Msg("Contextx.Error - Service error:" + v.Error())
		return c.ErrorService(v)
	default:
		IrisxLog.Error().Err(err).Msg("Contextx.Error - Internal error:" + err.Error())
		c.StatusCode(500)
	}
	return c.JSON(HandleError(err))
}

func (c *Contextx) ErrorService(err error) error {
	c.StatusCode(406)
	return c.JSON(HandleError(err))
}

func (c *Contextx) ErrorParam(err error) error {
	if es, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]string, len(es))
		for _, v := range es {
			fieldErrors[v.Field()] = v.ActualTag()
			// fieldErrors[v.Field()] = v.Error()
		}
		return c.ErrorFields(fieldErrors)
	}
	c.StatusCode(400)
	return c.JSON(commons.Result{State: states.InvalidParam, Message: err.Error()})
}
func (c *Contextx) ErrorFields(fieldErrors map[string]string) error {
	c.StatusCode(400)
	return c.JSON(commons.Result{State: states.InvalidParam, Message: "request parameter error", FieldErrors: fieldErrors})
}

func (c *Contextx) GetUID() string {
	claims := jwt.Get(c.Context).(*RbacClaims)
	if claims == nil {
		return ""
	}
	return claims.Subject
}

func (c *Contextx) IsLogined() bool {
	t := jwt.Get(c.Context)
	if t == nil {
		return false
	}

	if claims, ok := t.(*RbacClaims); ok {
		return len(claims.Subject) > 0
	}
	return false
}
func (c *Contextx) GetIP() string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	return c.Context.RemoteAddr()
}
