package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SasaranPemdaController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllWithPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSasaranPemdaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	// Dual target
	FindSasaranPemdaRankhirDual(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSasaranPemdaPenetapanDual(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	// Target layer
	CreateTargetRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTargetRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateTargetPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTargetPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	LockSasaranPemda(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UnlockSasaranPemda(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	IsSasaranPemdaLocked(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllLockSasaranPemda(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
