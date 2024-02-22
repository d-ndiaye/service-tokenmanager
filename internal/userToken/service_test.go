package userToken

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"service-token/pkg/dto"
	"testing"
)

func TestServiceUserToken_GetByTenantAndCategory(t *testing.T) {
	mockUserTokens := []UserToken{
		{Username: "test",
			DeviceToken: []DeviceToken{
				{Id: "1",
					Platform: 2,
					TenantNotificationConfig: []TenantNotificationConfig{
						{TenantId: "3",
							Category: []string{"news"},
						},
					},
				},
			},
		},
		{Username: "test",
			DeviceToken: []DeviceToken{
				{Id: "2",
					Platform: 2,
					TenantNotificationConfig: []TenantNotificationConfig{
						{TenantId: "3",
							Category: []string{"news"},
						},
					},
				},
			},
		},
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetByTenantAndCategory("3", "news").Return(mockUserTokens, nil).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	userTokens, err := sut.GetByTenantAndCategory("3", "news")
	assert.Equal(t, 2, len(userTokens))
	assert.Nil(t, err)
}

func TestServiceUserToken_Delete(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("3").Return(nil).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	err := sut.Delete("3")
	assert.Nil(t, err)
}

func TestServiceUserToken_Save(t *testing.T) {
	mockUserTokens := UserToken{
		Username: "test",
		DeviceToken: []DeviceToken{
			{Id: "1",
				Platform: 2,
				TenantNotificationConfig: []TenantNotificationConfig{
					{TenantId: "3",
						Category: []string{"news"},
					},
				},
			},
		},
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Save(mockUserTokens).Return(mockUserTokens, nil).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	userTokens, err := sut.Save(mockUserTokens)
	assert.Equal(t, "test", userTokens.Username)
	assert.Equal(t, mockUserTokens, userTokens)
	assert.Nil(t, err)
}

func TestServiceUserToken_GetTenantConfiguration(t *testing.T) {
	mockUserTokenDto := dto.UserTokenDto{
		Username: "test",
		DeviceToken: dto.DeviceToken{
			Id:       "1",
			Platform: 2,
			TenantNotificationConfig: dto.TenantNotificationConfig{
				TenantId: "3",
				Category: []string{"news"},
			},
		},
	}
	mockUserTokens := UserToken{
		Username: "test",
		DeviceToken: []DeviceToken{
			{Id: "1",
				Platform: 2,
				TenantNotificationConfig: []TenantNotificationConfig{
					{TenantId: "3",
						Category: []string{"news"},
					},
				},
			},
		},
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().FindByUsername("test").Return(mockUserTokens, nil).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	userTokens, err := sut.GetTenantConfiguration("3", "test", "1")
	assert.Equal(t, mockUserTokenDto, userTokens)
	assert.Nil(t, err)
}

func TestServiceUserToken_GetByTenantAndCategory_Error(t *testing.T) {
	var mockUserTokens []UserToken
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().GetByTenantAndCategory("3", "news").Return(mockUserTokens, fmt.Errorf("error GetByTenantAndCategory")).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	userTokens, err := sut.GetByTenantAndCategory("3", "news")
	assert.Equal(t, 0, len(userTokens))
	assert.NotNil(t, err)

}

func TestServiceUserToken_Delete_Error(t *testing.T) {
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Delete("25").Return(fmt.Errorf("error Delete userToken")).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	err := sut.Delete("25")
	assert.NotNil(t, err)
}

func TestServiceUserToken_Save_Error(t *testing.T) {
	userToken := UserToken{}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().Save(userToken).Return(UserToken{}, fmt.Errorf("error saving userToken")).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	userTokens, err := sut.Save(userToken)
	assert.NotNil(t, err)
	assert.Equal(t, UserToken{}, userTokens)

}

func TestServiceUserToken_GetTenantConfiguration_Error(t *testing.T) {
	mockUserTokens := UserToken{
		Username: "test",
		DeviceToken: []DeviceToken{
			{Id: "1",
				Platform: 2,
				TenantNotificationConfig: []TenantNotificationConfig{
					{TenantId: "3",
						Category: []string{"news"},
					},
				},
			},
		},
	}
	repoMock := NewRepositoryMock(t)
	repoMock.EXPECT().FindByUsername("test").Return(mockUserTokens, fmt.Errorf("error FindByUsername")).Once()
	sut := ServiceUserToken{
		repo: repoMock,
	}
	userTokens, err := sut.GetTenantConfiguration("3", "test", "1")
	assert.Equal(t, "", userTokens.Username)
	assert.NotNil(t, err)
}
