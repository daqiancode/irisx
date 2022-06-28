package irisx

import (
	"github.com/daqiancode/gocommons/commons"
	"github.com/daqiancode/gocommons/commons/states"
)

func OK(data interface{}) commons.Result {
	return commons.Result{
		State: 0,
		Data:  data,
	}
}

func HandleError(err error) commons.Result {
	if err == nil {
		return commons.Result{State: 0}
	}
	if e, ok := err.(*commons.ServiceError); ok {
		return commons.Result{State: states.ServiceError, Message: e.Message}
	}
	return commons.Result{State: states.InternalError, Message: err.Error()}
}
