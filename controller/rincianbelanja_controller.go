package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RincianBelanjaController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindRincianBelanjaAsn(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	LaporanRincianBelanjaOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	LaporanRincianBelanjaPegawai(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
