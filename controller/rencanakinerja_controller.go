package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RencanaKinerjaController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllRencanaKinerja(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllRincianKak(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindRekinSasaranOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateSasaranOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateSasaranOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
