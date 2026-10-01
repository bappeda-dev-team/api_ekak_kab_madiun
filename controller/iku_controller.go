package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IkuController interface {
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllIkuOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIkuActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIkuOpdActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllIkuRenjaOpdRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllIkuRenjaOpdRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllIkuRenjaOpdPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	// v2 — IKU pemda filter tematik.tahun
	FindIkuPemdaRanwalV2(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindIkuPemdaRankhirDualV2(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindIkuPemdaPenetapanDualV2(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
