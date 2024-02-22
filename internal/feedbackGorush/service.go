package feedbackGorush

import (
	"fmt"
	"service-token/internal/userToken"
)

type userTokenService struct {
	s userToken.Service
}

type Service interface {
	DeleteExpiredToken(errorLog ErrorLog) error
}

func New(uService userToken.Service) Service {
	us := userTokenService{
		s: uService,
	}
	return us
}

func (uts userTokenService) DeleteExpiredToken(errorLog ErrorLog) error {
	err := uts.s.Delete(errorLog.Token)
	if err != nil {
		return err
	}
	fmt.Println("Tokens with errors deleted successfully.")
	return nil
}
