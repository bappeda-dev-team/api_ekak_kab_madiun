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
