package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/datamaster"
	"ekak_kabupaten_madiun/service"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type DataMasterControllerImpl struct {
	DataMasterService service.DataMasterService
}

func NewDataMasterControllerImpl(dataMasterService service.DataMasterService) *DataMasterControllerImpl {
	return &DataMasterControllerImpl{
		DataMasterService: dataMasterService,
	}
}

func (controller *DataMasterControllerImpl) DataRB(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tahunNextParams := r.URL.Query().Get("tahun_next")
	if tahunNextParams == "" {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "tahun_next params is missing",
		})
		return
	}
	tahunInt, err := strconv.Atoi(tahunNextParams)

	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "tahun_next params is malformatted",
		})
		return
	}

	response, err := controller.DataMasterService.DataRBByTahun(r.Context(), tahunInt)
	if err != nil {
		log.Printf("Error: %v", err)
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   500,
			Status: "ERROR",
			Data:   "Terjadi kesalahan server saat mengambil data RB.",
		})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   response,
	})
}

func (c *DataMasterControllerImpl) CreateRB(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.Header().Set("Content-Type", "application/json")

	// ===============================
	// 1. Ambil User dari JWT Claims
	// ===============================
	claims := helper.GetUserInfo(request.Context())
	if claims.UserId == 0 {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token invalid",
		})
		return
	}

	userId := claims.UserId

	// ===============================
	// 2. Decode JSON Request Body
	// ===============================
	rb := datamaster.RBRequest{}
	helper.ReadFromRequestBody(request, &rb)

	// ===============================
	// 3. Validasi Minimal
	// ===============================
	if rb.JenisRB == "" || rb.KegiatanUtama == "" || rb.TahunBaseline == 0 || rb.TahunNext == 0 {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "Field wajib tidak boleh kosong",
		})
		return
	}

	if len(rb.Indikator) == 0 {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "Indikator minimal 1",
		})
		return
	}

	// ===============================
	// 4. Panggil Service SaveRB
	// ===============================
	response, err := c.DataMasterService.SaveRB(request.Context(), rb, userId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "ERROR",
			Data:   err.Error(),
		})
		return
	}

	// ===============================
	// 5. Response sukses
	// ===============================
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   response,
	})
}

func (c *DataMasterControllerImpl) UpdateRB(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.Header().Set("Content-Type", "application/json")

	// ===============================
	// 1. Ambil User dari JWT Claims
	// ===============================
	claims := helper.GetUserInfo(request.Context())
	if claims.UserId == 0 {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token invalid",
		})
		return
	}

	userId := claims.UserId

	rbId := params.ByName("rb_id")

	if rbId == "" {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "ID TIDAK VALID",
		})
		return
	}

	rbIdNum, err := strconv.Atoi(rbId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "ID TIDAK VALID",
		})
		return
	}

	// ===============================
	// 2. Decode JSON Request Body
	// ===============================
	rb := datamaster.RBRequest{}
	helper.ReadFromRequestBody(request, &rb)

	// ===============================
	// 3. Validasi Minimal
	// ===============================
	if rb.JenisRB == "" || rb.KegiatanUtama == "" || rb.TahunBaseline == 0 || rb.TahunNext == 0 {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "Field wajib tidak boleh kosong",
		})
		return
	}

	if len(rb.Indikator) == 0 {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "Indikator minimal 1",
		})
		return
	}

	// ===============================
	// 4. Panggil Service SaveRB
	// ===============================
	response, err := c.DataMasterService.UpdateRB(request.Context(), rb, userId, rbIdNum)
	if err != nil {

		if err.Error() == "rb_not_found" {
			helper.WriteToResponseBody(writer, web.WebResponse{
				Code:   http.StatusNotFound,
				Status: "NOT_FOUND",
				Data:   "Data RB tidak ditemukan",
			})
			return
		}

		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "ERROR",
			Data:   err.Error(),
		})
		return
	}

	// ===============================
	// 5. Response sukses
	// ===============================
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   response,
	})
}

func (c *DataMasterControllerImpl) DeleteRB(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Guard, Hanya authorization tertinggi yang boleh
	// ===============================
	// 1. Ambil User dari JWT Claims
	// ===============================
	claims := helper.GetUserInfo(r.Context())
	if claims.UserId == 0 {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token invalid",
		})
		return
	}
	roles := claims.Roles
	allowed := slices.Contains(roles, "super_admin")

	if !allowed {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Data:   "Akses ditolak, role tidak diizinkan",
		})
		return
	}

	rbId := params.ByName("rb_id")
	if rbId == "" {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "ID TIDAK VALID",
		})
		return
	}

	rbIdNum, err := strconv.Atoi(rbId)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD_REQUEST",
			Data:   "ID TIDAK VALID",
		})
		return
	}

	err = c.DataMasterService.DeleteRB(r.Context(), rbIdNum)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "ERROR",
			Data:   err.Error(),
		})
		return
	}

	// ===============================
	// 5. Response sukses
	// ===============================
	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   http.StatusNoContent,
		Status: "SUCCESS",
		Data:   "RB Dihapus",
	})
}
