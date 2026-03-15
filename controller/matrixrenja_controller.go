package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenjaController interface {
	GetRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetRenjaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetRenjaPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertBatchIndikatorRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertBatchIndikatorRenjaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertBatchIndikatorRenjaPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpsertAnggaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
