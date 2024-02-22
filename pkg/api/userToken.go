package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"service-token/internal/userToken"
	"service-token/pkg/dto"
)

type UserTokenHandler struct {
	S userToken.Service
}

func (uth *UserTokenHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Put("/", uth.SaveToken)
	router.Get("/{username}/tenant/{tenantId}", uth.GetTenantConfiguration)
	return router
}

func (uth *UserTokenHandler) GetTenantConfiguration(response http.ResponseWriter, request *http.Request) {
	tenantId := chi.URLParam(request, "tenantId")
	username := chi.URLParam(request, "username")
	deviceTokenId := request.Header.Get("X-firebaseToken")
	if deviceTokenId == "" {
		http.Error(response, "Missing required header X-firebaseToken", http.StatusBadRequest)
	}
	dto, err := uth.S.GetTenantConfiguration(tenantId, username, deviceTokenId)
	if err != nil {
		http.Error(response, "could not retrieve userToken", http.StatusInternalServerError)
		return
	}
	render.JSON(response, request, dto)
}

func (uth *UserTokenHandler) SaveToken(response http.ResponseWriter, request *http.Request) {
	dto := dto.UserTokenDto{}
	if err := uth.decodeBody(request, &dto); err != nil {
		http.Error(response, "could not decode body", http.StatusBadRequest)
		return
	}
	d := &userToken.UserToken{}
	err := d.FromDto(dto)
	if err != nil {
		http.Error(response, "could not convert dto to userToken", http.StatusInternalServerError)
		return
	}
	_, err = uth.S.Save(*d)
	if err != nil {
		http.Error(response, "could not Create userToken", http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (uth *UserTokenHandler) decodeBody(request *http.Request, target interface{}) (err error) {
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(target)
	return
}
