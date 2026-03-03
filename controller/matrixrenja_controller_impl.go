package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
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

func (controller *MatrixRenjaControllerImpl) GetRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetRenjaRanwal(request.Context(), kodeOpd, tahun)
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

func (controller *MatrixRenjaControllerImpl) GetRenjaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetRenjaRankhir(request.Context(), kodeOpd, tahun)
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
		Status: "Succes Get Matrix Renja Rankhir",
		Data:   matrixRenjaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *MatrixRenjaControllerImpl) CreateOrUpdateTarget(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	TargetRenjaRequest := programkegiatan.TargetRenjaRequest{}
	helper.ReadFromRequestBody(request, &TargetRenjaRequest)

	TargetRenjaResponse, err := controller.MatrixRenjaService.CreateOrUpdateTarget(request.Context(), TargetRenjaRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   TargetRenjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
