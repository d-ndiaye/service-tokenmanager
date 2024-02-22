package feedbackGorush

import (
	"github.com/jinzhu/copier"
	"service-token/pkg/dto"
)

type ErrorLog struct {
	Type     string `bson:"type"`
	Platform string `bson:"platform"`
	Token    string `bson:"token"`
	Message  string `bson:"message"`
	Error    string `bson:"error"`
}

func (n *ErrorLog) ToDto() (dto dto.ErrorLogDto, err error) {
	err = copier.Copy(&dto, n)
	return dto, err
}

func (n *ErrorLog) FromDto(dto dto.ErrorLogDto) (err error) {
	err = copier.Copy(n, &dto)
	return err
}
