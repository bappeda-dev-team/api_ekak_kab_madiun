package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"log"
	"net/http"
	"strconv"

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

func (controller *DataMasterControllerImpl) DataRB(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tahunNextParams := r.URL.Query().Get("tahun_next")
	if tahunNextParams == "" {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "tahun_next params is missing",
		})
		return
	}
	tahunInt, err := strconv.Atoi(tahunNextParams)

	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "tahun_next params is malformatted",
		})
		return
	}

	response, err := controller.DataMasterService.DataRBByTahun(r.Context(), tahunInt)
	if err != nil {
		log.Printf("Error: %v", err)
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   500,
			Status: "ERROR",
			Data:   "Terjadi kesalahan server saat mengambil data RB.",
		})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "success",
		Data:   response,
	})
}
