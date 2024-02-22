package news

import (
	"github.com/jinzhu/copier"
	"service-token/pkg/dto"
)

type News struct {
	NewsId       string
	Notification string
	TenantId     string
}

func (n *News) ToDto() (dto dto.NewsDto, err error) {
	err = copier.Copy(&dto, n)
	return dto, err
}

func (n *News) FromDto(dto dto.NewsDto) (err error) {
	err = copier.Copy(n, &dto)
	return err
}
