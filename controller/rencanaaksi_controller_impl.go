package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/rencanaaksi"
	"ekak_kabupaten_madiun/service"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type RencanaAksiControllerImpl struct {
	RencanaAksiService service.RencanaAksiService
}

func NewRencanaAksiControllerImpl(rencanaAksiService service.RencanaAksiService) *RencanaAksiControllerImpl {
	return &RencanaAksiControllerImpl{
		RencanaAksiService: rencanaAksiService,
	}
}

func (controller *RencanaAksiControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaKinerjaId := params.ByName("rekin_id")

	rencanaAksiCreateRequest := rencanaaksi.RencanaAksiCreateRequest{}
	helper.ReadFromRequestBody(request, &rencanaAksiCreateRequest)

	rencanaAksiCreateRequest.RencanaKinerjaId = rencanaKinerjaId

	rencanaAksiResponse, err := controller.RencanaAksiService.Create(request.Context(), rencanaAksiCreateRequest)
	if err != nil {
		webResponse := web.WebRencanaAksiResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Create Rencana Aksi",
		Data:   rencanaAksiResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaAksiUpdateRequest := rencanaaksi.RencanaAksiUpdateRequest{}
	helper.ReadFromRequestBody(request, &rencanaAksiUpdateRequest)

	rencanaAksiId := params.ByName("rencanaaksiId")
	rencanaAksiUpdateRequest.Id = rencanaAksiId

	rencanaAksiResponse, err := controller.RencanaAksiService.Update(request.Context(), rencanaAksiUpdateRequest)
	if err != nil {
		webResponse := web.WebRencanaAksiResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Update Rencana Aksi",
		Data:   rencanaAksiResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// FindAllByRencanaKinerjaId godoc
// @Summary      Rencana Aksi by Rencana Kinerja ID
// @Description  Mendapatkan daftar rencana aksi berdasarkan ID rencana kinerja.
// @Tags         Rencana Aksi
// @Tags         IntegrasiKabMadiun
// @Accept       json
// @Produce      json
// @Param        rencana_kinerja_id  path  string  true  "ID Rencana Kinerja"
// @Success      200  {object}  web.WebRencanaAksiResponse{renaksi=[]rencanaaksi.RencanaAksiResponse}
// @Failure      400  {object}  web.WebRencanaAksiResponse
// @Security     None
// @Router       /rencana_aksi/findall/{rencana_kinerja_id} [get]
func (controller *RencanaAksiControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaKinerjaId := params.ByName("rencana_kinerja_id")

	rencanaAksiResponses, err := controller.RencanaAksiService.FindAll(request.Context(), rencanaKinerjaId)
	if err != nil {
		webResponse := web.WebRencanaAksiResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Get Rencana Aksi",
		Data:   rencanaAksiResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// FindById godoc
// @Summary      Rencana Aksi by ID
// @Description  Mendapatkan detail rencana aksi berdasarkan ID.
// @Tags         Rencana Aksi
// @Tags         IntegrasiKabMadiun
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "ID Rencana Aksi"
// @Success      200  {object}  web.WebRencanaAksiResponse{renaksi=[]rencanaaksi.RencanaAksiResponse}
// @Failure      400  {object}  web.WebRencanaAksiResponse
// @Security     None
// @Router       /detail-rencana_aksi/{id} [get]
func (controller *RencanaAksiControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaAksiId := params.ByName("rencanaaksiId")

	rencanaAksiResponse, err := controller.RencanaAksiService.FindById(request.Context(), rencanaAksiId)
	if err != nil {
		webResponse := web.WebRencanaAksiResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Get Rencana Aksi",
		Data:   rencanaAksiResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaAksiId := params.ByName("rencanaaksiId")

	err := controller.RencanaAksiService.Delete(request.Context(), rencanaAksiId)
	if err != nil {
		webResponse := web.WebRencanaAksiResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Delete Rencana Aksi",
		Data:   "Rencana Aksi berhasil dihapus",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *RencanaAksiControllerImpl) FindAllByRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rencanaKinerjaId := params.ByName("rencana_kinerja_id")

	rencanaAksiResponses, err := controller.RencanaAksiService.FindAll(request.Context(), rencanaKinerjaId)
	if err != nil {
		webResponse := web.WebRencanaAksiResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Create Rencana Aksi",
			Method:     "POST",
			Url:        fmt.Sprintf("%s:%s/rencana_aksi/create/rencanaaksi/:rekin_id", host, port),
		},
	}

	webResponse := web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Get Rencana Aksi",
		Action: buttonActions,
		Data:   rencanaAksiResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// FindByRekin godoc
// @Summary      Rencana Aksi by Rencana Kinerja
// @Description  Mengambil seluruh rencana aksi beserta bobot pelaksanaan berdasarkan ID rencana kinerja.
// @Tags         Rencana Aksi
// @Tags         IntegrasiKabMadiun
// @Accept       json
// @Produce      json
// @Param        rencana_kinerja_id  path  string  true  "ID Rencana Kinerja"
// @Success      200  {object}  web.WebRencanaAksiResponse{renaksi=[]rencanaaksi.RencanaAksiResponse}
// @Failure      400  {object}  web.WebRencanaAksiResponse
// @Security     None
// @Router       /rencana_aksi/rekin/{rencana_kinerja_id} [get]
func (controller *RencanaAksiControllerImpl) FindByRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rekinId := params.ByName("rencana_kinerja_id")

	responses, err := controller.RencanaAksiService.FindByRekin(request.Context(), rekinId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebRencanaAksiResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Get Rencana Aksi by Rencana Kinerja",
		Data:   responses,
	})
}

// FindByPegawai godoc
// @Summary      Rencana Aksi by Pegawai
// @Description  Mengambil seluruh rencana aksi beserta bobot pelaksanaan berdasarkan NIP, kode OPD, dan tahun. Data difilter dari rencana kinerja pegawai (pegawai_id = NIP).
// @Tags         Rencana Aksi
// @Tags         IntegrasiKabMadiun
// @Accept       json
// @Produce      json
// @Param        nip       path  string  true  "NIP Pegawai"
// @Param        kode_opd  path  string  true  "Kode OPD"  example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path  string  true  "Tahun"     example("2025")
// @Success      200  {object}  web.WebRencanaAksiResponse{renaksi=[]rencanaaksi.RencanaAksiPegawaiResponse}
// @Failure      400  {object}  web.WebRencanaAksiResponse
// @Security     None
// @Router       /rencana_aksi/pegawai/{nip}/{kode_opd}/{tahun} [get]
func (controller *RencanaAksiControllerImpl) FindByPegawai(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	nip := params.ByName("nip")
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	rencanaAksiResponses, err := controller.RencanaAksiService.FindByPegawai(request.Context(), nip, kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebRencanaAksiResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebRencanaAksiResponse{
		Code:   http.StatusOK,
		Status: "Success Get Rencana Aksi Pegawai",
		Data:   rencanaAksiResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
