package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SasaranOpdController interface {
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindIdPokinSasaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSasaranRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSasaranRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSasaranRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateIndikatorRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateIndikatorRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIndikatorRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIndikatorRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteIndikatorTargetRenja(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSasaranPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateIndikatorPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIndikatorPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
