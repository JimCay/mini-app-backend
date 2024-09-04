package util

import (
	"encoding/json"
	"html/template"
	"io"
	"tg-backend/pkg/log"

	"net/http"
	"time"
)

func ResponseWrapper(f HttpHandler) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		result := f(r)
		if result.HasError() {
			responseError(result.Error, w)
			return
		}

		responseOk(result, w)
	}

	return handler
}

func responseError(handleError *HandleError, w http.ResponseWriter) {
	if handleError.Type == ErrorInternal {
		log.Error("Handler Error: %+v\n", handleError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	errorResp := struct {
		Error *HandleError `json:"error"`
	}{handleError}
	responseJson, err := json.Marshal(errorResp)
	if err != nil {
		log.Error("Handler Error: %+v\n json: %+v\n", err, errorResp)
		http.Error(w, "can't encode json response error", handleError.GetHttpStatus())
		return
	}
	w.WriteHeader(handleError.GetHttpStatus())
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJson); err != nil {
		log.Error("error writing response", err)
	}
}

func responseOk(result HandleResult, w http.ResponseWriter) {
	if result.Type == ResponseTypeJson {
		responseJson, err := json.Marshal(result.Payload)
		if err != nil {
			log.Error("can't encode json", err)
			http.Error(w, "can't encode json response error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if n, err := w.Write(responseJson); err != nil {
			log.Error("error writing response", "err", err, "bytesWritten", n)
		}
		return
	}

	if result.Type == ResponseTypeHtml {
		tmpl := result.Payload.(*template.Template)
		if err := tmpl.Execute(w, nil); err != nil {
			log.Error("error executing template", err)
			http.Error(w, "error executing template", http.StatusInternalServerError)
		}
		return
	}

	if result.Type == ResponseTypeJpeg {
		now := time.Now()
		reader := result.Payload.(io.ReadCloser)
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Last-Modified", now.Format(http.TimeFormat))
		w.Header().Set("ETag", "image/jpeg")
		w.Header().Set("Expires", now.Add(time.Hour).Format(http.TimeFormat))
		w.Header().Set("Cache-Control", "public,max-age=86400;")
		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, reader); err != nil {
			log.Error("error writing response", err)
		}
		_ = reader.Close()
		return
	}

	responseJson, _ := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseJson); err != nil {
		log.Error("error writing response", err)
	}
}

func Error(message string, errType errorType) HandleResult {
	return ErrorWith(message, errType, nil)
}

func ErrorWith(message string, errType errorType, err error) HandleResult {
	return HandleResult{
		Error: &HandleError{
			Message: message,
			Type:    errType,
			Err:     err,
		},
	}
}

func Success(payload any) HandleResult {
	return HandleResult{
		Payload: payload,
		Type:    ResponseTypeJson,
	}
}
