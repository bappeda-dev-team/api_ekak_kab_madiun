package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TargetProgramKegiatanControllerImpl struct {
	svc service.TargetProgramKegiatanService
}

func NewTargetProgramKegiatanControllerImpl(svc service.TargetProgramKegiatanService) *TargetProgramKegiatanControllerImpl {
	return &TargetProgramKegiatanControllerImpl{svc: svc}
}

// ─────────────────────────────────────────────────────────────────────────────
// helpers internal
// ─────────────────────────────────────────────────────────────────────────────

func (c *TargetProgramKegiatanControllerImpl) handleCreate(w http.ResponseWriter, r *http.Request, jenis string) {
	var req programkegiatan.TargetProgramKegiatanCreateRequest
	helper.ReadFromRequestBody(r, &req)

	resp, err := c.svc.Create(r.Context(), jenis, req)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: resp})
}

func (c *TargetProgramKegiatanControllerImpl) handleUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: "id tidak valid"})
		return
	}

	var req programkegiatan.TargetProgramKegiatanUpdateRequest
	helper.ReadFromRequestBody(r, &req)

	resp, err := c.svc.Update(r.Context(), id, req)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: resp})
}

// ─────────────────────────────────────────────────────────────────────────────
// CREATE per jenis
// ─────────────────────────────────────────────────────────────────────────────

// CreateRenstra godoc
// @Summary      Buat Target Renstra
// @Description  Membuat target program kegiatan untuk fase Renstra (jenis='renstra').
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        request  body  programkegiatan.TargetProgramKegiatanCreateRequest  true  "Data target"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/renstra/create [post]
func (c *TargetProgramKegiatanControllerImpl) CreateRenstra(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.handleCreate(w, r, "renstra")
}

// CreateRenja godoc
// @Summary      Buat Target Renja
// @Description  Membuat target program kegiatan untuk fase Renja (jenis='renja').
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        request  body  programkegiatan.TargetProgramKegiatanCreateRequest  true  "Data target"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/renja/create [post]
func (c *TargetProgramKegiatanControllerImpl) CreateRenja(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.handleCreate(w, r, "renja")
}

// CreateRankhir godoc
// @Summary      Buat Target Rancangan Akhir
// @Description  Membuat target program kegiatan untuk fase Rancangan Akhir (jenis='rankhir').
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        request  body  programkegiatan.TargetProgramKegiatanCreateRequest  true  "Data target"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/rankhir/create [post]
func (c *TargetProgramKegiatanControllerImpl) CreateRankhir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.handleCreate(w, r, "rankhir")
}

// CreatePenetapan godoc
// @Summary      Buat Target Penetapan
// @Description  Membuat target program kegiatan untuk fase Penetapan (jenis='penetapan').
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        request  body  programkegiatan.TargetProgramKegiatanCreateRequest  true  "Data target"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/penetapan/create [post]
func (c *TargetProgramKegiatanControllerImpl) CreatePenetapan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c.handleCreate(w, r, "penetapan")
}

// ─────────────────────────────────────────────────────────────────────────────
// UPDATE per jenis (hanya target & satuan)
// ─────────────────────────────────────────────────────────────────────────────

// UpdateRenstra godoc
// @Summary      Update Target Renstra
// @Description  Mengubah target dan satuan untuk data dengan jenis='renstra'. Kolom jenis dan tahun tidak dapat diubah.
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        id       path  int                                                   true  "ID target"
// @Param        request  body  programkegiatan.TargetProgramKegiatanUpdateRequest  true  "target & satuan baru"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/renstra/update/{id} [put]
func (c *TargetProgramKegiatanControllerImpl) UpdateRenstra(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.handleUpdate(w, r, p)
}

// UpdateRenja godoc
// @Summary      Update Target Renja
// @Description  Mengubah target dan satuan untuk data dengan jenis='renja'. Kolom jenis dan tahun tidak dapat diubah.
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        id       path  int                                                   true  "ID target"
// @Param        request  body  programkegiatan.TargetProgramKegiatanUpdateRequest  true  "target & satuan baru"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/renja/update/{id} [put]
func (c *TargetProgramKegiatanControllerImpl) UpdateRenja(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.handleUpdate(w, r, p)
}

// UpdateRankhir godoc
// @Summary      Update Target Rancangan Akhir
// @Description  Mengubah target dan satuan untuk data dengan jenis='rankhir'. Kolom jenis dan tahun tidak dapat diubah.
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        id       path  int                                                   true  "ID target"
// @Param        request  body  programkegiatan.TargetProgramKegiatanUpdateRequest  true  "target & satuan baru"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/rankhir/update/{id} [put]
func (c *TargetProgramKegiatanControllerImpl) UpdateRankhir(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.handleUpdate(w, r, p)
}

// UpdatePenetapan godoc
// @Summary      Update Target Penetapan
// @Description  Mengubah target dan satuan untuk data dengan jenis='penetapan'. Kolom jenis dan tahun tidak dapat diubah.
// @Tags         Target Program Kegiatan
// @Accept       json
// @Produce      json
// @Param        id       path  int                                                   true  "ID target"
// @Param        request  body  programkegiatan.TargetProgramKegiatanUpdateRequest  true  "target & satuan baru"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/penetapan/update/{id} [put]
func (c *TargetProgramKegiatanControllerImpl) UpdatePenetapan(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.handleUpdate(w, r, p)
}

// ─────────────────────────────────────────────────────────────────────────────
// SHARED — Delete, FindById, FindAll
// ─────────────────────────────────────────────────────────────────────────────

// Delete godoc
// @Summary      Hapus Target Program Kegiatan
// @Description  Menghapus satu target berdasarkan id.
// @Tags         Target Program Kegiatan
// @Produce      json
// @Param        id  path  int  true  "ID target"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/delete/{id} [delete]
func (c *TargetProgramKegiatanControllerImpl) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: "id tidak valid"})
		return
	}

	if err := c.svc.Delete(r.Context(), id); err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: "berhasil dihapus"})
}

// FindById godoc
// @Summary      Detail Target Program Kegiatan
// @Description  Mendapatkan satu target berdasarkan id.
// @Tags         Target Program Kegiatan
// @Produce      json
// @Param        id  path  int  true  "ID target"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/detail/{id} [get]
func (c *TargetProgramKegiatanControllerImpl) FindById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: "id tidak valid"})
		return
	}

	resp, err := c.svc.FindById(r.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: resp})
}

// FindByKodeIndikatorAndKodeOpd godoc
// @Summary      Daftar Target berdasarkan Kode Indikator + Kode OPD
// @Description  Mendapatkan semua target (semua jenis) untuk kombinasi kode_indikator dan kode_opd.
// @Tags         Target Program Kegiatan
// @Produce      json
// @Param        kode_indikator  path  string  true  "Kode indikator"
// @Param        kode_opd        path  string  true  "Kode OPD"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.TargetProgramKegiatanResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /target_program_kegiatan/findall/{kode_indikator}/{kode_opd} [get]
func (c *TargetProgramKegiatanControllerImpl) FindByKodeIndikatorAndKodeOpd(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	kodeIndikator := p.ByName("kode_indikator")
	kodeOpd := p.ByName("kode_opd")

	resp, err := c.svc.FindByKodeIndikatorAndKodeOpd(r.Context(), kodeIndikator, kodeOpd)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusOK, Status: "OK", Data: resp})
}
