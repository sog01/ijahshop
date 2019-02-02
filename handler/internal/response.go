package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
)

// ConstructRespSucces to construct success resp
func ConstructRespSucces(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	type List struct {
		List interface{} `json:"list"`
	}

	// reflection need to check data is slice or not
	s := reflect.ValueOf(data)

	constructData := make(map[string]interface{})
	if s.Kind() != reflect.Slice {
		constructData[name] = data
	} else {
		constructData[name] = List{
			List: data,
		}
	}

	return json.NewEncoder(w).Encode(constructData)
}

// ConstructRespSuccesWithMeta to construct success resp with meta
func ConstructRespSuccesWithMeta(w http.ResponseWriter, name string, data interface{}, meta interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	constructData := map[string]interface{}{
		name: struct {
			List interface{} `json:"list"`
			Meta interface{} `json:"meta"`
		}{
			List: data,
			Meta: meta,
		},
	}

	return json.NewEncoder(w).Encode(constructData)
}

// ConstructRespError to construct error resp
func ConstructRespError(w http.ResponseWriter, errorCode int, message string) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)

	type respError struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	constructData := map[string]respError{
		"error": respError{
			Message: message,
			Code:    errorCode,
		},
	}

	return json.NewEncoder(w).Encode(constructData)
}

// ConstructRespErrorWithDetail to construct error resp with detail
func ConstructRespErrorWithDetail(w http.ResponseWriter, errorCode int, message string, errorDetail interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)

	type respError struct {
		Message     string      `json:"message"`
		Code        int         `json:"code"`
		ErrorDetail interface{} `json:"error_detail"`
	}

	constructData := map[string]respError{
		"error": respError{
			Message:     message,
			Code:        errorCode,
			ErrorDetail: errorDetail,
		},
	}

	return json.NewEncoder(w).Encode(constructData)
}

// DownloadFile will download file from files
func DownloadFile(w http.ResponseWriter, filename string) error {

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	out, err := os.Open("files/data/" + filename)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(w, out)
	return err
}
