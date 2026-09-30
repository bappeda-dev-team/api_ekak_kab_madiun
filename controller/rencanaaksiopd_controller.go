package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RencanaAksiOpdController interface {
	FindBySasaranOpdAndTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SyncJadwalPelaksanaan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllSasaranByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
