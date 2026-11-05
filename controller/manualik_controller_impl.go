package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ManualIKControllerImpl struct {
	ManualIKService service.ManualIKService
}

func NewManualIKControllerImpl(manualIKService service.ManualIKService) *ManualIKControllerImpl {
	return &ManualIKControllerImpl{
		ManualIKService: manualIKService,
	}
}

func (controller *ManualIKControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	indikatorId := params.ByName("indikatorId")

	manualIKCreateRequest := rencanakinerja.ManualIKCreateRequest{}
	helper.ReadFromRequestBody(request, &manualIKCreateRequest)

	manualIKResponse, err := controller.ManualIKService.Create(request.Context(), manualIKCreateRequest, indikatorId)
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
		Data:   manualIKResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ManualIKControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	indikatorId := params.ByName("indikatorId")

	manualIKUpdateRequest := rencanakinerja.ManualIKUpdateRequest{}
	helper.ReadFromRequestBody(request, &manualIKUpdateRequest)

	manualIKResponse, err := controller.ManualIKService.Update(request.Context(), manualIKUpdateRequest, indikatorId)
	helper.PanicIfError(err)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   manualIKResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ManualIKControllerImpl) FindManualIKByIndikatorId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	indikatorId := params.ByName("indikatorId")

	manualIKResponses, err := controller.ManualIKService.FindManualIKByIndikatorId(request.Context(), indikatorId)
	helper.PanicIfError(err)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   manualIKResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ManualIKControllerImpl) FindManualIKSasaranOpdByIndikatorId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	indikatorId := params.ByName("indikatorId")
	tahun := params.ByName("tahun")

	manualIKResponses, err := controller.ManualIKService.FindManualIKSasaranOpdByIndikatorId(request.Context(), indikatorId, tahun)
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
		Data:   manualIKResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// FindManualIKByRekin godoc
// @Summary      Manual IK by Rencana Kinerja
// @Description  Mengambil seluruh Manual IK yang terikat pada sebuah Rencana Kinerja berdasarkan ID rencana kinerja.
// @Tags         Manual IK
// @Tags         IntegrasiKabMadiun
// @Accept       json
// @Produce      json
// @Param        rencana_kinerja_id  path  string  true  "ID Rencana Kinerja"
// @Success      200  {object}  web.WebResponse{data=[]rencanakinerja.ManualIKResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     None
// @Router       /manual_ik/rekin/{rencana_kinerja_id} [get]
func (controller *ManualIKControllerImpl) FindManualIKByRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rekinId := params.ByName("rencana_kinerja_id")

	responses, err := controller.ManualIKService.FindManualIKByRekin(request.Context(), rekinId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   responses,
	})
}

// FindManualIKByPegawai godoc
// @Summary      Manual IK by Pegawai
// @Description  Mengambil daftar Manual IK berdasarkan NIP pegawai, kode OPD, dan tahun. Manual IK diambil dari indikator yang melekat pada rencana kinerja pegawai.
// @Tags         Manual IK
// @Tags         IntegrasiKabMadiun
// @Accept       json
// @Produce      json
// @Param        nip       path  string  true  "NIP Pegawai"
// @Param        kode_opd  path  string  true  "Kode OPD"  example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path  string  true  "Tahun"     example("2025")
// @Success      200  {object}  web.WebResponse{data=[]rencanakinerja.ManualIKResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     None
// @Router       /manual_ik/pegawai/{nip}/{kode_opd}/{tahun} [get]
func (controller *ManualIKControllerImpl) FindManualIKByPegawai(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	nip := params.ByName("nip")
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	manualIKResponses, err := controller.ManualIKService.FindManualIKByPegawai(request.Context(), nip, kodeOpd, tahun)
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
		Data:   manualIKResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
