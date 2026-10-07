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

func (c *OutcomeMatrixControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	q := request.URL.Query()
	resp, err := c.OutcomeMatrixService.FindAll(request.Context(), q.Get("kode"), q.Get("kode_subkegiatan"))
	if err != nil {
		writeOutcomeMatrixError(writer, http.StatusInternalServerError, err)
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "OK", Data: resp,
	})
}

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
