package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
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

func (controller *MatrixRenstraControllerImpl) CreateIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	IndikatorRenstraCreateRequest := programkegiatan.IndikatorRenstraCreateRequest{}
	helper.ReadFromRequestBody(request, &IndikatorRenstraCreateRequest)

	IndikatorRenstraResponse, err := controller.MatrixRenstraService.CreateIndikator(request.Context(), IndikatorRenstraCreateRequest)
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
		Code:   http.StatusOK,
		Status: "success create indikator renstra",
		Data:   IndikatorRenstraResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *MatrixRenstraControllerImpl) UpdateIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	IndikatorUpdateRequest := programkegiatan.UpdateIndikatorRequest{}
	helper.ReadFromRequestBody(request, &IndikatorUpdateRequest)

	IndikatorUpdateRequest.Id = params.ByName("id")

	IndikatorUpdateResponse, err := controller.MatrixRenstraService.UpdateIndikator(request.Context(), IndikatorUpdateRequest)
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
		Code:   http.StatusOK,
		Status: "success update indikator renstra",
		Data:   IndikatorUpdateResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *MatrixRenstraControllerImpl) DeleteIndikator(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	IndikatorId := params.ByName("id")

	err := controller.MatrixRenstraService.DeleteIndikator(request.Context(), IndikatorId)
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
		Code:   http.StatusOK,
		Status: "success delete indikator",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *MatrixRenstraControllerImpl) FindIndikatorById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	IndikatorId := params.ByName("id")
	indikatorResponse, err := controller.MatrixRenstraService.FindIndikatorById(request.Context(), IndikatorId)
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
		Code:   http.StatusOK,
		Status: "success delete indikator",
		Data:   indikatorResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
