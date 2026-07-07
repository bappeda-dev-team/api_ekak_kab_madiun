package controller

import (
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/internal"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type StrategicArahKebijakanPemdaControllerImpl struct {
	StrategicArahKebijakanPemdaService service.StrategicArahKebijakanPemdaService
	IsuStrategisClient internal.IsustrategicClient
}

func NewStrategicArahKebijakanPemdaControllerImpl(strategicArahKebijakanPemdaService service.StrategicArahKebijakanPemdaService, isuStrategisClient internal.IsustrategicClient) *StrategicArahKebijakanPemdaControllerImpl {
	return &StrategicArahKebijakanPemdaControllerImpl{
		StrategicArahKebijakanPemdaService: strategicArahKebijakanPemdaService,
		IsuStrategisClient: isuStrategisClient,
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
func (controller *StrategicArahKebijakanPemdaControllerImpl) FindIsu(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	strategicarahkebijakanResponse, err := controller.IsuStrategisClient.GetDataIsuStrategic(request.Context(), kodeOpd, tahun)
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
		Status: "Success Get All Isu Strategis",
		Data:   strategicarahkebijakanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *StrategicArahKebijakanPemdaControllerImpl) ExportExcel(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")

	if tahunAwal == "" || tahunAkhir == "" {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "tahun_awal dan tahun_akhir wajib diisi",
		}

		writer.WriteHeader(http.StatusBadRequest)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	buffer, err := controller.StrategicArahKebijakanPemdaService.ExportExcel(
		request.Context(),
		tahunAwal,
		tahunAkhir,
	)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}

		writer.WriteHeader(http.StatusInternalServerError)
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	filename := fmt.Sprintf(
		"Strategic_Arah_Kebijakan_%s_%s.xlsx",
		tahunAwal,
		tahunAkhir,
	)

	writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	writer.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))

	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}