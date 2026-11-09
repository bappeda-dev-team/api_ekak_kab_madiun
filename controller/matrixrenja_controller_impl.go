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
// @Success      200  {object}  web.WebResponse{data=[]renja.RanwalMatrixResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/ranwal/{kode_opd}/{tahun} [get]
func (controller *MatrixRenjaControllerImpl) GetRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisPagu := "renstra"
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetRenja(request.Context(), kodeOpd, tahun, jenisPagu)
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
// @Success      200  {object}  web.WebResponse{data=[]renja.RankhirMatrixResponse}
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
// @Success      200  {object}  web.WebResponse{data=[]renja.PenetapanMatrixResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/penetapan/{kode_opd}/{tahun} [get]
func (controller *MatrixRenjaControllerImpl) GetRenjaPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisPagu := "penetapan"
	matrixRenjaResponses, err := controller.MatrixRenjaService.GetRenjaPenetapan(request.Context(), kodeOpd, tahun, jenisPagu)
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
// @Description endpoint sudah tidak digunakan
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body   []programkegiatan.IndikatorRenjaCreateRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.IndikatorUpsertResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/indikator/ranwal/upsert [post]
func (controller *MatrixRenjaControllerImpl) UpsertBatchIndikatorRenjaRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// BatchIndikatorRenjaRequest := []programkegiatan.IndikatorRenjaCreateRequest{}
	// helper.ReadFromRequestBody(request, &BatchIndikatorRenjaRequest)
	// for i := range BatchIndikatorRenjaRequest {
	// 	BatchIndikatorRenjaRequest[i].Jenis = "ranwal"
	// }
	// BatchIndikatorRenjaResponse, err := controller.MatrixRenjaService.UpsertBatchIndikatorRenja(request.Context(), BatchIndikatorRenjaRequest)
	// if err != nil {
	// 	webResponse := web.WebResponse{
	// 		Code:   http.StatusBadRequest,
	// 		Status: "BAD REQUEST",
	// 		Data:   err.Error(),
	// 	}
	// 	helper.WriteToResponseBody(writer, webResponse)
	// 	return
	// }

	// webResponse := web.WebResponse{
	// 	Code:   200,
	// 	Status: "OK",
	// 	Data:   BatchIndikatorRenjaResponse,
	// }
	// helper.WriteToResponseBody(writer, webResponse)
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusBadRequest, // atau http.StatusMethodNotAllowed (405)
		Status: "BAD REQUEST",
		Data:   "endpoint sudah tidak digunakan",
	})
}

// @Summary      Upsert Batch Indikator Renja Rankhir
// @Description  endpoint sudah tidak digunakan.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body    []programkegiatan.IndikatorRenjaCreateRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.IndikatorUpsertResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/indikator/rankhir/upsert [post]
func (controller *MatrixRenjaControllerImpl) UpsertBatchIndikatorRenjaRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// BatchIndikatorRenjaRequest := []programkegiatan.IndikatorRenjaCreateRequest{}
	// helper.ReadFromRequestBody(request, &BatchIndikatorRenjaRequest)

	// for i := range BatchIndikatorRenjaRequest {
	// 	BatchIndikatorRenjaRequest[i].Jenis = "rankhir"
	// }

	// BatchIndikatorRenjaResponse, err := controller.MatrixRenjaService.UpsertBatchIndikatorRenja(request.Context(), BatchIndikatorRenjaRequest)
	// if err != nil {
	// 	webResponse := web.WebResponse{
	// 		Code:   http.StatusBadRequest,
	// 		Status: "BAD REQUEST",
	// 		Data:   err.Error(),
	// 	}
	// 	helper.WriteToResponseBody(writer, webResponse)
	// 	return
	// }

	// webResponse := web.WebResponse{
	// 	Code:   200,
	// 	Status: "OK",
	// 	Data:   BatchIndikatorRenjaResponse,
	// }
	// helper.WriteToResponseBody(writer, webResponse)
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusBadRequest, // atau http.StatusMethodNotAllowed (405)
		Status: "BAD REQUEST",
		Data:   "endpoint sudah tidak digunakan",
	})
}

// @Summary      Upsert Batch Indikator Renja Penetapan
// @Description  endpoint sudah tidak digunakan.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body   []programkegiatan.IndikatorRenjaCreateRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=[]programkegiatan.IndikatorUpsertResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/indikator/penetapan/upsert [post]
func (controller *MatrixRenjaControllerImpl) UpsertBatchIndikatorRenjaPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// BatchIndikatorRenjaRequest := []programkegiatan.IndikatorRenjaCreateRequest{}
	// helper.ReadFromRequestBody(request, &BatchIndikatorRenjaRequest)

	// for i := range BatchIndikatorRenjaRequest {
	// 	BatchIndikatorRenjaRequest[i].Jenis = "penetapan"
	// }

	// BatchIndikatorRenjaResponse, err := controller.MatrixRenjaService.UpsertBatchIndikatorRenjaPenetapan(request.Context(), BatchIndikatorRenjaRequest)
	// if err != nil {
	// 	webResponse := web.WebResponse{
	// 		Code:   http.StatusBadRequest,
	// 		Status: "BAD REQUEST",
	// 		Data:   err.Error(),
	// 	}
	// 	helper.WriteToResponseBody(writer, webResponse)
	// 	return
	// }

	// webResponse := web.WebResponse{
	// 	Code:   200,
	// 	Status: "OK",
	// 	Data:   BatchIndikatorRenjaResponse,
	// }
	// helper.WriteToResponseBody(writer, webResponse)
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusBadRequest, // atau http.StatusMethodNotAllowed (405)
		Status: "BAD REQUEST",
		Data:   "endpoint sudah tidak digunakan",
	})
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

