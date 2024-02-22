package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"service-token/internal/feedbackGorush"
	"service-token/pkg/dto"
)

type FeedbackGorushHandler struct {
	S feedbackGorush.Service
}

func (fgh *FeedbackGorushHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", fgh.DeleteExpiredToken)
	return router
}

func (fgh *FeedbackGorushHandler) DeleteExpiredToken(response http.ResponseWriter, request *http.Request) {
	dto := dto.ErrorLogDto{}
	if err := fgh.decodeBody(request, &dto); err != nil {
		http.Error(response, "could not decode body", http.StatusBadRequest)
		return
	}
	n := &feedbackGorush.ErrorLog{}
	err := n.FromDto(dto)
	if err != nil {
		http.Error(response, "could not convert dto to feedbackGorush", http.StatusInternalServerError)
		return
	}
	err = fgh.S.DeleteExpiredToken(*n)
	if err != nil {
		http.Error(response, "could not DeleteExpiredToken", http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func (fgh *FeedbackGorushHandler) decodeBody(request *http.Request, target interface{}) (err error) {
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(target)
	return
}
