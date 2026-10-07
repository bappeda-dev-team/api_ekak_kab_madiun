package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/outcomematrix"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type OutcomeMatrixControllerImpl struct {
	OutcomeMatrixService service.OutcomeMatrixService
}

func NewOutcomeMatrixControllerImpl(outcomeMatrixService service.OutcomeMatrixService) *OutcomeMatrixControllerImpl {
	return &OutcomeMatrixControllerImpl{OutcomeMatrixService: outcomeMatrixService}
}

func writeOutcomeMatrixError(writer http.ResponseWriter, code int, err error) {
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   code,
		Status: http.StatusText(code),
		Data:   err.Error(),
	})
}

// @Summary      Create Outcome Matrix
// @Description  Membuat data outcome matrix baru (kode OPD, kode hierarki, jenis, outcome).
// @Tags         Outcome Matrix
// @Accept       json
// @Produce      json
// @Param        request  body  outcomematrix.OutcomeMatrixCreateRequest  true  "Payload create outcome matrix"
// @Success      201  {object}  web.WebResponse{data=outcomematrix.OutcomeMatrixResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /outcome_matrix/create [post]
func (c *OutcomeMatrixControllerImpl) Create(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var req outcomematrix.OutcomeMatrixCreateRequest
	helper.ReadFromRequestBody(request, &req)
	resp, err := c.OutcomeMatrixService.Create(request.Context(), req)
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusCreated, Status: "success create outcome matrix", Data: resp,
	})
}

// @Summary      Update Outcome Matrix
// @Description  Memperbarui outcome matrix berdasarkan ID path.
// @Tags         Outcome Matrix
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "ID outcome matrix"
// @Param        request  body  outcomematrix.OutcomeMatrixUpdateRequest  true  "Payload update outcome matrix"
// @Success      200  {object}  web.WebResponse{data=outcomematrix.OutcomeMatrixResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /outcome_matrix/update/{id} [put]
func (c *OutcomeMatrixControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	var req outcomematrix.OutcomeMatrixUpdateRequest
	helper.ReadFromRequestBody(request, &req)
	req.Id = id
	resp, err := c.OutcomeMatrixService.Update(request.Context(), req)
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success update outcome matrix", Data: resp,
	})
}

// @Summary      Delete Outcome Matrix
// @Description  Menghapus outcome matrix berdasarkan ID.
// @Tags         Outcome Matrix
// @Produce      json
// @Param        id  path  int  true  "ID outcome matrix"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /outcome_matrix/delete/{id} [delete]
func (c *OutcomeMatrixControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	if err := c.OutcomeMatrixService.Delete(request.Context(), id); err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success delete outcome matrix", Data: nil,
	})
}

// @Summary      Detail Outcome Matrix
// @Description  Mendapatkan satu outcome matrix berdasarkan ID.
// @Tags         Outcome Matrix
// @Produce      json
// @Param        id  path  int  true  "ID outcome matrix"
// @Success      200  {object}  web.WebResponse{data=outcomematrix.OutcomeMatrixResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /outcome_matrix/detail/{id} [get]
func (c *OutcomeMatrixControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	resp, err := c.OutcomeMatrixService.FindById(request.Context(), id)
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "OK", Data: resp,
	})
}

// @Summary      Daftar Outcome Matrix
// @Description  Mendapatkan daftar outcome matrix; filter opsional kode, kode_opd, dan jenis (query kosong = tanpa filter field tersebut).
// @Tags         Outcome Matrix
// @Produce      json
// @Param        kode      query  string  false  "Filter kode hierarki (urusan/program/kegiatan/subkegiatan)"
// @Param        kode_opd  query  string  false  "Filter kode OPD"
// @Param        jenis     query  string  false  "Filter jenis"
// @Success      200  {object}  web.WebResponse{data=[]outcomematrix.OutcomeMatrixResponse}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /outcome_matrix [get]
func (c *OutcomeMatrixControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	q := request.URL.Query()
	resp, err := c.OutcomeMatrixService.FindAll(request.Context(), q.Get("kode"), q.Get("kode_opd"), q.Get("jenis"))
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusInternalServerError, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "OK", Data: resp,
	})
}

// @Summary      Upsert Batch Outcome Matrix
// @Description  id = 0 atau kosong untuk create; id > 0 untuk update baris yang sudah ada.
// @Tags         Outcome Matrix
// @Accept       json
// @Produce      json
// @Param        request  body  []outcomematrix.OutcomeMatrixBatchItemRequest  true  "Array item outcome matrix"
// @Success      200  {object}  web.WebResponse{data=[]outcomematrix.OutcomeMatrixResponse}
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /outcome_matrix/batch [post]
func (c *OutcomeMatrixControllerImpl) UpsertBatch(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var reqs []outcomematrix.OutcomeMatrixBatchItemRequest
	helper.ReadFromRequestBody(request, &reqs)
	resp, err := c.OutcomeMatrixService.UpsertBatch(request.Context(), reqs)
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusBadRequest, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success upsert batch outcome matrix", Data: resp,
	})
}
