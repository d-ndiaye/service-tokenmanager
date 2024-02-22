package userToken

import (
	"github.com/jinzhu/copier"
	"service-token/pkg/dto"
)

type DeviceToken struct {
	Id                       string                     `bson:"id"`
	Platform                 int                        `bson:"platform"`
	TenantNotificationConfig []TenantNotificationConfig `bson:"tenantNotificationConfig"`
}
type TenantNotificationConfig struct {
	TenantId string   `bson:"tenantId"`
	Category []string `bson:"category"`
}
type UserToken struct {
	Username    string        `bson:"username"`
	DeviceToken []DeviceToken `bson:"deviceToken"`
}

func (t *UserToken) ToDto() (dto dto.UserTokenDto, err error) {
	err = copier.Copy(&dto, t)
	return dto, err
}

func (t *UserToken) FromDto(dto dto.UserTokenDto) (err error) {
	err = copier.Copy(t, &dto)
	return err
}