// UpsertTargetRankhir godoc
// @Summary      Upsert Target Rankhir
// @Description  Simpan atau update target rancangan akhir untuk satu indikator.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body     programkegiatan.TargetRenjaLayerRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=renja.TargetLayerResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/target/rankhir [put]
func (c *MatrixRenjaControllerImpl) UpsertTargetRankhir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req programkegiatan.TargetRenjaLayerRequest
	helper.ReadFromRequestBody(r, &req)
	data, err := c.MatrixRenjaService.UpsertTargetRenjaLayer(r.Context(), "rankhir", req)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: 200, Status: "OK", Data: data})
}

// UpsertTargetPenetapan godoc
// @Summary      Upsert Target Penetapan
// @Description  Simpan atau update target penetapan untuk satu indikator.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body     programkegiatan.TargetRenjaLayerRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse{data=renja.TargetLayerResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/target/penetapan [put]
func (c *MatrixRenjaControllerImpl) UpsertTargetPenetapan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req programkegiatan.TargetRenjaLayerRequest
	helper.ReadFromRequestBody(r, &req)
	data, err := c.MatrixRenjaService.UpsertTargetRenjaLayer(r.Context(), "penetapan", req)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: 200, Status: "OK", Data: data})
}

// DeleteTargetRankhir godoc
// @Summary      Delete Target Rankhir
// @Description  Hapus target rankhir berdasarkan kode indikator dan tahun.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        kode_indikator  path  string  true  "Kode Indikator"
// @Param        tahun           path  string  true  "Tahun"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/target/rankhir/{kode_indikator}/{tahun} [delete]
func (c *MatrixRenjaControllerImpl) DeleteTargetRankhir(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	kodeIndikator := p.ByName("kode_indikator")
	tahun := p.ByName("tahun")
	if err := c.MatrixRenjaService.DeleteTargetRenjaLayer(r.Context(), kodeIndikator, tahun, "rankhir"); err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: 200, Status: "OK", Data: "target rankhir berhasil dihapus"})
}

// DeleteTargetPenetapan godoc
// @Summary      Delete Target Penetapan
// @Description  Hapus target penetapan berdasarkan kode indikator dan tahun.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        kode_indikator  path  string  true  "Kode Indikator"
// @Param        tahun           path  string  true  "Tahun"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/target/penetapan/{kode_indikator}/{tahun} [delete]
func (c *MatrixRenjaControllerImpl) DeleteTargetPenetapan(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	kodeIndikator := p.ByName("kode_indikator")
	tahun := p.ByName("tahun")
	if err := c.MatrixRenjaService.DeleteTargetRenjaLayer(r.Context(), kodeIndikator, tahun, "penetapan"); err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: 200, Status: "OK", Data: "target penetapan berhasil dihapus"})
}

// UpsertPaguRanwal godoc
// @Summary      Upsert Pagu Ranwal (nonaktif)
// @Description  Endpoint tidak digunakan. Pagu ranwal mengikuti pagu renstra (tb_pagu jenis 'renstra') per tahun.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body     programkegiatan.PaguRenjaLayerRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/pagu/ranwal [post]
func (c *MatrixRenjaControllerImpl) UpsertPaguRanwal(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// var req programkegiatan.PaguRenjaLayerRequest
	// helper.ReadFromRequestBody(r, &req)
	// if err := c.MatrixRenjaService.UpsertPaguRenjaByJenis(r.Context(), "ranwal", req); err != nil {
	// 	helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
	// 	return
	// }
	// helper.WriteToResponseBody(w, web.WebResponse{Code: 200, Status: "OK", Data: "pagu ranwal berhasil disimpan"})
	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   http.StatusBadRequest, // atau http.StatusMethodNotAllowed (405)
		Status: "BAD REQUEST",
		Data:   "endpoint tidak digunakan: pagu ranwal renja bersumber dari pagu renstra (jenis 'renstra'), bukan input manual di sini",
	})
}

// UpsertPaguRankhir godoc
// @Summary      Upsert Pagu Rankhir (nonaktif)
// @Description  Endpoint tidak digunakan. Pagu rankhir dihitung otomatis dari SUM rincian belanja (rencana aksi). Hanya target rankhir yang dapat diubah via PUT /matrix_renja/target/rankhir.
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body     programkegiatan.PaguRenjaLayerRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/pagu/rankhir [post]
func (c *MatrixRenjaControllerImpl) UpsertPaguRankhir(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   http.StatusBadRequest,
		Status: "BAD REQUEST",
		Data:   "endpoint tidak digunakan: pagu rankhir dihitung dari total rincian belanja, bukan input manual. Gunakan PUT /matrix_renja/target/rankhir untuk mengubah target.",
	})
}

// UpsertPaguPenetapan godoc
// @Summary      Upsert Pagu Penetapan
// @Description  Simpan atau update pagu penetapan per subkegiatan (input manual seperti renstra).
// @Tags         Matrix Renja
// @Accept       json
// @Produce      json
// @Param        request  body     programkegiatan.PaguRenjaLayerRequest  true  "Request Body"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /matrix_renja/pagu/penetapan [post]
func (c *MatrixRenjaControllerImpl) UpsertPaguPenetapan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req programkegiatan.PaguRenjaLayerRequest
	helper.ReadFromRequestBody(r, &req)
	if err := c.MatrixRenjaService.UpsertPaguRenjaByJenis(r.Context(), "penetapan", req); err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error()})
		return
	}
	helper.WriteToResponseBody(w, web.WebResponse{Code: 200, Status: "OK", Data: "pagu penetapan berhasil disimpan"})
}
