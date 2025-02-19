package entity

import "helpdesk/model"

type BalcaoEntity struct {
	model.Balcao
}
type NewBalcaoEntity struct {
	ID              int64  `json:"id"`
	NomeAtendente   string `json:"nome_atendente"`
	FilaAtendimento int    `json:"fila_atendimento"`
}
