package utils

import (
	"helpdesk/entity"
	"helpdesk/model"
)

func ConvertBalcaoEntityToBalcao(balcaoEntity *entity.BalcaoEntity) *model.Balcao {
	return &model.Balcao{
		ID:              balcaoEntity.ID,
		NomeAtendente:   balcaoEntity.NomeAtendente,
		FilaAtendimento: balcaoEntity.FilaAtendimento,
	}
}
