package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type SubKegiatanControllerImpl struct {
	SubKegiatanService service.SubKegiatanService
}

func NewSubKegiatanControllerImpl(subKegiatanService service.SubKegiatanService) *SubKegiatanControllerImpl {
	return &SubKegiatanControllerImpl{
		SubKegiatanService: subKegiatanService,
	}
}

func (controller *SubKegiatanControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subKegiatanCreateRequest := subkegiatan.SubKegiatanCreateRequest{}
	helper.ReadFromRequestBody(request, &subKegiatanCreateRequest)

	subKegiatanResponse, err := controller.SubKegiatanService.Create(request.Context(), subKegiatanCreateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
		Code:   http.StatusCreated,
		Status: "success create data sub kegiatan",
		Data:   subKegiatanResponse,
	})
}

func (controller *SubKegiatanControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if id == "" {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "ID tidak boleh kosong",
		})
		return
	}

	// Decode request body
	subKegiatanUpdateRequest := subkegiatan.SubKegiatanUpdateRequest{}
	err := json.NewDecoder(request.Body).Decode(&subKegiatanUpdateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Format JSON tidak valid",
		})
		return
	}

	// Set id dari params ke request
	subKegiatanUpdateRequest.Id = id

	// Panggil service untuk update gambaran umum
	subKegiatanResponse, err := controller.SubKegiatanService.Update(request.Context(), subKegiatanUpdateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	// Kirim response
	helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
		Code:   http.StatusOK,
		Status: "success update data sub kegiatan",
		Data:   subKegiatanResponse,
	})
}

func (controller *SubKegiatanControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subKegiatanId := params.ByName("id")
	if subKegiatanId == "" {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "ID sub kegiatan tidak valid",
		})
		return
	}

	subKegiatanResponse, err := controller.SubKegiatanService.FindById(request.Context(), subKegiatanId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
		Code:   http.StatusOK,
		Status: "success get data sub kegiatan",
		Data:   subKegiatanResponse,
	})
}

func (controller *SubKegiatanControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subKegiatanResponses, err := controller.SubKegiatanService.FindAll(request.Context())

	if err != nil {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
		Code:   http.StatusOK,
		Status: "success get data sub kegiatan",
		Data:   subKegiatanResponses,
	})
}

func (controller *SubKegiatanControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subKegiatanId := params.ByName("id")

	err := controller.SubKegiatanService.Delete(request.Context(), subKegiatanId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
		Code:   http.StatusOK,
		Status: "success delete data sub kegiatan",
		Data:   "Data sub kegiatan berhasil dihapus",
	})
}

func (controller *SubKegiatanControllerImpl) FindAllByRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Panggil service untuk mendapatkan data sub kegiatan
	subKegiatanResponses, err := controller.SubKegiatanService.FindAll(request.Context())

	if err != nil {
		helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
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
			NameAction: "Create Subkegiatan",
			Method:     "POST",
			Url:        fmt.Sprintf("%s:%s/subkegiatan/create", host, port),
		},
		{
			NameAction: "Update Subkegiatan",
			Method:     "PUT",
			Url:        fmt.Sprintf("%s:%s/sub_kegiatan/update/:id", host, port),
		},
		{
			NameAction: "Delete Subkegiatan",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/sub_kegiatan/delete/:id", host, port),
		},
		{
			NameAction: "Pilihan Subkegiatan",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/sub_kegiatan/pilihan/:kode_opd", host, port),
		},
		{
			NameAction:  "Create Subkegiatan Yang Dipilih",
			Method:      "POST",
			Url:         fmt.Sprintf("%s:%s/subkegiatanterpilih/create/:rencana_kinerja_id", host, port),
			JenisUsulan: "subkegiatan",
		},
	}

	helper.WriteToResponseBody(writer, web.WebSubKegiatanResponse{
		Code:   http.StatusOK,
		Status: "success get data sub kegiatan",
		Data:   subKegiatanResponses,
		Action: buttonActions,
	})
}

func (controller *SubKegiatanControllerImpl) FindSubKegiatanKAK(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeSubKegiatan := params.ByName("kode_subkegiatan")
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	subKegiatanKAKResponse, err := controller.SubKegiatanService.FindSubKegiatanKAK(request.Context(), kodeOpd, kodeSubKegiatan, tahun)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success find sub kegiatan kak",
		Data:   subKegiatanKAKResponse,
	})
}
