package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"service-token/internal/news"
	"service-token/pkg/dto"
)

type NewsHandler struct {
	S news.Service
}

func (nh *NewsHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", nh.NotifyNews)
	return router
}

func (nh *NewsHandler) NotifyNews(response http.ResponseWriter, request *http.Request) {
	dto := dto.NewsDto{}
	if err := nh.decodeBody(request, &dto); err != nil {
		http.Error(response, "could not decode body", http.StatusBadRequest)
		return
	}
	n := &news.News{}
	err := n.FromDto(dto)
	if err != nil {
		http.Error(response, "could not convert dto to news", http.StatusInternalServerError)
		return
	}
	err = nh.S.NotifyNews(*n)
	if err != nil {
		http.Error(response, "could not NotifyNews", http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (nh *NewsHandler) decodeBody(request *http.Request, target interface{}) (err error) {
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(target)
	return
}
