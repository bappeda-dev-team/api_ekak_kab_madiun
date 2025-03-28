package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SasaranOpdController interface {
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdRencanaKinerja(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindIdPokinSasaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
