package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CascadingOpdController interface {
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
