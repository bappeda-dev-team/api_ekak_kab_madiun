package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/sasaranpemda"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SasaranPemdaControllerImpl struct {
	sasaranPemdaService service.SasaranPemdaService
}

func NewSasaranPemdaControllerImpl(sasaranPemdaService service.SasaranPemdaService) *SasaranPemdaControllerImpl {
	return &SasaranPemdaControllerImpl{sasaranPemdaService: sasaranPemdaService}
}

// ── Helper private ───────────────────────────────────────────────
func sasaranPemdaErr(err error) web.WebResponse {
	code := http.StatusInternalServerError
	status := "INTERNAL SERVER ERROR"
	if helper.IsTargetValidationError(err) {
		code = http.StatusBadRequest
		status = "BAD REQUEST"
	}
	return web.WebResponse{Code: code, Status: status, Data: err.Error()}
}

// ── CRUD ─────────────────────────────────────────────────────────
func (c *SasaranPemdaControllerImpl) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req sasaranpemda.SasaranPemdaCreateRequest
	if err := helper.DecodeJSONBody(r, &req); err != nil {
		helper.WriteBadRequest(w, err)
		return
	}
	result, err := c.sasaranPemdaService.Create(r.Context(), req)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusCreated, "success create sasaran pemda", result)
}
func (c *SasaranPemdaControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var req sasaranpemda.SasaranPemdaUpdateRequest
	if err := helper.DecodeJSONBody(r, &req); err != nil {
		helper.WriteBadRequest(w, err)
		return
	}
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, "BAD REQUEST", "Invalid ID format")
		return
	}
	req.Id = id
	result, err := c.sasaranPemdaService.Update(r.Context(), req)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "success update sasaran pemda", result)
}
func (c *SasaranPemdaControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, "BAD REQUEST", "Invalid ID format")
		return
	}
	if err := c.sasaranPemdaService.Delete(r.Context(), id); err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "success delete sasaran pemda", nil)
}
func (c *SasaranPemdaControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, "BAD REQUEST", "Invalid ID format")
		return
	}
	result, err := c.sasaranPemdaService.FindById(r.Context(), id)
	if err != nil {
		helper.WriteJSON(w, http.StatusNotFound, "NOT FOUND", err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, "OK", result)
}
func (c *SasaranPemdaControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := c.sasaranPemdaService.FindAll(r.Context(), params.ByName("tahun"))
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "OK", result)
}
func (c *SasaranPemdaControllerImpl) FindAllWithPokin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := c.sasaranPemdaService.FindAllWithPokin(
		r.Context(),
		params.ByName("tahun_awal"),
		params.ByName("tahun_akhir"),
		params.ByName("jenis_periode"),
	)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "OK", result)
}

// ── Dual ─────────────────────────────────────────────────────────
func (c *SasaranPemdaControllerImpl) FindSasaranPemdaRankhirDual(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := c.sasaranPemdaService.FindSasaranPemdaRankhirDual(
		r.Context(), params.ByName("tahun"), params.ByName("jenis_periode"),
	)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "OK", result)
}
func (c *SasaranPemdaControllerImpl) FindSasaranPemdaPenetapanDual(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := c.sasaranPemdaService.FindSasaranPemdaPenetapanDual(
		r.Context(), params.ByName("tahun"), params.ByName("jenis_periode"),
	)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "OK", result)
}

