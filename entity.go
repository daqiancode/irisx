package irisx

type Result struct {
	State       int               `json:"state"`
	Data        interface{}       `json:"data,omitempty"`
	ErrorCode   string            `json:"errorCode,omitempty"`
	Error       string            `json:"error,omitempty"`
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
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
