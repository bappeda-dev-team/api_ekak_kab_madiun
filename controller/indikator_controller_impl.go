package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/indikator"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IndikatorControllerImpl struct {
	IkmService service.IkmService
}

func NewIndikatorControllerImpl(ikmService service.IkmService) *IndikatorControllerImpl {
	return &IndikatorControllerImpl{
		IkmService: ikmService,
	}
}

// @Summary     Get IKM by Periode
// @Description Get all indikator IKM berdasarkan tahun awal dan akhir
// @Tags        Indikator IKM
// @Accept      json
// @Produce     json
// @Param       tahun_awal   query string true "Tahun Awal"
// @Param       tahun_akhir  query string true "Tahun Akhir"
// @Success     200 {object} web.WebResponse{data=[]indikator.IkmResponse}
// @Failure     400 {object} web.WebResponse
// @Router      /ikm [GET]
func (c *IndikatorControllerImpl) FindAllPeriode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	tahunAwal := r.URL.Query().Get("tahun_awal")
	tahunAkhir := r.URL.Query().Get("tahun_akhir")

	if tahunAwal == "" || tahunAkhir == "" {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "tahun_awal dan tahun_akhir wajib diisi",
		})
		return
	}

	result, err := c.IkmService.FindAllByPeriode(r.Context(), tahunAwal, tahunAkhir)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   result,
	})
}

// @Summary     Get IKM by ID
// @Description Get indikator IKM berdasarkan ID
// @Tags        Indikator IKM
// @Produce     json
// @Param       id path string true "ID IKM"
// @Success     200 {object} web.WebResponse{data=indikator.IkmResponse}
// @Failure     404 {object} web.WebResponse
// @Router      /ikm/{id} [GET]
func (c *IndikatorControllerImpl) FindById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	id := params.ByName("id")
	if id == "" {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   "id wajib diisi",
		})
		return
	}

	result, err := c.IkmService.FindById(r.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   404,
			Status: "Not Found",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   result,
	})
}

// @Summary     Create IKM
// @Description Membuat indikator IKM baru
// @Tags        Indikator IKM
// @Accept      json
// @Produce     json
// @Param       body body indikator.IkmRequest true "Request Body"
// @Success     200 {object} web.WebResponse{data=indikator.IkmResponse}
// @Failure     400 {object} web.WebResponse
// @Router      /ikm [POST]
func (c *IndikatorControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	request := indikator.IkmRequest{}
	helper.ReadFromRequestBody(r, &request)

	result, err := c.IkmService.Create(r.Context(), request)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   result,
	})
}

// @Summary     Update IKM
// @Description Update indikator IKM berdasarkan ID
// @Tags        Indikator IKM
// @Accept      json
// @Produce     json
// @Param       id path string true "ID IKM"
// @Param       body body indikator.IkmRequest true "Request Body"
// @Success     200 {object} web.WebResponse{data=indikator.IkmResponse}
// @Failure     400 {object} web.WebResponse
// @Router      /ikm/{id} [PUT]
func (c *IndikatorControllerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	id := params.ByName("id")

	request := indikator.IkmRequest{}
	helper.ReadFromRequestBody(r, &request)

	result, err := c.IkmService.Update(r.Context(), request, id)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   result,
	})
}

// @Summary     Delete IKM
// @Description Hapus indikator IKM berdasarkan ID
// @Tags        Indikator IKM
// @Produce     json
// @Param       id path string true "ID IKM"
// @Success     200 {object} web.WebResponse
// @Failure     400 {object} web.WebResponse
// @Router      /ikm/{id} [DELETE]
func (c *IndikatorControllerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	id := params.ByName("id")

	err := c.IkmService.Delete(r.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(w, web.WebResponse{
			Code:   400,
			Status: "Bad Request",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(w, web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   "berhasil dihapus",
	})
}
