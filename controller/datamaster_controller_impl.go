package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type DataMasterControllerImpl struct {
	DataMasterService service.DataMasterService
}

func NewDataMasterControllerImpl(dataMasterService service.DataMasterService) *DataMasterControllerImpl {
	return &DataMasterControllerImpl{
		DataMasterService: dataMasterService,
	}
}

func (controler *DataMasterControllerImpl) DataRB(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Hello, World",
	})
}
