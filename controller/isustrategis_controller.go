package controller

import (
	"github.com/julienschmidt/httprouter"

	"net/http"
)

type CSFController interface {
	AllCsfsByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
