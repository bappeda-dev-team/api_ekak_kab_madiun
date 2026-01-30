package controller

import (
	"net/http"

	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pkopd"
	"ekak_kabupaten_madiun/service"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PkControllerImpl struct {
	pkOpdService service.PkService
}

func NewPkControllerImpl(pkOpdService service.PkService) *PkControllerImpl {
	return &PkControllerImpl{
		pkOpdService: pkOpdService,
	}
}

func (controller *PkControllerImpl) FindAllPkOpdTahunan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahunStr := params.ByName("tahun")

	if kodeOpd == "" {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Kode OPD harus terisi",
			Data:   nil,
		}
		helper.WriteToResponseBody(w, webResponse)
		return
	}

	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Tahun tidak sesuai",
			Data:   nil,
		}
		helper.WriteToResponseBody(w, webResponse)
		return
	}

	response, err := controller.pkOpdService.FindByKodeOpdTahun(r.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "[ERROR] PK OPD sedang tidak dapat diakses",
			Data:   nil,
		}
		helper.WriteToResponseBody(w, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   response,
	}

	helper.WriteToResponseBody(w, webResponse)
}

func (controller *PkControllerImpl) HubungkanRekin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	hubungkanRekinRequest := pkopd.PkOpdRequest{}
	helper.ReadFromRequestBody(r, &hubungkanRekinRequest)

	hubungkanResponse, err := controller.pkOpdService.HubungkanRekin(r.Context(), hubungkanRekinRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "[ERROR] Terjadi kesalahan di hubungkan rekin",
			Data:   nil,
		}
		helper.WriteToResponseBody(w, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   hubungkanResponse,
	}

	helper.WriteToResponseBody(w, webResponse)
}

func (controller *PkControllerImpl) HubungkanAtasan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	hubungkanAtasanRequest := pkopd.HubungkanAtasanRequest{}
	helper.ReadFromRequestBody(r, &hubungkanAtasanRequest)

	hubungkanResponse, err := controller.pkOpdService.HubungkanAtasan(r.Context(), hubungkanAtasanRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "[ERROR] Terjadi kesalahan di hubungkan atasan",
			Data:   nil,
		}
		helper.WriteToResponseBody(w, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   hubungkanResponse,
	}

	helper.WriteToResponseBody(w, webResponse)
}
