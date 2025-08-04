package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IkuController interface {
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllIkuOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIkuActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
