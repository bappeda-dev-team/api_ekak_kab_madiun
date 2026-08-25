package controller

import (
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pptk"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PptkControllerImpl struct {
	PptkService service.PptkService
}

func NewPptkControllerImpl(pptkService service.PptkService) *PptkControllerImpl {
	return &PptkControllerImpl{
		PptkService: pptkService,
	}
}
func (controller *PptkControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pptkCreateRequest := pptk.PptkCreateRequest{}
	helper.ReadFromRequestBody(request, &pptkCreateRequest)

	// TODO: guard jika request invalid
	// return 400 Invalid

	pptkResponse, err := controller.PptkService.Create(request.Context(), pptkCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			// TODO: CODE: AMBIL DARI http
			Code: http.StatusInternalServerError,
			// TODO: STATUS: TERJEMAHKAN DARI code
			Status: http.StatusText(http.StatusInternalServerError),
			// TODO: buat nil saja
			Data: err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		// TODO: CODE AMBIL DARI http
		Code:   201,
		Status: "Success Created PPTK",
		Data:   pptkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PptkControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pptkUpdateRequest := pptk.PptkUpdateRequest{}
	helper.ReadFromRequestBody(request, &pptkUpdateRequest)

	idStr := params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	pptkUpdateRequest.Id = id

	pptkResponse, err := controller.PptkService.Update(request.Context(), pptkUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Updated PPTK",
		Data:   pptkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PptkControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pptkId := params.ByName("id")
	id, err := strconv.Atoi(pptkId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.PptkService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Deleted PPTK",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
func (controller *PptkControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pptkId := params.ByName("id")
	id, err := strconv.Atoi(pptkId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	pptkResponse, err := controller.PptkService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Found PPTK",
		Data:   pptkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *PptkControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeSubkegiatan := params.ByName("kode_sub_kegiatan")
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	// Jika kodeOpd atau tahun kosong, kembalikan response null
	if kodeOpd == "" || tahun == "" {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service FindAll
	pptkResponse, err := controller.PptkService.FindAll(request.Context(),kodeSubkegiatan, kodeOpd, tahun)
	if err != nil {
		// Jika tidak ada data, kembalikan response sukses dengan data null
		if err == sql.ErrNoRows {
			webResponse := web.WebResponse{
				Code:   200,
				Status: "OK",
				Data:   nil,
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Untuk error lainnya
		webResponse := web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		}
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Get All PPTK",
		Data:   pptkResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
// func (controller *PptkControllerImpl) FindAllByNip(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
// 	nip := params.ByName("nip")
// 	tahun := params.ByName("tahun")

// 	// Jika kodeOpd atau tahun kosong, kembalikan response null
// 	if nip == "" || tahun == "" {
// 		webResponse := web.WebResponse{
// 			Code:   200,
// 			Status: "OK",
// 			Data:   nil,
// 		}
// 		helper.WriteToResponseBody(writer, webResponse)
// 		return
// 	}

// 	// Panggil service FindAll
// 	pptkResponse, err := controller.PptkService.FindAllByNip(request.Context(), nip, tahun)
// 	if err != nil {
// 		// Jika tidak ada data, kembalikan response sukses dengan data null
// 		if err == sql.ErrNoRows {
// 			webResponse := web.WebResponse{
// 				Code:   200,
// 				Status: "OK",
// 				Data:   nil,
// 			}
// 			helper.WriteToResponseBody(writer, webResponse)
// 			return
// 		}

// 		// Untuk error lainnya
// 		webResponse := web.WebResponse{
// 			Code:   404,
// 			Status: "Not Found",
// 			Data:   err.Error(),
// 		}
// 		writer.WriteHeader(http.StatusNotFound)
// 		helper.WriteToResponseBody(writer, webResponse)
// 		return
// 	}

// 	// Kirim response sukses
// 	webResponse := web.WebResponse{
// 		Code:   200,
// 		Status: "Success Get All PPTK",
// 		Data:   pptkResponse,
// 	}
// 	helper.WriteToResponseBody(writer, webResponse)
// }