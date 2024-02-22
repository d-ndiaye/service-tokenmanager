package feedbackGorush

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"service-token/internal/userToken"
	"testing"
)

func TestUserTokenService_DeleteExpiredToken(t *testing.T) {
	errorLog := ErrorLog{
		Type:     "failed-push",
		Platform: "2",
		Token:    "123",
		Message:  "Hello World Test!",
		Error:    "ExpiredToken",
	}
	repoMock := userToken.NewRepositoryMock(t)
	repoMock.EXPECT().Delete(errorLog.Token).Return(nil).Once()
	us := userTokenService{
		s: repoMock,
	}
	err := us.DeleteExpiredToken(errorLog)
	assert.Nil(t, err)
}

func TestUserTokenService_DeleteExpiredToken_Error(t *testing.T) {
	errorLog := ErrorLog{
		Type:     "failed-push",
		Platform: "2",
		Token:    "123",
		Message:  "Hello World Test!",
		Error:    "ExpiredToken",
	}
	repoMock := userToken.NewRepositoryMock(t)
	repoMock.EXPECT().Delete("123").Return(fmt.Errorf("error DeleteExpiredToken")).Once()
	us := userTokenService{
		s: repoMock,
	}
	err := us.DeleteExpiredToken(errorLog)
	assert.NotNil(t, err)
}
