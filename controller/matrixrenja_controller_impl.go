package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MatrixRenjaControllerImpl struct {
	MatrixRenjaService service.MatrixRenjaService
}

func NewMatrixRenjaControllerImpl(matrixRenjaService service.MatrixRenjaService) *MatrixRenjaControllerImpl {
	return &MatrixRenjaControllerImpl{MatrixRenjaService: matrixRenjaService}
}

// GetRenjaRanwal godoc
// @Summary      Matrix Renja Rancangan Awal
// @Description  Mendapatkan data matrix renja rancangan awal berdasarkan kode OPD dan tahun. Anggaran diambil dari tb_pagu jenis='ranwal'.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.UrusanDetailResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/ranwal/{kode_opd}/{tahun} [get]
func (controller *MatrixRenjaControllerImpl) GetRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisIndikator := "ranwal"
	jenisPagu := "renstra"
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetRenja(request.Context(), kodeOpd, tahun, jenisIndikator, jenisPagu)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   matrixRenjaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Matrix Renja Rankhir
// @Description  Mendapatkan data matrix renja rancangan rankhir berdasarkan kode OPD dan tahun. Anggaran diambil dari tb_pagu jenis='rankhir'.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.UrusanDetailResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/rankhir/{kode_opd}/{tahun} [get]
func (controller *MatrixRenjaControllerImpl) GetRenjaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetRenjaRankhir(request.Context(), kodeOpd, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Succes Get Matrix Renja Rankhir",
		Data:   matrixRenjaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Matrix Renja Penetapan
// @Description  Mendapatkan data matrix renja rancangan penetapan berdasarkan kode OPD dan tahun. Anggaran diambil dari tb_pagu jenis='penetapan'.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.UrusanDetailResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/penetapan/{kode_opd}/{tahun} [get]
func (controller *MatrixRenjaControllerImpl) GetRenjaPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisIndikator := "penetapan"
	jenisPagu := "penetapan"
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetRenja(request.Context(), kodeOpd, tahun, jenisIndikator, jenisPagu)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Succes Get Matrix Renja Penetapan",
		Data:   matrixRenjaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Upsert Batch Indikator Renja Rancangan Awal
// @Description  Upsert batch indikator renja rancangan awal.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body   []programkegiatan.IndikatorRenjaCreateRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.IndikatorUpsertResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/indikator/ranwal/upsert [post]
func (controller *MatrixRenjaControllerImpl) UpsertBatchIndikatorRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	BatchIndikatorRenjaRequest := []programkegiatan.IndikatorRenjaCreateRequest{}
	helper.ReadFromRequestBody(request, &BatchIndikatorRenjaRequest)
	for i := range BatchIndikatorRenjaRequest {
		BatchIndikatorRenjaRequest[i].Jenis = "ranwal"
	}
	BatchIndikatorRenjaResponse, err := controller.MatrixRenjaService.UpsertBatchIndikatorRenja(request.Context(), BatchIndikatorRenjaRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   BatchIndikatorRenjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Upsert Batch Indikator Renja Rankhir
// @Description  Upsert batch indikator renja rankhir.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body    []programkegiatan.IndikatorRenjaCreateRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.IndikatorUpsertResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/indikator/rankhir/upsert [post]
func (controller *MatrixRenjaControllerImpl) UpsertBatchIndikatorRenjaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	BatchIndikatorRenjaRequest := []programkegiatan.IndikatorRenjaCreateRequest{}
	helper.ReadFromRequestBody(request, &BatchIndikatorRenjaRequest)

	for i := range BatchIndikatorRenjaRequest {
		BatchIndikatorRenjaRequest[i].Jenis = "rankhir"
	}

	BatchIndikatorRenjaResponse, err := controller.MatrixRenjaService.UpsertBatchIndikatorRenja(request.Context(), BatchIndikatorRenjaRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   BatchIndikatorRenjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Upsert Batch Indikator Renja Penetapan
// @Description  Upsert batch indikator renja penetapan.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body   []programkegiatan.IndikatorRenjaCreateRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.IndikatorUpsertResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/indikator/penetapan/upsert [post]
func (controller *MatrixRenjaControllerImpl) UpsertBatchIndikatorRenjaPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	BatchIndikatorRenjaRequest := []programkegiatan.IndikatorRenjaCreateRequest{}
	helper.ReadFromRequestBody(request, &BatchIndikatorRenjaRequest)

	for i := range BatchIndikatorRenjaRequest {
		BatchIndikatorRenjaRequest[i].Jenis = "penetapan"
	}

	BatchIndikatorRenjaResponse, err := controller.MatrixRenjaService.UpsertBatchIndikatorRenja(request.Context(), BatchIndikatorRenjaRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   BatchIndikatorRenjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Upsert Anggaran Renja
// @Description  Upsert anggaran renja.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body     programkegiatan.AnggaranRenjaRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=programkegiatan.AnggaranRenjaResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/anggaran_penetapan/upsert [post]
func (controller *MatrixRenjaControllerImpl) UpsertAnggaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	AnggaranRenjaRequest := programkegiatan.AnggaranRenjaRequest{}
	helper.ReadFromRequestBody(request, &AnggaranRenjaRequest)

	AnggaranRenjaResponse, err := controller.MatrixRenjaService.UpsertAnggaran(request.Context(), AnggaranRenjaRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   AnggaranRenjaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
