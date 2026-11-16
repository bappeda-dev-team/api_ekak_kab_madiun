package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ManualIKController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindManualIKByIndikatorId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindManualIKSasaranOpdByIndikatorId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindManualIKByPegawai(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindManualIKByRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
