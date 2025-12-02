package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type DataMasterController interface {
	DataRB(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	CreateRB(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	UpdateRB(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}
