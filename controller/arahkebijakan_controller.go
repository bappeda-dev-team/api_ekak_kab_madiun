package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ArahKebijakanController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}