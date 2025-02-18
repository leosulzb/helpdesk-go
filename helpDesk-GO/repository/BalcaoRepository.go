package repository

import (
	"helpdesk/entity"
)

type BalcaoRepository interface {
	FindAll() ([]entity.BalcaoEntity, error)
	Save(balcao entity.BalcaoEntity) (entity.BalcaoEntity, error)
	FindById(id int64) (*entity.BalcaoEntity, error)
	FindByCustomerId(customerId int64) ([]entity.BalcaoEntity, error)
}

type BalcaoRepositoryImpl struct {
}
