package controller

import (
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type StrategicArahKebijakanPemdaControllerImpl struct {
	StrategicArahKebijakanPemdaService service.StrategicArahKebijakanPemdaService
}

func NewStrategicArahKebijakanPemdaControllerImpl(strategicArahKebijakanPemdaService service.StrategicArahKebijakanPemdaService) *StrategicArahKebijakanPemdaControllerImpl {
	return &StrategicArahKebijakanPemdaControllerImpl{
		StrategicArahKebijakanPemdaService: strategicArahKebijakanPemdaService,
	}
}

func (controller *StrategicArahKebijakanPemdaControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")

	// Jika kodeOpd atau tahun kosong, kembalikan response null
	if tahunAwal == "" || tahunAkhir == "" {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service FindAll
	strategicarahkebijakanResponse, err := controller.StrategicArahKebijakanPemdaService.FindAll(request.Context(), tahunAwal, tahunAkhir)
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
		Status: "Success Get All Strategic Arah Kebijakan",
		Data:   strategicarahkebijakanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}