// ── Target Layer ─────────────────────────────────────────────────
func (c *SasaranPemdaControllerImpl) CreateTargetRankhir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req sasaranpemda.LayerTargetBatchRequest
	if err := helper.DecodeJSONBody(r, &req); err != nil {
		helper.WriteBadRequest(w, err)
		return
	}
	result, err := c.sasaranPemdaService.CreateTargetSasaranLayer(r.Context(), "rankhir", req)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusCreated, "success create target sasaran pemda", result)
}
func (c *SasaranPemdaControllerImpl) CreateTargetPenetapan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req sasaranpemda.LayerTargetBatchRequest
	if err := helper.DecodeJSONBody(r, &req); err != nil {
		helper.WriteBadRequest(w, err)
		return
	}
	result, err := c.sasaranPemdaService.CreateTargetSasaranLayer(r.Context(), "penetapan", req)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusCreated, "success create target sasaran pemda", result)
}
func (c *SasaranPemdaControllerImpl) UpdateTargetRankhir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req sasaranpemda.LayerTargetUpdateBatchRequest
	if err := helper.DecodeJSONBody(r, &req); err != nil {
		helper.WriteBadRequest(w, err)
		return
	}
	result, err := c.sasaranPemdaService.UpdateTargetSasaranLayer(r.Context(), "rankhir", req)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "success update target sasaran pemda", result)
}
func (c *SasaranPemdaControllerImpl) UpdateTargetPenetapan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req sasaranpemda.LayerTargetUpdateBatchRequest
	if err := helper.DecodeJSONBody(r, &req); err != nil {
		helper.WriteBadRequest(w, err)
		return
	}
	result, err := c.sasaranPemdaService.UpdateTargetSasaranLayer(r.Context(), "penetapan", req)
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "success update target sasaran pemda", result)
}

// LockSasaranPemda godoc
// @Summary      Lock Data Sasaran Pemda
// @Description  Mengunci data sasaran pemda untuk tahun tertentu (jenis: sasaran_pemda). Setelah lock: create/update/delete sasaran diblokir, target penetapan tidak bisa diubah. Target rankhir masih boleh diubah.
// @Tags         Sasaran Pemda Lock
// @Accept       json
// @Produce      json
// @Param        tahun  path  string  true  "Tahun yang akan di-lock"  example("2025")
// @Success      200  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_pemda/lock/{tahun} [post]
func (c *SasaranPemdaControllerImpl) LockSasaranPemda(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := c.sasaranPemdaService.LockSasaranPemda(r.Context(), params.ByName("tahun"))
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "success lock sasaran pemda", result)
}

// UnlockSasaranPemda godoc
// @Summary      Unlock Data Sasaran Pemda
// @Description  Membuka kunci data sasaran pemda untuk tahun tertentu.
// @Tags         Sasaran Pemda Lock
// @Accept       json
// @Produce      json
// @Param        tahun  path  string  true  "Tahun yang akan di-unlock"  example("2025")
// @Success      200  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_pemda/lock/{tahun} [delete]
func (c *SasaranPemdaControllerImpl) UnlockSasaranPemda(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := c.sasaranPemdaService.UnlockSasaranPemda(r.Context(), params.ByName("tahun"))
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "success unlock sasaran pemda", result)
}

// IsSasaranPemdaLocked godoc
// @Summary      Cek Status Lock Sasaran Pemda
// @Description  Mengecek apakah data sasaran pemda untuk tahun tertentu sedang terkunci.
// @Tags         Sasaran Pemda Lock
// @Accept       json
// @Produce      json
// @Param        tahun  path  string  true  "Tahun yang dicek"  example("2025")
// @Success      200  {object}  web.WebResponse{data=sasaranpemda.LockDataPemdaResponse}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_pemda/lock/{tahun} [get]
func (c *SasaranPemdaControllerImpl) IsSasaranPemdaLocked(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := c.sasaranPemdaService.IsSasaranPemdaLocked(r.Context(), params.ByName("tahun"))
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "OK", result)
}

// FindAllLockSasaranPemda godoc
// @Summary      Daftar Semua Lock Sasaran Pemda
// @Description  Mengambil seluruh daftar tahun yang sedang di-lock untuk modul sasaran pemda (jenis: sasaran_pemda).
// @Tags         Sasaran Pemda Lock
// @Accept       json
// @Produce      json
// @Success      200  {object}  web.WebResponse{data=[]sasaranpemda.LockDataPemdaResponse}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_pemda/lock [get]
func (c *SasaranPemdaControllerImpl) FindAllLockSasaranPemda(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result, err := c.sasaranPemdaService.FindAllLockSasaranPemda(r.Context())
	if err != nil {
		helper.WriteToResponseBodyWstatus(w, sasaranPemdaErr(err))
		return
	}
	helper.WriteJSON(w, http.StatusOK, "OK", result)
}
