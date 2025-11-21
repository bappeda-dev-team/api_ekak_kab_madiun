package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CascadingOpdController interface {
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByRekinPegawaiAndId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByNip(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByMultipleRekinPegawai(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	MultiRekinDetailByOpdAndTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
