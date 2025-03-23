package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SubKegiatanTerpilihController interface {
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteSubKegiatanTerpilih(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllSubkegiatanByBidangUrusanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
