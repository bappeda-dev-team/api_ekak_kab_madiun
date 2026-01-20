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
	RincianKakByBulanTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// FindRekinSasaranOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateRekinLevel1(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateRekinLevel1(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindIdRekinLevel1(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindRekinLevel3(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindRekinAtasan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindBatchDetails(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}
