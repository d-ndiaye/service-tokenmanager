package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"service-token/internal/news"
	"testing"
)

func TestNewsHandler_NotifyNews(t *testing.T) {
	mockNews := news.News{
		NewsId:       "1",
		Notification: "test",
		TenantId:     "2",
	}
	repoMock := news.NewServiceMock(t)
	repoMock.EXPECT().NotifyNews(mockNews).Return(nil).Once()
	rr := postNews("1", "test", "2", t, repoMock)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestNewsHandler_NotifyNews_Error(t *testing.T) {
	mockNews := news.News{
		NewsId:       "1",
		Notification: "test",
		TenantId:     "2",
	}
	repoMock := news.NewServiceMock(t)
	repoMock.EXPECT().NotifyNews(mockNews).Return(fmt.Errorf("error NotifyNews")).Once()
	rr := postNews("1", "test", "2", t, repoMock)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func postNews(newsId string, notification string, tenantId string, t *testing.T, newsService news.Service) *httptest.ResponseRecorder {
	params := map[string]string{"id": newsId, "notification": notification, "tenantId": tenantId}
	jsonReader := parseJSONBody(t, params)
	return prepareNewsTestServer(newsService, http.MethodPost, "/", jsonReader)
}

func parseJSONBody(t *testing.T, input interface{}) (reader *bytes.Reader) {
	jsonBody, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	reader = bytes.NewReader(jsonBody)
	return
}

func prepareNewsTestServer(mockedService news.Service, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	newsHandler := NewsHandler{S: mockedService}
	handler := newsHandler.Routes()
	handler.ServeHTTP(rr, req)
	return rr
}
