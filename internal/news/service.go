package news

import "service-token/internal/notification"

const (
	category string = "news"
)

type ServiceNews struct {
	ns notification.Service
}
type Service interface {
	NotifyNews(news News) error
}

func New(nService notification.Service) Service {
	s := ServiceNews{
		ns: nService,
	}
	return s
}

func (uth ServiceNews) NotifyNews(news News) error {
	err := uth.ns.SendPushNotification(news.TenantId, news.Notification, category)
	if err != nil {
		return err
	}
	return nil
}
