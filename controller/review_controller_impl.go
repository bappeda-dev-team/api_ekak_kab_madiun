package controller

import (
	"context"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ReviewControllerImpl struct {
	ReviewService service.ReviewService
}

func NewReviewControllerImpl(reviewService service.ReviewService) *ReviewControllerImpl {
	return &ReviewControllerImpl{
		ReviewService: reviewService,
	}
}

// @Summary      Create Review Pohon Kinerja
// @Description  Membuat review untuk pohon kinerja; id_pohon_kinerja diisi dari path pokinId. created_by diambil dari JWT.
// @Tags         Review Pokin
// @Accept       json
// @Produce      json
// @Param        pokinId  path  int  true  "ID pohon kinerja"
// @Param        request  body  pohonkinerja.ReviewCreateRequest  true  "Payload review (id_pohon_kinerja dari path)"
// @Success      201  {object}  web.WebResponse{data=pohonkinerja.ReviewResponse}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_pokin/create/{pokinId} [post]
func (controller *ReviewControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	reviewCreateRequest := pohonkinerja.ReviewCreateRequest{}
	helper.ReadFromRequestBody(request, &reviewCreateRequest)

	idStr := params.ByName("pokinId")
	id, err := strconv.Atoi(idStr)
	helper.PanicIfError(err)
	reviewCreateRequest.IdPohonKinerja = id

	claims := request.Context().Value(helper.UserInfoKey).(web.JWTClaim)
	ctx := context.WithValue(request.Context(), helper.UserInfoKey, claims)

	reviewResponse, err := controller.ReviewService.Create(ctx, reviewCreateRequest)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
		return
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create review",
		Data:   reviewResponse,
	})
}

// @Summary      Update Review Pohon Kinerja
// @Description  Memperbarui review berdasarkan ID.
// @Tags         Review Pokin
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "ID review"
// @Param        request  body  pohonkinerja.ReviewUpdateRequest  true  "Payload update review"
// @Success      200  {object}  web.WebResponse{data=pohonkinerja.ReviewResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_pokin/update/{id} [put]
func (controller *ReviewControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	paramId := params.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
	}

	reviewUpdateRequest := pohonkinerja.ReviewUpdateRequest{}
	helper.ReadFromRequestBody(request, &reviewUpdateRequest)

	reviewUpdateRequest.Id = id

	reviewResponse, err := controller.ReviewService.Update(request.Context(), reviewUpdateRequest)
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
		Status: "success update review",
		Data:   reviewResponse,
	})
}

// @Summary      Delete Review Pohon Kinerja
// @Description  Menghapus review berdasarkan ID.
// @Tags         Review Pokin
// @Produce      json
// @Param        id  path  int  true  "ID review"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_pokin/delete/{id} [delete]
func (controller *ReviewControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	paramId := params.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
		return
	}

	controller.ReviewService.Delete(request.Context(), id)
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success delete review",
	})
}

// @Summary      Daftar Review per Pohon Kinerja
// @Description  Mendapatkan semua review untuk satu pohon kinerja.
// @Tags         Review Pokin
// @Produce      json
// @Param        pokin_id  path  int  true  "ID pohon kinerja"
// @Success      200  {object}  web.WebResponse{data=[]pohonkinerja.ReviewResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_pokin/findall/{pokin_id} [get]
func (controller *ReviewControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pohonkinerjaId := params.ByName("pokin_id")
	id, err := strconv.Atoi(pohonkinerjaId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
		return
	}

	reviewResponse, err := controller.ReviewService.FindAll(request.Context(), id)
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
		Status: "success get all review",
		Data:   reviewResponse,
	})
}

// @Summary      Detail Review Pohon Kinerja
// @Description  Mendapatkan satu review berdasarkan ID.
// @Tags         Review Pokin
// @Produce      json
// @Param        id  path  int  true  "ID review"
// @Success      200  {object}  web.WebResponse{data=pohonkinerja.ReviewResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_pokin/detail/{id} [get]
func (controller *ReviewControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	paramId := params.ByName("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		})
	}

	reviewResponse, err := controller.ReviewService.FindById(request.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: "success get review by id",
		Data:   reviewResponse,
	})
}

// @Summary      Daftar Review Tematik
// @Description  Mendapatkan review yang dikelompokkan per tematik untuk tahun tertentu.
// @Tags         Review Pokin
// @Produce      json
// @Param        tahun  path  string  true  "Tahun"  example("2025")
// @Success      200  {object}  web.WebResponse{data=[]pohonkinerja.ReviewTematikResponse}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_pokin/tematik/{tahun} [get]
func (controller *ReviewControllerImpl) FindAllReviewByTematik(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")

	reviewResponse, err := controller.ReviewService.FindAllReviewByTematik(request.Context(), tahun)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: fmt.Sprintf("success get review tematik tahun %v", tahun),
		Data:   reviewResponse,
	})
}

