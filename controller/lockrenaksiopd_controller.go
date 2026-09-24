package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LockRenaksiOpdController interface {
	Lock(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Unlock(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
