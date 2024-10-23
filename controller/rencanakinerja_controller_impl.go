package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RencanaKinerjaControllerImpl struct {
	rencanaKinerjaService service.RencanaKinerjaService
}

func NewRencanaKinerjaControllerImpl(rencanaKinerjaService service.RencanaKinerjaService) *RencanaKinerjaControllerImpl {
	return &RencanaKinerjaControllerImpl{
		rencanaKinerjaService: rencanaKinerjaService,
	}
}

func (controller *RencanaKinerjaControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaKinerjaCreateRequest := rencanakinerja.RencanaKinerjaCreateRequest{}
	helper.ReadFromRequestBody(request, &rencanaKinerjaCreateRequest)

	rencanaKinerjaResponse, err := controller.rencanaKinerjaService.Create(request.Context(), rencanaKinerjaCreateRequest)
	if err != nil {
		webResponse := web.WebRencanaKinerjaResponse{
			Code:   400,
			Status: "failed create rencana kinerja",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebRencanaKinerjaResponse{
		Code:   http.StatusCreated,
		Status: "success create rencana kinerja",
		Data:   rencanaKinerjaResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaKinerjaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaKinerjaUpdateRequest := rencanakinerja.RencanaKinerjaUpdateRequest{}
	helper.ReadFromRequestBody(request, &rencanaKinerjaUpdateRequest)

	rencanaKinerjaUpdateRequest.Id = params.ByName("id")

	rencanaKinerjaResponse, err := controller.rencanaKinerjaService.Update(request.Context(), rencanaKinerjaUpdateRequest)
	if err != nil {
		webResponse := web.WebRencanaKinerjaResponse{
			Code:   400,
			Status: "failed update rencana kinerja",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebRencanaKinerjaResponse{
		Code:   200,
		Status: "success update rencana kinerja",
		Data:   rencanaKinerjaResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaKinerjaControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiId := params.ByName("pegawai_id")

	query := request.URL.Query()
	kodeOpd := query.Get("kode_opd")
	tahun := query.Get("tahun")

	// Membuat map untuk menyimpan parameter opsional
	filterParams := make(map[string]string)

	if pegawaiId != "" {
		filterParams["pegawai_id"] = pegawaiId
	}
	if kodeOpd != "" {
		filterParams["kode_opd"] = kodeOpd
	}
	if tahun != "" {
		filterParams["tahun"] = tahun
	}

	rencanaKinerjaResponses, err := controller.rencanaKinerjaService.FindAll(request.Context(), pegawaiId, kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebRencanaKinerjaResponse{
			Code:   http.StatusBadRequest,
			Status: "failed get rencana kinerja",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebRencanaKinerjaResponse{
		Code:   http.StatusOK,
		Status: "success get rencana kinerja",
		Data:   rencanaKinerjaResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaKinerjaControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("rencana_kinerja_id")
	kodeOPD := request.URL.Query().Get("kode_opd")
	tahun := request.URL.Query().Get("tahun")

	result, err := controller.rencanaKinerjaService.FindById(request.Context(), id, kodeOPD, tahun)
	if err != nil {
		// Handle error
		webResponse := web.WebRencanaKinerjaResponse{
			Code:   404,
			Status: http.StatusText(http.StatusNotFound),
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim respons sukses
	webResponse := web.WebRencanaKinerjaResponse{
		Code:   http.StatusOK,
		Status: "success get rencana kinerja",
		Data:   result,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaKinerjaControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaKinerjaId := params.ByName("id")

	controller.rencanaKinerjaService.Delete(request.Context(), rencanaKinerjaId)
	webResponse := web.WebRencanaKinerjaResponse{
		Code:   200,
		Status: "success delete rencana kinerja",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaKinerjaControllerImpl) FindAllRencanaKinerja(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiId := params.ByName("pegawai_id")

	query := request.URL.Query()
	kodeOpd := query.Get("kode_opd")
	tahun := query.Get("tahun")

	// Membuat map untuk menyimpan parameter opsional
	filterParams := make(map[string]string)

	if pegawaiId != "" {
		filterParams["pegawai_id"] = pegawaiId
	}
	if kodeOpd != "" {
		filterParams["kode_opd"] = kodeOpd
	}
	if tahun != "" {
		filterParams["tahun"] = tahun
	}

	rencanaKinerjaResponses, err := controller.rencanaKinerjaService.FindAll(request.Context(), pegawaiId, kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebRencanaKinerjaResponse{
			Code:   http.StatusBadRequest,
			Status: "failed get rencana kinerja",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	actionButton := []web.ActionButton{
		{
			NameAction: "Create Rencana Kinerja",
			Method:     "POST",
			Url:        "/rencana_kinerja/create",
		},
	}
	webResponse := web.WebRencanaKinerjaResponse{
		Code:   http.StatusOK,
		Status: "success get rencana kinerja",
		Action: actionButton,
		Data:   rencanaKinerjaResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
