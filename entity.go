package irisx

type Result struct {
	State          int               `json:"state"`
	Data           interface{}       `json:"data,omitempty"`
	ErrorCode      string            `json:"errorCode,omitempty"`
	ErrorInfo      string            `json:"error,omitempty"`
	FieldErrors    map[string]string `json:"fieldErrors,omitempty"`
	HttpStatusCode int               `json:"-"`
	ErrorKey       string            `json:"-"` // error message key
	ErrorParams    []any             `json:"-"` // error message params
}

func (s Result) Error() string {
	return s.ErrorInfo
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

type ErrorKeyGetter interface {
	GetErrorKey() string
}
type ErrorParamsGetter interface {
	GetErrorParams() []any
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
	r := Result{State: 1, ErrorInfo: err.Error()}
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

func decapitalize(s string) string {
	bs := []byte(s)
	i := 0
	for ; i < len(s); i++ {
		if !isUpper(s[i]) {
			break
		}
	}
	if i != len(s) {
		i--
	}
	for j := 0; j < i; j++ {
		bs[j] += 32
	}

	return string(bs)
}
func isUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}
