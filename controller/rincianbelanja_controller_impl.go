package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/rincianbelanja"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RincianBelanjaControllerImpl struct {
	rincianBelanjaService service.RincianBelanjaService
}

func NewRincianBelanjaControllerImpl(rincianBelanjaService service.RincianBelanjaService) *RincianBelanjaControllerImpl {
	return &RincianBelanjaControllerImpl{
		rincianBelanjaService: rincianBelanjaService,
	}
}

func (controller *RincianBelanjaControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rincianBelanjaCreateRequest := rincianbelanja.RincianBelanjaCreateRequest{}
	helper.ReadFromRequestBody(request, &rincianBelanjaCreateRequest)

	rincianBelanjaResponse, err := controller.rincianBelanjaService.Create(request.Context(), rincianBelanjaCreateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create rincian belanja",
		Data:   rincianBelanjaResponse,
	})
}

func (controller *RincianBelanjaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rincianBelanjaUpdateRequest := rincianbelanja.RincianBelanjaUpdateRequest{}
	helper.ReadFromRequestBody(request, &rincianBelanjaUpdateRequest)

	rincianBelanjaUpdateRequest.RenaksiId = params.ByName("renaksiId")
	rincianBelanjaResponse, err := controller.rincianBelanjaService.Update(request.Context(), rincianBelanjaUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Update Rincian Belanja",
		Data:   rincianBelanjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RincianBelanjaControllerImpl) FindRincianBelanjaAsn(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiId := params.ByName("pegawai_id")
	tahun := params.ByName("tahun")

	response := controller.rincianBelanjaService.FindRincianBelanjaAsn(request.Context(), pegawaiId, tahun)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success",
		Data:   response,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RincianBelanjaControllerImpl) LaporanRincianBelanjaOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	query := request.URL.Query()
	kodeOpd := query.Get("kode_opd")
	tahun := query.Get("tahun")

	filterParams := make(map[string]string)

	if kodeOpd != "" {
		filterParams["kode_opd"] = kodeOpd
	}
	if tahun != "" {
		filterParams["tahun"] = tahun
	}

	rincianBelanjaResponses, err := controller.rincianBelanjaService.LaporanRincianBelanjaOpd(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "failed get laporan rincian belanja",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success get laporan rincian belanja " + kodeOpd + " tahun " + tahun,
		Data:   rincianBelanjaResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
