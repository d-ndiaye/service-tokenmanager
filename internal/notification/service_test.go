package notification

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"service-token/internal/userToken"
	"testing"
)

func TestNotificationService_SendPushNotification(t *testing.T) {
	mockUserTokens := []userToken.UserToken{
		{Username: "test",
			DeviceToken: []userToken.DeviceToken{
				{Id: "1",
					Platform: 2,
					TenantNotificationConfig: []userToken.TenantNotificationConfig{
						{TenantId: "3",
							Category: []string{"news"},
						},
					},
				},
			},
		},
	}
	repoMock := userToken.NewRepositoryMock(t)
	repoMock.EXPECT().GetByTenantAndCategory("123", "news").Return(mockUserTokens, nil).Once()
	sn := ServiceNotification{
		s: repoMock,
	}
	err := sn.SendPushNotification("123", "Test notification", "news")
	assert.Nil(t, err)
}

func TestNotificationService_SendPushNotification_Error(t *testing.T) {
	mockUserTokens := []userToken.UserToken{}
	repoMock := userToken.NewRepositoryMock(t)
	repoMock.EXPECT().GetByTenantAndCategory("", "").Return(mockUserTokens, fmt.Errorf("error NotifyNews")).Once()
	sn := ServiceNotification{
		s: repoMock,
	}
	err := sn.SendPushNotification("", "", "")
	assert.NotNil(t, err)
}
