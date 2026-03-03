package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenjaController interface {
	GetRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetRenjaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateOrUpdateTarget(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
