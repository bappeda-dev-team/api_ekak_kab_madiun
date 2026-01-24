package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type BidangUrusanController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByKodeOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	CreateOPD(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteOPD(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindBidangUrusanTerpilihByKodeOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
