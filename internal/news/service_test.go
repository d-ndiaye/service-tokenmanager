package news

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"service-token/internal/notification"
	"testing"
)

func TestServiceNews_NotifyNews(t *testing.T) {
	news := News{
		NewsId:       "123",
		Notification: "Test notification",
		TenantId:     "news",
	}
	repoMock := notification.NewServiceMock(t)
	repoMock.EXPECT().SendPushNotification(news.TenantId, news.Notification, category).Return(nil).Once()
	sn := ServiceNews{
		ns: repoMock,
	}
	err := sn.NotifyNews(news)
	assert.Nil(t, err)
}

func TestServiceNews_NotifyNews_Error(t *testing.T) {
	news := News{}
	repoMock := notification.NewServiceMock(t)
	repoMock.EXPECT().SendPushNotification(news.TenantId, news.Notification, category).Return(fmt.Errorf("error NotifyNews")).Once()
	sn := ServiceNews{
		ns: repoMock,
	}
	err := sn.NotifyNews(news)
	assert.NotNil(t, err)
}
