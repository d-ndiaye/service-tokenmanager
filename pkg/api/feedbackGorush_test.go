package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"service-token/internal/feedbackGorush"
	"testing"
)

func TestFeedbackGorushHandler_DeleteExpiredToken(t *testing.T) {
	errorLog := feedbackGorush.ErrorLog{
		Type:     "failed-push",
		Platform: "2",
		Token:    "123",
		Message:  "Hello World Test!",
		Error:    "ExpiredToken",
	}
	repoMock := feedbackGorush.NewServiceMock(t)
	repoMock.EXPECT().DeleteExpiredToken(errorLog).Return(nil).Once()
	rr := postFeedbackGorush("failed-push", "2", "123", "Hello World Test!", "ExpiredToken", t, repoMock)
	assert.Equal(t, 200, rr.Code)
}

func TestFeedbackGorushHandler_DeleteExpiredToken_error(t *testing.T) {
	errorLog := feedbackGorush.ErrorLog{
		Type:     "failed-push",
		Platform: "2",
		Token:    "",
		Message:  "Hello World Test!",
		Error:    "ExpiredToken",
	}
	repoMock := feedbackGorush.NewServiceMock(t)
	repoMock.EXPECT().DeleteExpiredToken(errorLog).Return(fmt.Errorf("error DeleteExpiredToken")).Once()
	rr := postFeedbackGorush("failed-push", "2", "", "Hello World Test!", "ExpiredToken", t, repoMock)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func postFeedbackGorush(typePush string, platform string, token string, message string, error string, t *testing.T, feedbackGorushService feedbackGorush.Service) *httptest.ResponseRecorder {
	params := map[string]string{
		"type":     typePush,
		"platform": platform,
		"token":    token,
		"message":  message,
		"error":    error,
	}
	jsonReader := parseJSONBody(t, params)
	return prepareFeedbackGorusTestServer(feedbackGorushService, http.MethodPost, "/", jsonReader)
}

func parseJSONBody(t *testing.T, input interface{}) (reader *bytes.Reader) {
	jsonBody, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	reader = bytes.NewReader(jsonBody)
	return
}

func prepareFeedbackGorusTestServer(mockedService feedbackGorush.Service, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	feedbackGorushHandler := FeedbackGorushHandler{S: mockedService}
	handler := feedbackGorushHandler.Routes()
	handler.ServeHTTP(rr, req)
	return rr
}
