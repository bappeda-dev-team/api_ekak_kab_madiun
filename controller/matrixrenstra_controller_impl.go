package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenstraControllerImpl struct {
	MatrixRenstraService service.MatrixRenstraService
}

func NewMatrixRenstraControllerImpl(matrixRenstraService service.MatrixRenstraService) *MatrixRenstraControllerImpl {
	return &MatrixRenstraControllerImpl{MatrixRenstraService: matrixRenstraService}
}

func (controller *MatrixRenstraControllerImpl) GetByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahunAwal := request.URL.Query().Get("tahun_awal")
	tahunAkhir := request.URL.Query().Get("tahun_akhir")

	matrixRenstraResponses, err := controller.MatrixRenstraService.GetByKodeSubKegiatan(request.Context(), kodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   matrixRenstraResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
