package controller

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type CSFController interface {
	FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
