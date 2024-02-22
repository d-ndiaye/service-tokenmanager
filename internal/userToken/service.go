package userToken

import (
	"service-token/pkg/dto"
)

type ServiceUserToken struct {
	repo Repository
}

type Service interface {
	Save(userToken UserToken) (UserToken, error)
	GetTenantConfiguration(tenantId string, username string, deviceTokenId string) (dto.UserTokenDto, error)
	GetByTenantAndCategory(tenantId string, category string) ([]UserToken, error)
	Delete(deviceToken string) error
}

func New(repository Repository) Service {
	s := ServiceUserToken{
		repo: repository,
	}
	return s
}

func (st ServiceUserToken) GetTenantConfiguration(tenantId string, username string, deviceTokenId string) (dto.UserTokenDto, error) {
	userToken, err := st.repo.FindByUsername(username)
	if err != nil {
		return dto.UserTokenDto{}, err
	}
	var deviceToken *DeviceToken
	for _, dt := range userToken.DeviceToken {
		if dt.Id == deviceTokenId {
			deviceToken = &dt
			break
		}
	}
	if deviceToken == nil {
		return dto.UserTokenDto{}, err
	}
	var tenantNotificationConfig *TenantNotificationConfig
	for _, tnc := range deviceToken.TenantNotificationConfig {
		if tnc.TenantId == tenantId {
			tenantNotificationConfig = &tnc
		}
	}
	return dto.UserTokenDto{
		Username: username,
		DeviceToken: dto.DeviceToken{
			Id:       deviceToken.Id,
			Platform: deviceToken.Platform,
			TenantNotificationConfig: dto.TenantNotificationConfig{
				TenantId: tenantNotificationConfig.TenantId,
				Category: tenantNotificationConfig.Category,
			},
		},
	}, nil
}

func (st ServiceUserToken) Save(userToken UserToken) (UserToken, error) {
	t, err := st.repo.Save(userToken)
	if err != nil {
		return UserToken{}, err
	}
	return t, nil
}
func (st ServiceUserToken) GetByTenantAndCategory(tenantID string, category string) ([]UserToken, error) {
	ut, err := st.repo.GetByTenantAndCategory(tenantID, category)
	if err != nil {
		return ut, err
	}
	return ut, nil

}
func (st ServiceUserToken) Delete(deviceToken string) error {
	err := st.repo.Delete(deviceToken)
	if err != nil {
		return err
	}
	return nil

}
