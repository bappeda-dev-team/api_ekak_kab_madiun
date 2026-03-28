package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/sasaranopd"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SasaranOpdControllerImpl struct {
	SasaranOpdService service.SasaranOpdService
}

func NewSasaranOpdControllerImpl(SasaranOpdService service.SasaranOpdService) *SasaranOpdControllerImpl {
	return &SasaranOpdControllerImpl{
		SasaranOpdService: SasaranOpdService,
	}
}

func (controller *SasaranOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	KodeOpd := params.ByName("kode_opd")
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")

	sasaranOpdResponse, err := controller.SasaranOpdService.FindAll(request.Context(), KodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
	} else {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "get all sasaran opd",
			Data:   sasaranOpdResponse,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

func (controller *SasaranOpdControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	helper.PanicIfError(err)

	sasaranOpdResponse, err := controller.SasaranOpdService.FindById(request.Context(), idInt)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
	} else {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "get sasaran opd by id",
			Data:   sasaranOpdResponse,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

// @Summary      Create Sasaran Opd
// @Description  Membuat data sasaran opd baru.
// @Tags         Sasaran Opd Renstra
// @Accept       json
// @Produce      json
// @Param        request  body      sasaranopd.SasaranOpdCreateRequest  true  "Payload Create Sasaran OPD"
// @Success      201  {object}  web.WebResponse{data=sasaranopd.SasaranOpdCreateResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/create [post]
func (controller *SasaranOpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	sasaranOpdCreateRequest := sasaranopd.SasaranOpdCreateRequest{}
	helper.ReadFromRequestBody(request, &sasaranOpdCreateRequest)

	sasaranOpdCreateResponse, err := controller.SasaranOpdService.Create(request.Context(), sasaranOpdCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed create sasaran opd",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create sasaran opd",
		Data:   sasaranOpdCreateResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Update Sasaran Opd Renstra
// @Description  Memperbarui data sasaran opd yang sudah ada.
// @Tags         Sasaran Opd Renstra
// @Accept       json
// @Produce      json
// @Param        sasaranopdId       path      int                              true  "ID Sasaran OPD" example(1)
// @Param        request  body      sasaranopd.SasaranOpdUpdateRequest  true  "Payload  Sasaran OPD"
// @Success      201  {object}  web.WebResponse{data=sasaranopd.SasaranOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/update/{id} [put]
func (controller *SasaranOpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	sasaranOpdUpdateRequest := sasaranopd.SasaranOpdUpdateRequest{}
	helper.ReadFromRequestBody(request, &sasaranOpdUpdateRequest)

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "ID harus berupa angka",
		}
		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	sasaranOpdUpdateRequest.IdSasaranOpd = id

	// Panggil service Update
	sasaranOpdResponse, err := controller.SasaranOpdService.Update(request.Context(), sasaranOpdUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Update Sasaran Opd",
		Data:   sasaranOpdResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Delete Sasaran Opd
// @Description  Menghapus data sasaran opd yang sudah ada.
// @Tags         Sasaran Opd Renstra
// @Accept       json
// @Produce      json
// @Param        id       path      int                              true  "ID Sasaran OPD" example(1)
// @Success      200  {object}  web.WebResponse{data=string}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/delete/{id} [delete]
func (controller *SasaranOpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	err := controller.SasaranOpdService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed delete sasaran opd",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
	} else {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "success delete sasaran opd",
			Data:   id,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

func (controller *SasaranOpdControllerImpl) FindByIdPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idPokinStr := params.ByName("id_pokin")
	tahun := params.ByName("tahun")

	// Validasi parameter tidak boleh kosong
	if idPokinStr == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "ID pohon kinerja tidak boleh kosong",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	if tahun == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Tahun tidak boleh kosong",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Validasi format tahun
	if len(tahun) != 4 {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Format tahun harus 4 digit (YYYY)",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	idPokin, err := strconv.Atoi(idPokinStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "ID pohon kinerja harus berupa angka",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	sasaranOpdResponse, err := controller.SasaranOpdService.FindByIdPokin(request.Context(), idPokin, tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "get sasaran opd by id pokin and tahun",
		Data:   sasaranOpdResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *SasaranOpdControllerImpl) FindIdPokinSasaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idPokinStr := params.ByName("id")
	idPokin, err := strconv.Atoi(idPokinStr)
	helper.PanicIfError(err)

	sasaranOpdResponse, err := controller.SasaranOpdService.FindIdPokinSasaran(request.Context(), idPokin)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
	} else {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "get all sasaran opd by id rencana kinerja",
			Data:   sasaranOpdResponse,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

func (controller *SasaranOpdControllerImpl) FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	jenisPeriode := params.ByName("jenis_periode")

	sasaranOpdResponses, err := controller.SasaranOpdService.FindByTahun(request.Context(), kodeOpd, tahun, jenisPeriode)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	} else {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "success find sasaran opd by tahun",
			Data:   sasaranOpdResponses,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

// GetRenjaRanwal godoc
// @Summary      Sasaran Opd Renstra
// @Description  Mendapatkan data sasaran opd renstra berdasarkan kode OPD dan tahun.
// @Tags         Sasaran Opd Renstra
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun_awal  path     string  true  "Tahun Awal"   example("2025")
// @Param        tahun_akhir  path     string  true  "Tahun Akhir"   example("2026")
// @Param        jenis_periode  path     string  true  "Jenis Periode"   example("RPJMD")
// @Success      200  {object}  web.WebResponse{data=[]sasaranopd.SasaranOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/renstra/{kode_opd}/{tahun_awal}/{tahun_akhir}/{jenis_periode} [get]
func (controller *SasaranOpdControllerImpl) FindSasaranRenstra(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")
	jenisPeriode := params.ByName("jenis_periode")
	sasaranOpdResponses, err := controller.SasaranOpdService.FindSasaranRenstra(
		request.Context(), kodeOpd, tahunAwal, tahunAkhir, jenisPeriode,
	)
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
		Status: "success find sasaran opd renstra",
		Data:   sasaranOpdResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// GetRenjaRanwal godoc
// @Summary      Sasaran Opd Ranwal
// @Description  Mendapatkan data sasaran opd ranwal berdasarkan kode OPD dan tahun.
// @Tags         Sasaran Opd Renja Ranwal
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]sasaranopd.SasaranOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/ranwal/{kode_opd}/{tahun} [get]
func (controller *SasaranOpdControllerImpl) FindSasaranRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	sasaranOpdResponses, err := controller.SasaranOpdService.FindSasaranRanwal(
		request.Context(), kodeOpd, tahun, "RPJMD",
	)
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
		Status: "success find sasaran opd ranwal",
		Data:   sasaranOpdResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// GetRenjaRanwal godoc
// @Summary      Sasaran Opd Renja Rankhir
// @Description  Mendapatkan data sasaran opd rankhir berdasarkan kode OPD dan tahun.
// @Tags         Sasaran Opd Renja Rankhir
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]sasaranopd.SasaranOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/rankhir/{kode_opd}/{tahun} [get]
func (controller *SasaranOpdControllerImpl) FindSasaranRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	sasaranOpdResponses, err := controller.SasaranOpdService.FindSasaranRankhir(
		request.Context(), kodeOpd, tahun, "RPJMD",
	)
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
		Status: "success find sasaran opd rankhir",
		Data:   sasaranOpdResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Create Indikator Sasaran Opd Ranwal
// @Description  Membuat data indikator sasaran opd ranwal baru.
// @Tags         Sasaran Opd Renja Ranwal
// @Accept       json
// @Produce      json
// @Param        sasaranopdId       path      int                              true  "ID Sasaran OPD" example(1)
// @Param        request  body      []sasaranopd.IndikatorCreateRequest  true  "Payload Create Indikator Sasaran OPD"
// @Success      201  {object}  web.WebResponse{data=[]sasaranopd.IndikatorResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/renja/ranwal/indikator/create/{sasaranopdId} [post]
func (controller *SasaranOpdControllerImpl) CreateIndikatorRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	sasaranopdId := params.ByName("sasaranopdId")
	sasaranopdIdInt, err := strconv.Atoi(sasaranopdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	indikatorCreateRequests := []sasaranopd.IndikatorCreateRequest{}
	helper.ReadFromRequestBody(request, &indikatorCreateRequests)

	indikatorCreateResponses, err := controller.SasaranOpdService.CreateRenjaIndikator(request.Context(), sasaranopdIdInt, "ranwal", indikatorCreateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed create indikator",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create indikator",
		Data:   indikatorCreateResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

// @Summary      Update Indikator Sasaran Opd Ranwal
// @Description  Memperbarui data indikator sasaran opd ranwal yang sudah ada.
// @Tags         Sasaran Opd Renja Ranwal
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Param        request  body      []sasaranopd.IndikatorUpdateRequest  true  "Payload Update Indikator Sasaran OPD"
// @Success      200  {object}  web.WebResponse{data=sasaranopd.IndikatorResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/renja/ranwal/indikator/update/{kodeIndikator} [put]
func (controller *SasaranOpdControllerImpl) UpdateIndikatorRanwal(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kodeIndikator")
	indikatorUpdateRequests := sasaranopd.IndikatorUpdateRequest{}
	helper.ReadFromRequestBody(request, &indikatorUpdateRequests)
	indikatorUpdateResponses, err := controller.SasaranOpdService.UpdateRenjaIndikator(request.Context(), kodeIndikator, "ranwal", indikatorUpdateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update indikator",
		Data:   indikatorUpdateResponses,
	})
}

// @Summary      Create Indikator Sasaran Opd Ranwal
// @Description  Membuat data indikator sasaran opd ranwal baru.
// @Tags         Sasaran Opd Renja Rankhir
// @Accept       json
// @Produce      json
// @Param        sasaranopdId       path      int                              true  "ID Sasaran OPD" example(1)
// @Param        request  body      []sasaranopd.IndikatorCreateRequest  true  "Payload Create Indikator Sasaran OPD"
// @Success      201  {object}  web.WebResponse{data=[]sasaranopd.IndikatorResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/renja/rankhir/indikator/create/{sasaranopdId} [post]
func (controller *SasaranOpdControllerImpl) CreateIndikatorRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	sasaranopdId := params.ByName("sasaranopdId")
	sasaranopdIdInt, err := strconv.Atoi(sasaranopdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	indikatorCreateRequests := []sasaranopd.IndikatorCreateRequest{}
	helper.ReadFromRequestBody(request, &indikatorCreateRequests)

	indikatorCreateResponses, err := controller.SasaranOpdService.CreateRenjaIndikator(request.Context(), sasaranopdIdInt, "rankhir", indikatorCreateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed create indikator",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create indikator",
		Data:   indikatorCreateResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

// @Summary      Update Indikator Sasaran Opd Ranwal
// @Description  Memperbarui data indikator sasaran opd ranwal yang sudah ada.
// @Tags         Sasaran Opd Renja Rankhir
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Param        request  body      []sasaranopd.IndikatorUpdateRequest  true  "Payload Update Indikator Sasaran OPD"
// @Success      200  {object}  web.WebResponse{data=sasaranopd.IndikatorResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/renja/rankhir/indikator/update/{kodeIndikator} [put]
func (controller *SasaranOpdControllerImpl) UpdateIndikatorRankhir(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kodeIndikator")
	indikatorUpdateRequests := sasaranopd.IndikatorUpdateRequest{}
	helper.ReadFromRequestBody(request, &indikatorUpdateRequests)
	indikatorUpdateResponses, err := controller.SasaranOpdService.UpdateRenjaIndikator(request.Context(), kodeIndikator, "rankhir", indikatorUpdateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update indikator",
		Data:   indikatorUpdateResponses,
	})
}

// @Summary      Delete Indikator Sasaran Opd Renja
// @Description  Menghapus data indikator sasaran opd renja yang sudah ada berdasarkan Kode Indikator.
// @Tags         Sasaran Opd Renja Ranwal
// @Tags         Sasaran Opd Renja Rankhir
// @Tags         Sasaran Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Success      200      {object}  web.WebResponse{data=string}
// @Failure      500      {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/indikator/delete/{kodeIndikator} [delete]
func (controller *SasaranOpdControllerImpl) DeleteIndikatorTargetRenja(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kodeIndikator")
	err := controller.SasaranOpdService.DeleteRenjaIndikator(request.Context(), kodeIndikator)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success delete indikator",
		Data:   nil,
	})
}

// GetRenjaPenetapan godoc
// @Summary      Sasaran Opd Penetapan
// @Description  Mendapatkan data sasaran opd penetapan berdasarkan kode OPD dan tahun.
// @Tags         Sasaran Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        kode_opd  path     string  true  "Kode OPD"   example("1.01.1.01.0.00.01.0000")
// @Param        tahun     path     string  true  "Tahun"      example("2025")
// @Success      200  {object}  web.WebResponse{data=[]sasaranopd.SasaranOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/penetapan/{kode_opd}/{tahun} [get]
func (controller *SasaranOpdControllerImpl) FindSasaranPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")
	sasaranOpdResponses, err := controller.SasaranOpdService.FindSasaranPenetapan(
		request.Context(), kodeOpd, tahun, "RPJMD",
	)
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
		Status: "success find sasaran opd penetapan",
		Data:   sasaranOpdResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

// @Summary      Create Indikator Sasaran Opd Penetapan
// @Description  Membuat data indikator sasaran opd penetapan baru.
// @Tags         Sasaran Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        sasaranopdId       path      int                              true  "ID Sasaran OPD" example(1)
// @Param        request  body      []sasaranopd.IndikatorCreateRequest  true  "Payload Create Indikator Sasaran OPD"
// @Success      201  {object}  web.WebResponse{data=[]sasaranopd.IndikatorResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/renja/penetapan/indikator/create/{sasaranopdId} [post]
func (controller *SasaranOpdControllerImpl) CreateIndikatorPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	sasaranopdId := params.ByName("sasaranopdId")
	sasaranopdIdInt, err := strconv.Atoi(sasaranopdId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	indikatorCreateRequests := []sasaranopd.IndikatorCreateRequest{}
	helper.ReadFromRequestBody(request, &indikatorCreateRequests)

	indikatorCreateResponses, err := controller.SasaranOpdService.CreateRenjaIndikator(request.Context(), sasaranopdIdInt, "penetapan", indikatorCreateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed create indikator",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create indikator",
		Data:   indikatorCreateResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

// @Summary      Update Indikator Sasaran Opd Penetapan
// @Description  Memperbarui data indikator sasaran opd ranwal yang sudah ada.
// @Tags         Sasaran Opd Renja Penetapan
// @Accept       json
// @Produce      json
// @Param        kodeIndikator       path      string                              true  "Kode Indikator"
// @Param        request  body      []sasaranopd.IndikatorUpdateRequest  true  "Payload Update Indikator Sasaran OPD"
// @Success      200  {object}  web.WebResponse{data=sasaranopd.IndikatorResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /sasaran_opd/renja/penetapan/indikator/update/{kodeIndikator} [put]
func (controller *SasaranOpdControllerImpl) UpdateIndikatorPenetapan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeIndikator := params.ByName("kodeIndikator")
	indikatorUpdateRequests := sasaranopd.IndikatorUpdateRequest{}
	helper.ReadFromRequestBody(request, &indikatorUpdateRequests)
	indikatorUpdateResponses, err := controller.SasaranOpdService.UpdateRenjaIndikator(request.Context(), kodeIndikator, "penetapan", indikatorUpdateRequests)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update indikator",
		Data:   indikatorUpdateResponses,
	})
}
