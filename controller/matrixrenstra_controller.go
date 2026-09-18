package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenstraController interface {
	GetByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetByKodeSubKegiatanVersiKedua(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateIndikatorV2(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIndikatorRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertTarget(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertAnggaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertBatchIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
