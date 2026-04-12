package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProgramPrioritasPusatController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByKodeProgramPrioritasPusat(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindUnusedByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByIdTerkait(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}