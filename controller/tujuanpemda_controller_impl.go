package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/tujuanpemda"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TujuanPemdaControllerImpl struct {
	TujuanPemdaService service.TujuanPemdaService
}

func NewTujuanPemdaControllerImpl(tujuanPemdaService service.TujuanPemdaService) *TujuanPemdaControllerImpl {
	return &TujuanPemdaControllerImpl{
		TujuanPemdaService: tujuanPemdaService,
	}
}

func (controller *TujuanPemdaControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Decode request body
	tujuanPemdaCreateRequest := tujuanpemda.TujuanPemdaCreateRequest{}
	helper.ReadFromRequestBody(request, &tujuanPemdaCreateRequest)

	// Panggil service create
	tujuanPemdaResponse, err := controller.TujuanPemdaService.Create(request.Context(), tujuanPemdaCreateRequest)
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
		Code:   http.StatusCreated,
		Status: "success create tujuan pemda",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// Update godoc
// @Summary      Update Tujuan Pemda
// @Description  Memperbarui data tujuan pemda.
// @Tags         Tujuan Pemda
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Tujuan Pemda ID"
// @Param        body   body  tujuanpemda.TujuanPemdaUpdateRequest  true  "Data tujuan pemda yang akan diupdate"
// @Success      200  {object}  web.WebResponse{data=tujuanpemda.TujuanPemdaResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/update/{id} [put]
func (controller *TujuanPemdaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaUpdateRequest := tujuanpemda.TujuanPemdaUpdateRequest{}
	helper.ReadFromRequestBody(request, &tujuanPemdaUpdateRequest)

	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	tujuanPemdaUpdateRequest.Id = idInt

	// Panggil service update
	tujuanPemdaResponse, err := controller.TujuanPemdaService.Update(request.Context(), tujuanPemdaUpdateRequest)
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
		Status: "success update tujuan pemda",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaId := params.ByName("id")
	id, err := strconv.Atoi(tujuanPemdaId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.TujuanPemdaService.Delete(request.Context(), id)
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
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaId := params.ByName("id")

	id, err := strconv.Atoi(tujuanPemdaId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	tujuanPemdaResponse, err := controller.TujuanPemdaService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindAll(request.Context(), tahun, jenisPeriode)
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
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) UpdatePeriode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaUpdateRequest := tujuanpemda.TujuanPemdaUpdateRequest{}
	helper.ReadFromRequestBody(request, &tujuanPemdaUpdateRequest)

	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	tujuanPemdaUpdateRequest.Id = idInt

	// Panggil service update
	tujuanPemdaResponse, err := controller.TujuanPemdaService.UpdatePeriode(request.Context(), tujuanPemdaUpdateRequest)
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
		Status: "success update periode tujuan pemda",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) FindAllWithPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindAllWithPokin(request.Context(), tahunAwal, tahunAkhir, jenisPeriode)
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
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindPokinWithPeriode godoc
// @Summary Get Pokin With Periode
// @Description Mendapatkan data Pohon Kinerja berdasarkan ID Pokin dan jenis periode
// @Tags Tujuan Pemda
// @Accept json
// @Produce json
// @Param pokin_id path int true "Pokin ID"
// @Param jenis_periode path string true "Jenis Periode (triwulan, semester, tahunan, dll)"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Security     BearerAuth
// @Router /tujuan_pemda/pokin_with_periode/{pokin_id}/{jenis_periode} [get]
func (controller *TujuanPemdaControllerImpl) FindPokinWithPeriode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pokinId := params.ByName("pokin_id")
	jenisPeriode := params.ByName("jenis_periode")
	id, err := strconv.Atoi(pokinId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pokinWithPeriodeResponse, err := controller.TujuanPemdaService.FindPokinWithPeriode(request.Context(), id, jenisPeriode)
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
		Data:   pokinWithPeriodeResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindAllWithPokinRenstra godoc
// @Summary Get Pokin With Periode Renstra
// @Description Mendapatkan data Pohon Kinerja berdasarkan tahun awal, tahun akhir, dan jenis periode
// @Tags Tujuan Pemda
// @Accept json
// @Produce json
// @Param tahun_awal path string true "Tahun Awal"
// @Param tahun_akhir path string true "Tahun Akhir"
// @Param jenis_periode path string true "Jenis Periode (triwulan, semester, tahunan, dll)"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 500 {object} web.WebResponse
// @Security     BearerAuth
// @Router /tujuan_pemda/findall_with_pokin/{tahun_awal}/{tahun_akhir}/{jenis_periode} [get]
func (controller *TujuanPemdaControllerImpl) FindAllWithPokinRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")
	pokinWithPeriodeResponse, err := controller.TujuanPemdaService.FindAllWithPokinRenstra(request.Context(), tahunAwal, tahunAkhir, jenisPeriode)
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
		Data:   pokinWithPeriodeResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindTujuanPemdaRanwal godoc
// @Summary      Tujuan Pemda Ranwal
// @Description  Mendapatkan data tujuan pemda ranwal berdasarkan tahun dan jenis periode. Response langsung ke tujuan pemda (tanpa wrapper pohon kinerja). Target ranwal menimpa renstra jika tersedia.
// @Tags         Tujuan Pemda
// @Accept       json
// @Produce      json
// @Param        tahun         path  string  true  "Tahun yang berada dalam range RPJMD"  example("2025")
// @Param        jenis_periode path  string  true  "Jenis Periode"                        example("renstra")
// @Success      200  {object}  web.WebResponse{data=[]tujuanpemda.TujuanPemdaResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/ranwal/{tahun}/{jenis_periode} [get]
func (controller *TujuanPemdaControllerImpl) FindTujuanPemdaRanwal(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindTujuanPemdaRanwal(
		request.Context(), tahun, jenisPeriode,
	)
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
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindTujuanPemdaRankhir godoc
func (controller *TujuanPemdaControllerImpl) FindTujuanPemdaRankhir(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindTujuanPemdaRankhir(
		request.Context(), tahun, jenisPeriode,
	)
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
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) FindTujuanPemdaPenetapan(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindTujuanPemdaPenetapan(
		request.Context(), tahun, jenisPeriode,
	)
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
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// UpsertTargetPemdaLayer godoc
// @Summary      Upsert Target Tujuan Pemda Layer RKPD
// @Description  Menyimpan atau memperbarui target tujuan pemda untuk layer ranwal, rankhir, atau penetapan. Hanya mengubah target — metadata indikator (nama, rumus, dll.) tidak diubah.
// @Tags         Tujuan Pemda
// @Accept       json
// @Produce      json
// @Param        jenis  path  string                              true  "Jenis layer: ranwal, rankhir, atau penetapan"  example("ranwal")
// @Param        body   body  tujuanpemda.LayerTargetBatchRequest true  "Daftar target yang akan di-upsert"
// @Success      200  {object}  web.WebResponse{data=[]tujuanpemda.TargetResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/{jenis}/target/upsert [post]
func (controller *TujuanPemdaControllerImpl) UpsertTargetPemdaLayer(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	jenis := params.ByName("jenis")
	layerTargetRequest := tujuanpemda.LayerTargetBatchRequest{}
	helper.ReadFromRequestBody(request, &layerTargetRequest)
	targetResponses, err := controller.TujuanPemdaService.UpsertTargetPemdaLayer(
		request.Context(), jenis, layerTargetRequest,
	)
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
		Status: "success upsert target tujuan pemda",
		Data:   targetResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindTujuanPemdaRankhirDual godoc
// @Summary      Tujuan Pemda Rankhir (Dual Target)
// @Description  Menampilkan 2 target per indikator: ranwal dan rankhir. Tanpa fallback antar jenis.
// @Tags         Tujuan Pemda
// @Param        tahun         path  string  true  "Tahun"          example("2025")
// @Param        jenis_periode path  string  true  "Jenis Periode"  example("renstra")
// @Success      200  {object}  web.WebResponse{data=[]tujuanpemda.TujuanPemdaResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/rankhir/{tahun}/{jenis_periode} [get]
func (controller *TujuanPemdaControllerImpl) FindTujuanPemdaRankhirDual(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindTujuanPemdaRankhirDual(
		request.Context(), tahun, jenisPeriode,
	)
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
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindTujuanPemdaPenetapanDual godoc
// @Summary      Tujuan Pemda Penetapan (Dual Target)
// @Description  Menampilkan 2 target per indikator: rankhir dan penetapan. Tanpa fallback antar jenis.
// @Tags         Tujuan Pemda
// @Param        tahun         path  string  true  "Tahun"          example("2025")
// @Param        jenis_periode path  string  true  "Jenis Periode"  example("renstra")
// @Success      200  {object}  web.WebResponse{data=[]tujuanpemda.TujuanPemdaResponse}
// @Security     BearerAuth
// @Router       /tujuan_pemda/penetapan/{tahun}/{jenis_periode} [get]
func (controller *TujuanPemdaControllerImpl) FindTujuanPemdaPenetapanDual(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	result, err := controller.TujuanPemdaService.FindTujuanPemdaPenetapanDual(
		request.Context(), params.ByName("tahun"), params.ByName("jenis_periode"),
	)
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
		Data:   result,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// CreateTargetPemdaRankhir godoc
// @Summary      Create Target Tujuan Pemda Rankhir
// @Description  Membuat target baru untuk rankhir. Gagal jika target sudah ada.
// @Tags         Tujuan Pemda
// @Accept       json
// @Produce      json
// @Param        body   body  tujuanpemda.LayerTargetBatchRequest  true  "Daftar target yang akan dibuat"
// @Success      201  {object}  web.WebResponse{data=[]tujuanpemda.TargetResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/target/rankhir/create [post]
func (controller *TujuanPemdaControllerImpl) CreateTargetRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	jenis := "rankhir"
	var req tujuanpemda.LayerTargetBatchRequest
	helper.ReadFromRequestBody(request, &req)
	result, err := controller.TujuanPemdaService.CreateTargetPemdaLayer(
		request.Context(), jenis, req,
	)
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
		Code:   http.StatusCreated,
		Status: "success create target tujuan pemda",
		Data:   result,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// CreateTargetPemdaLayer godoc
// @Summary      Create Target Tujuan Pemda Penetapan
// @Description  Membuat target baru untuk penetapan. Gagal jika target sudah ada.
// @Tags         Tujuan Pemda
// @Accept       json
// @Produce      json
// @Success      201  {object}  web.WebResponse{data=[]tujuanpemda.TargetResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/target/penetapan/create [post]
func (controller *TujuanPemdaControllerImpl) CreateTargetPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	jenis := "penetapan"
	var req tujuanpemda.LayerTargetBatchRequest
	helper.ReadFromRequestBody(request, &req)
	result, err := controller.TujuanPemdaService.CreateTargetPemdaLayer(
		request.Context(), jenis, req,
	)
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
		Code:   http.StatusCreated,
		Status: "success create target tujuan pemda",
		Data:   result,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// UpdateTargetPemdaRankhir godoc
// @Summary      Update Target Tujuan Pemda Rankhir
// @Description  Memperbarui target dan satuan saja. Jenis, kode_indikator, dan tahun tidak diubah.
// @Tags         Tujuan Pemda
// @Accept       json
// @Produce      json
// @Param        body   body  tujuanpemda.LayerTargetUpdateBatchRequest  true  "Daftar target (wajib ada id)"
// @Success      200  {object}  web.WebResponse{data=[]tujuanpemda.TargetResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/target/rankhir/update [put]
func (controller *TujuanPemdaControllerImpl) UpdateTargetRankhir(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	jenis := "rankhir"
	var req tujuanpemda.LayerTargetUpdateBatchRequest
	helper.ReadFromRequestBody(request, &req)
	result, err := controller.TujuanPemdaService.UpdateTargetPemdaLayer(
		request.Context(), jenis, req,
	)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusInternalServerError, Status: "INTERNAL SERVER ERROR", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success update target tujuan pemda", Data: result,
	})
}

// UpdateTargetPemdaPenetapan godoc
// @Summary      Update Target Tujuan Pemda Penetapan
// @Description  Memperbarui target dan satuan saja. Jenis, kode_indikator, dan tahun tidak diubah.
// @Tags         Tujuan Pemda
// @Accept       json
// @Produce      json
// @Param        body   body  tujuanpemda.LayerTargetUpdateBatchRequest  true  "Daftar target (wajib ada id)"
// @Success      200  {object}  web.WebResponse{data=[]tujuanpemda.TargetResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/target/penetapan/update [put]
func (controller *TujuanPemdaControllerImpl) UpdateTargetPenetapan(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	jenis := "penetapan"
	var req tujuanpemda.LayerTargetUpdateBatchRequest
	helper.ReadFromRequestBody(request, &req)
	result, err := controller.TujuanPemdaService.UpdateTargetPemdaLayer(
		request.Context(), jenis, req,
	)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusInternalServerError, Status: "INTERNAL SERVER ERROR", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success update target tujuan pemda", Data: result,
	})
}

// lock
// LockTujuanPemda godoc
// @Summary      Lock Data Tujuan Pemda
// @Description  Mengunci data tujuan pemda untuk tahun tertentu. Setelah lock: indikator tidak bisa ditambah/diubah, delete diblokir, target ranwal & penetapan tidak bisa diubah. Target renstra & rankhir masih boleh diubah.
// @Tags         Tujuan Pemda Lock
// @Accept       json
// @Produce      json
// @Param        tahun  path  string  true  "Tahun yang akan di-lock"  example("2025")
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/lock/{tahun} [post]
func (controller *TujuanPemdaControllerImpl) LockTujuanPemda(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	tahun := params.ByName("tahun")
	err := controller.TujuanPemdaService.LockTujuanPemda(request.Context(), tahun)
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
		Status: "success lock tujuan pemda",
		Data: tujuanpemda.LockDataPemdaResponse{
			Tahun:  tahun,
			Locked: true,
		},
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// UnlockTujuanPemda godoc
// @Summary      Unlock Data Tujuan Pemda
// @Description  Membuka kunci data tujuan pemda untuk tahun tertentu.
// @Tags         Tujuan Pemda Lock
// @Accept       json
// @Produce      json
// @Param        tahun  path  string  true  "Tahun yang akan di-unlock"  example("2025")
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/lock/{tahun} [delete]
func (controller *TujuanPemdaControllerImpl) UnlockTujuanPemda(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	tahun := params.ByName("tahun")
	err := controller.TujuanPemdaService.UnlockTujuanPemda(request.Context(), tahun)
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
		Status: "success unlock tujuan pemda",
		Data: tujuanpemda.LockDataPemdaResponse{
			Tahun:  tahun,
			Locked: false,
		},
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// IsTujuanPemdaLocked godoc
// @Summary      Cek Status Lock Tujuan Pemda
// @Description  Mengecek apakah data tujuan pemda untuk tahun tertentu sedang terkunci.
// @Tags         Tujuan Pemda Lock
// @Accept       json
// @Produce      json
// @Param        tahun  path  string  true  "Tahun yang dicek"  example("2025")
// @Success      200  {object}  web.WebResponse{data=tujuanpemda.LockDataPemdaResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/lock/{tahun} [get]
func (controller *TujuanPemdaControllerImpl) IsTujuanPemdaLocked(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	tahun := params.ByName("tahun")
	locked, err := controller.TujuanPemdaService.IsTujuanPemdaLocked(request.Context(), tahun)
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
		Data: tujuanpemda.LockDataPemdaResponse{
			Tahun:  tahun,
			Locked: locked,
		},
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// FindAllLockTujuanPemda godoc
// @Summary      Daftar Semua Lock Tujuan Pemda
// @Description  Mengambil seluruh daftar tahun yang sedang di-lock untuk modul tujuan pemda.
// @Tags         Tujuan Pemda Lock
// @Accept       json
// @Produce      json
// @Success      200  {object}  web.WebResponse{data=[]tujuanpemda.LockDataPemdaResponse}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /tujuan_pemda/lock [get]
func (controller *TujuanPemdaControllerImpl) FindAllLockTujuanPemda(
	writer http.ResponseWriter, request *http.Request, params httprouter.Params,
) {
	result, err := controller.TujuanPemdaService.FindAllLockTujuanPemda(request.Context())
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
		Data:   result,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