// @Summary      Daftar Review OPD
// @Description  Mendapatkan review pohon kinerja OPD berdasarkan kode OPD dan tahun.
// @Tags         Review Pokin
// @Produce      json
// @Param        kode_opd  path  string  true  "Kode OPD"
// @Param        tahun     path  string  true  "Tahun"  example("2025")
// @Success      200  {object}  web.WebResponse{data=[]pohonkinerja.ReviewOpdResponse}
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_pokin/opd/{kode_opd}/{tahun} [get]
func (controller *ReviewControllerImpl) FindAllReviewOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	reviewResponse, err := controller.ReviewService.FindAllReviewOpd(request.Context(), kodeOpd, tahun)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		})
	}

	helper.WriteToResponseBody(writer, web.WebResponse{
		Code:   http.StatusOK,
		Status: fmt.Sprintf("success get review tematik tahun %v", tahun),
		Data:   reviewResponse,
	})
}

// @Summary      Create Review Tujuan OPD
// @Description  Membuat review untuk tujuan OPD; id_tujuan_opd dari path. id_pohon_kinerja default 0, jenis_pokin tidak perlu dikirim.
// @Tags         Review Tujuan OPD
// @Accept       json
// @Produce      json
// @Param        tujuanOpdId  path  int  true  "ID tujuan OPD"
// @Param        request      body  pohonkinerja.ReviewTujuanOpdCreateRequest  true  "Payload review tujuan OPD"
// @Success      201  {object}  web.WebResponse{data=pohonkinerja.ReviewTujuanOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_tujuan_opd/create/{tujuanOpdId} [post]
func (controller *ReviewControllerImpl) CreateTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	req := pohonkinerja.ReviewTujuanOpdCreateRequest{}
	helper.ReadFromRequestBody(request, &req)

	id, err := strconv.Atoi(params.ByName("tujuanOpdId"))
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error(),
		})
		return
	}
	req.IdTujuanOpd = id

	claims := request.Context().Value(helper.UserInfoKey).(web.JWTClaim)
	ctx := context.WithValue(request.Context(), helper.UserInfoKey, claims)

	resp, err := controller.ReviewService.CreateTujuanOpd(ctx, req)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusInternalServerError, Status: "INTERNAL SERVER ERROR", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusCreated, Status: "success create review tujuan opd", Data: resp,
	})
}

// @Summary      Update Review Tujuan OPD
// @Description  Memperbarui review tujuan OPD berdasarkan ID review.
// @Tags         Review Tujuan OPD
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "ID review"
// @Param        request  body  pohonkinerja.ReviewTujuanOpdUpdateRequest  true  "Payload update"
// @Success      200  {object}  web.WebResponse{data=pohonkinerja.ReviewTujuanOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_tujuan_opd/update/{id} [put]
func (controller *ReviewControllerImpl) UpdateTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error(),
		})
		return
	}
	req := pohonkinerja.ReviewTujuanOpdUpdateRequest{}
	helper.ReadFromRequestBody(request, &req)
	req.Id = id

	resp, err := controller.ReviewService.UpdateTujuanOpd(request.Context(), req)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusInternalServerError, Status: "INTERNAL SERVER ERROR", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success update review tujuan opd", Data: resp,
	})
}

// @Summary      Delete Review Tujuan OPD
// @Description  Menghapus review tujuan OPD berdasarkan ID.
// @Tags         Review Tujuan OPD
// @Produce      json
// @Param        id  path  int  true  "ID review"
// @Success      200  {object}  web.WebResponse
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_tujuan_opd/delete/{id} [delete]
func (controller *ReviewControllerImpl) DeleteTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error(),
		})
		return
	}
	if err := controller.ReviewService.DeleteTujuanOpd(request.Context(), id); err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusInternalServerError, Status: "INTERNAL SERVER ERROR", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success delete review tujuan opd",
	})
}

// @Summary      Daftar Review Tujuan OPD
// @Description  Mendapatkan semua review untuk satu tujuan OPD.
// @Tags         Review Tujuan OPD
// @Produce      json
// @Param        tujuan_opd_id  path  int  true  "ID tujuan OPD"
// @Success      200  {object}  web.WebResponse{data=[]pohonkinerja.ReviewTujuanOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_tujuan_opd/findall/{tujuan_opd_id} [get]
func (controller *ReviewControllerImpl) FindAllTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("tujuan_opd_id"))
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error(),
		})
		return
	}
	resp, err := controller.ReviewService.FindAllTujuanOpd(request.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusInternalServerError, Status: "INTERNAL SERVER ERROR", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success get all review tujuan opd", Data: resp,
	})
}

// @Summary      Detail Review Tujuan OPD
// @Description  Mendapatkan satu review tujuan OPD berdasarkan ID review.
// @Tags         Review Tujuan OPD
// @Produce      json
// @Param        id  path  int  true  "ID review"
// @Success      200  {object}  web.WebResponse{data=pohonkinerja.ReviewTujuanOpdResponse}
// @Failure      400  {object}  web.WebResponse
// @Failure      500  {object}  web.WebResponse
// @Security     BearerAuth
// @Router       /review_tujuan_opd/detail/{id} [get]
func (controller *ReviewControllerImpl) FindByIdTujuanOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusBadRequest, Status: "BAD REQUEST", Data: err.Error(),
		})
		return
	}
	resp, err := controller.ReviewService.FindByIdTujuanOpd(request.Context(), id)
	if err != nil {
		helper.WriteToResponseBody(writer, web.WebResponse{
			Code: http.StatusInternalServerError, Status: "INTERNAL SERVER ERROR", Data: err.Error(),
		})
		return
	}
	helper.WriteToResponseBody(writer, web.WebResponse{
		Code: http.StatusOK, Status: "success get review tujuan opd", Data: resp,
	})
}
