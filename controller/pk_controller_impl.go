package controller

import (
	"net/http"
	"slices"
	"strings"

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
	query := r.URL.Query()
	mode := query.Get("mode")

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

	claims := helper.GetUserInfo(r.Context())
	roleIndividu := []string{
		"level_1",
		"level_2",
		"level_3",
		"level_4",
	}

	roleAdmin := []string{
		"super_admin",
		"admin_opd",
		"admin_kecamatan",
		"reviewer",
	}

	isIndividu := slices.ContainsFunc(claims.Roles, func(role string) bool {
		return slices.Contains(roleIndividu, role)
	})

	isAdmin := slices.ContainsFunc(claims.Roles, func(role string) bool {
		return slices.Contains(roleAdmin, role)
	})
	idPegawai := claims.Nip

	var response pkopd.PkOpdResponse
	var filterMode pkopd.PkFilterMode

	switch strings.ToUpper(strings.TrimSpace(mode)) {
	case "SELF":
		filterMode = pkopd.PkFilterSelf

	case "DIRECT_SUPERVISOR":
		filterMode = pkopd.PkFilterDirectSupervisor

	default:
		filterMode = pkopd.PkFilterAllCandidates
	}

	switch {
	case isAdmin:
		response, err = controller.pkOpdService.FindByKodeOpdTahun(
			r.Context(), kodeOpd, tahun,
		)

	case isIndividu:
		response, err = controller.pkOpdService.FindByPegawaiTahun(
			r.Context(), idPegawai, kodeOpd, tahun, filterMode,
		)

	default:
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   http.StatusForbidden,
			Status: "[ERROR] Tidak memiliki hak akses",
			Data:   nil,
		})
		return
	}

	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "[ERROR] PK OPD sedang tidak dapat diakses",
			Data:   nil,
		})
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

func (controller *PkControllerImpl) KunciPK(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	kunciPkRequest := pkopd.KunciPkRequest{}
	helper.ReadFromRequestBody(r, &kunciPkRequest)

	resp, err := controller.pkOpdService.KunciPK(r.Context(), kunciPkRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "[ERROR] Terjadi kesalahan saat mengunci PK",
			Data:   nil,
		}
		helper.WriteToResponseBodyWstatus(w, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   resp,
	}

	helper.WriteToResponseBodyWstatus(w, webResponse)
}

func (controller *PkControllerImpl) BukaKunciPK(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	kunciPkRequest := pkopd.KunciPkRequest{}
	helper.ReadFromRequestBody(r, &kunciPkRequest)

	resp, err := controller.pkOpdService.BukaKunciPK(r.Context(), kunciPkRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "[ERROR] Terjadi kesalahan saat membuka kunci PK",
			Data:   nil,
		}
		helper.WriteToResponseBodyWstatus(w, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   resp,
	}

	helper.WriteToResponseBodyWstatus(w, webResponse)
}

func (controller *PkControllerImpl) PkPenetapan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := r.URL.Query()
	idPegawai := req.Get("id_pegawai")
	kodeOpd := req.Get("kode_opd")
	tahunStr := req.Get("tahun")

	if idPegawai == "" {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Data:   nil,
			Error:  "id_pegawai harus terisi",
		}
		helper.WriteToResponseBodyWstatus(w, webResponse)
		return

	}

	if kodeOpd == "" {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Data:   nil,
			Error:  "kode_opd harus terisi",
		}
		helper.WriteToResponseBodyWstatus(w, webResponse)
		return
	}

	tahun, err := strconv.Atoi(tahunStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: http.StatusText(http.StatusBadRequest),
			Data:   nil,
			Error:  "tahun tidak valid",
		}
		helper.WriteToResponseBodyWstatus(w, webResponse)
		return
	}

	response, err := controller.pkOpdService.FindPkPenetapan(r.Context(), idPegawai, kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: http.StatusText(http.StatusInternalServerError),
			Data:   nil,
			Error:  "PK Penetapan TIDAK DAPAT DIAKSES",
		}
		helper.WriteToResponseBodyWstatus(w, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   response,
	}

	helper.WriteToResponseBodyWstatus(w, webResponse)
}
