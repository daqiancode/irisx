package irisx

import "github.com/go-playground/validator/v10"

type Result struct {
	State          int               `json:"state"`
	Data           interface{}       `json:"data,omitempty"`
	ErrorCode      string            `json:"errorCode,omitempty"`
	Error          string            `json:"error,omitempty"`
	FieldErrors    map[string]string `json:"fieldErrors,omitempty"`
	HttpStatusCode int               `json:"-"`
}

type Page struct {
	PageIndex  int         `json:"pageIndex"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
	Total      int         `json:"total"`
	Items      interface{} `json:"items"`
}

func NewPage(items interface{}, pageIndex, pageSize, total int) Page {
	totalPages := 0
	if pageSize != 0 {
		totalPages = total / pageSize
		if totalPages*pageSize != total {
			totalPages += 1
		}
	}
	return Page{
		PageIndex:  pageIndex,
		PageSize:   pageIndex,
		TotalPages: totalPages,
		Total:      total,
		Items:      items,
	}
}

type ErrorCodeGetter interface {
	GetErrorCode() string
}
type FieldErrorsGetter interface {
	GetFieldErrors() map[string]string
}
type HttpStatusCodeGetter interface {
	GetHttpStatusCode() int
}
type StateGetter interface {
	GetState() int
}

// https://stackoverflow.com/questions/3050518/what-http-status-response-code-should-i-use-if-the-request-is-missing-a-required
var ParameterErrorHttpStatusCode = 422

// ParseError error to Result & http status code
func ParseError(err error) Result {
	if err == nil {
		return Result{}
	}
	r := Result{State: 1, Error: err.Error()}
	if v, ok := err.(FieldErrorsGetter); ok {
		r.HttpStatusCode = ParameterErrorHttpStatusCode
		r.FieldErrors = v.GetFieldErrors()
	}
	if v, ok := err.(StateGetter); ok {
		r.State = v.GetState()
	}
	if v, ok := err.(ErrorCodeGetter); ok {
		r.ErrorCode = v.GetErrorCode()
	}
	if v, ok := err.(HttpStatusCodeGetter); ok {
		r.HttpStatusCode = v.GetHttpStatusCode()
	}
	return r
}

type ValidationErrors struct {
	Err         string
	FieldErrors map[string]string
}

func (s ValidationErrors) Error() string {
	return s.Err
}

func (s ValidationErrors) GetFieldErrors() map[string]string {
	return s.FieldErrors
}

func ParseValidationErrors(err error) error {
	if es, ok := err.(*validator.ValidationErrors); ok {
		fieldErrors := make(map[string]string, len(*es))
		for _, v := range *es {
			fieldErrors[v.Field()] = v.ActualTag()
			// fieldErrors[v.Field()] = v.Error()
		}
		return &ValidationErrors{FieldErrors: fieldErrors, Err: "validation error"}
	}
	return err
}
