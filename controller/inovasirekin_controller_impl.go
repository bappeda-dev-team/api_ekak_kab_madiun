package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/inovasirekin"
	"ekak_kabupaten_madiun/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type InovasiRekinControllerImpl struct {
	InovasiRekinService service.InovasiRekinService
}

func NewInovasiRekinControllerImpl(inovasiRekinService service.InovasiRekinService) *InovasiRekinControllerImpl {
	return &InovasiRekinControllerImpl{
		InovasiRekinService: inovasiRekinService,
	}
}

func (controller *InovasiRekinControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil rekinId dari params URL
	rekinId := params.ByName("rencana_kinerja_id")
	if rekinId == "" {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "RekinId tidak boleh kosong",
		})
		return
	}

	// Decode request body
	inovasiRekinCreateRequest := inovasirekin.InovasiRekinCreateRequest{}
	helper.ReadFromRequestBody(request, &inovasiRekinCreateRequest)

	// Set rekinId dari params ke request
	inovasiRekinCreateRequest.RekinId = rekinId

	// Panggil service untuk membuat gambaran umum
	inovasiRekinResponse, err := controller.InovasiRekinService.Create(request.Context(), inovasiRekinCreateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	// Kirim response
	helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
		Code:   http.StatusCreated,
		Status: "success create data inovasi",
		Data:   inovasiRekinResponse,
	})
}

func (controller *InovasiRekinControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil id dari params URL
	id := params.ByName("id")
	if id == "" {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "ID tidak boleh kosong",
		})
		return
	}

	// Decode request body
	inovasiRekinUpdateRequest := inovasirekin.InovasiRekinUpdateRequest{}
	err := json.NewDecoder(request.Body).Decode(&inovasiRekinUpdateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebGambaranUmumResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Format JSON tidak valid",
		})
		return
	}

	// Set id dari params ke request
	inovasiRekinUpdateRequest.Id = id

	// Panggil service untuk update gambaran umum
	inovasiRekinResponse, err := controller.InovasiRekinService.Update(request.Context(), inovasiRekinUpdateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	// Kirim response
	helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
		Code:   http.StatusOK,
		Status: "success update data inovasi",
		Data:   inovasiRekinResponse,
	})
}

func (controller *InovasiRekinControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil id dari params URL
	id := params.ByName("id")
	if id == "" {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "ID tidak boleh kosong",
		})
		return
	}

	// Panggil service untuk menghapus gambaran umum
	err := controller.InovasiRekinService.Delete(request.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	// Kirim response
	helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
		Code:   http.StatusOK,
		Status: "success delete data inovasi",
		Data:   "Inovasi berhasil dihapus",
	})
}

func (controller *InovasiRekinControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil rekinId dari params URL
	rekinId := params.ByName("rencana_kinerja_id")
	if rekinId == "" {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "RekinId tidak boleh kosong",
		})
		return
	}

	// Panggil service untuk mendapatkan semua gambaran umum
	inovasiRekinResponses, err := controller.InovasiRekinService.FindAll(request.Context(), rekinId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	// Kirim response
	helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
		Code:   http.StatusOK,
		Status: "success get data inovasi",
		Data:   inovasiRekinResponses,
	})
}

func (controller *InovasiRekinControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil id dari params URL
	id := params.ByName("id")
	if id == "" {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "ID tidak boleh kosong",
		})
		return
	}

	// Panggil service untuk mendapatkan gambaran umum berdasarkan ID
	inovasiRekinResponse, err := controller.InovasiRekinService.FindById(request.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	// Kirim response
	helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
		Code:   http.StatusOK,
		Status: "success get data inovasi by id",
		Data:   inovasiRekinResponse,
	})
}

func (controller *InovasiRekinControllerImpl) FindAllByRekinId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil rekinId dari params URL
	rekinId := params.ByName("rencana_kinerja_id")
	if rekinId == "" {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "RekinId tidak boleh kosong",
		})
		return
	}

	// Panggil service untuk mendapatkan semua gambaran umum
	inovasiRekinResponses, err := controller.InovasiRekinService.FindAll(request.Context(), rekinId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Create Inovasi",
			Method:     "POST",
			Url:        fmt.Sprintf("%s:%s/inovasi_rekin/create/:rencana_kinerja_id", host, port),
		},
	}

	helper.WriteToResponseBody(writer, web.WebInovasiRekinResponse{
		Code:   http.StatusOK,
		Status: "success get data inovasi",
		Action: buttonActions,
		Data:   inovasiRekinResponses,
	})
}
