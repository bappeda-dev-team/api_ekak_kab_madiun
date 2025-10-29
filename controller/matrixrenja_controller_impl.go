package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenjaControllerImpl struct {
	MatrixRenjaService service.MatrixRenjaService
}

func NewMatrixRenjaControllerImpl(matrixRenjaService service.MatrixRenjaService) *MatrixRenjaControllerImpl {
	return &MatrixRenjaControllerImpl{MatrixRenjaService: matrixRenjaService}
}

func (controller *MatrixRenjaControllerImpl) GetByKodeOpdAndTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetByKodeOpdAndTahun(request.Context(), kodeOpd, tahun)
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
		Data:   matrixRenjaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
