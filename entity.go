package irisx

type Result struct {
	State       int               `json:"state"`
	Data        interface{}       `json:"data,omitempty"`
	ErrorCode   string            `json:"errorCode,omitempty"`
	Error       string            `json:"error,omitempty"`
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
}

type Page struct {
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
	Total     int         `json:"total"`
	Items     interface{} `json:"items"`
}
