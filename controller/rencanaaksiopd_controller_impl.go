package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/renaksiopd"
	"ekak_kabupaten_madiun/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type RencanaAksiOpdControllerImpl struct {
	RencanaAksiOpdService service.RencanaAksiOpdService
}

func NewRencanaAksiOpdControllerImpl(rencanaAksiOpdService service.RencanaAksiOpdService) *RencanaAksiOpdControllerImpl {
	return &RencanaAksiOpdControllerImpl{
		RencanaAksiOpdService: rencanaAksiOpdService,
	}
}

func (controller *RencanaAksiOpdControllerImpl) FindBySasaranOpdAndTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	sasaranOpdId := params.ByName("sasaran_opd_id")
	sasaranOpdIdInt, err := strconv.Atoi(sasaranOpdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	tahun := params.ByName("tahun")
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.FindBySasaranOpdAndTahun(request.Context(), sasaranOpdIdInt, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiOpdControllerImpl) SyncJadwalPelaksanaan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rekinId := params.ByName("rekin_id")
	err := controller.RencanaAksiOpdService.SyncJadwalPelaksanaan(request.Context(), rekinId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Jadwal pelaksanaan berhasil disinkronkan",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiOpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaAksiOpdCreateRequest := renaksiopd.RencanaAksiOpdCreateRequest{}
	err := json.NewDecoder(request.Body).Decode(&rencanaAksiOpdCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.Create(request.Context(), rencanaAksiOpdCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiOpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaAksiOpdUpdateRequest := renaksiopd.RencanaAksiOpdUpdateRequest{}
	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdUpdateRequest.Id = idInt

	err = json.NewDecoder(request.Body).Decode(&rencanaAksiOpdUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.Update(request.Context(), rencanaAksiOpdUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiOpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	err = controller.RencanaAksiOpdService.Delete(request.Context(), idInt)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Rencana aksi opd berhasil dihapus",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiOpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	rencanaAksiOpdResponse, err := controller.RencanaAksiOpdService.FindById(request.Context(), idInt)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   rencanaAksiOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiOpdControllerImpl) FindAllSasaranByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	sasaranList, err := controller.RencanaAksiOpdService.FindAllSasaranByTahun(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   sasaranList,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
