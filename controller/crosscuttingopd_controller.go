package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CrosscuttingOpdController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ApproveOrReject(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteUnused(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinByCrosscuttingStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindOPDCrosscuttingFrom(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteCrosscuttingDiterima(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
