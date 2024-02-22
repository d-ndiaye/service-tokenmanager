package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"service-token/internal/userToken"
)

type ServiceNotification struct {
	s userToken.Service
}
type Service interface {
	SendPushNotification(tenantId string, message string, category string) error
}

func New(uService userToken.Service) Service {
	s := ServiceNotification{
		s: uService,
	}
	return s
}

// SendPushNotification this function sends the push notification using Gorush
// return nil
func (uth ServiceNotification) SendPushNotification(tenantId string, message string, category string) error {
	users, err := uth.s.GetByTenantAndCategory(tenantId, category)
	if err != nil {
		fmt.Println(users)
		return err
	}
	var deviceTokens []string
	for _, user := range users {
		for _, device := range user.DeviceToken {
			deviceTokens = append(deviceTokens, device.Id)
		}
	}
	var platform int
	for _, user := range users {
		for _, device := range user.DeviceToken {
			platform = device.Platform
		}
	}
	payload := createPushPayload(tenantId, message, deviceTokens, platform)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post("http://localhost:8088/api/push", "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func createPushPayload(tenantID, message string, deviceTokens []string, platform int) Payload {
	notifications := []Notification{
		{
			Topic:    "news_" + tenantID,
			Message:  message,
			Tokens:   deviceTokens,
			Platform: platform,
		},
	}
	return Payload{
		Notifications: notifications,
	}
}
