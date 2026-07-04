package helper

import (
	"ekak_kabupaten_madiun/model/web"
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError(err)
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError(err)
}

func WriteToResponseBodyWstatus(writer http.ResponseWriter, response web.WebResponse) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(response.Code)
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError(err)
}

func DecodeJSONBody(request *http.Request, dest interface{}) error {
	return json.NewDecoder(request.Body).Decode(dest)
}
func WriteJSON(writer http.ResponseWriter, code int, status string, data interface{}) {
	WriteToResponseBodyWstatus(writer, web.WebResponse{
		Code: code, Status: status, Data: data,
	})
}
func WriteBadRequest(writer http.ResponseWriter, err error) {
	WriteJSON(writer, http.StatusBadRequest, "BAD REQUEST", err.Error())
}
func WriteInternalError(writer http.ResponseWriter, err error) {
	WriteJSON(writer, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err.Error())
}

// WriteServiceError — error dari service; target validation → 400, sisanya → 500
func WriteServiceError(writer http.ResponseWriter, err error) {
	if IsTargetValidationError(err) {
		WriteBadRequest(writer, err)
		return
	}
	WriteInternalError(writer, err)
}
