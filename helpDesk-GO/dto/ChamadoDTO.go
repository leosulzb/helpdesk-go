package dto

import (
	"helpdesk/model"
)

type ChamadoDTO struct {
	model.Chamado
}
type StatusChamado int

const (
	Aberto StatusChamado = iota
	EmAndamento
	Resolvido
	Fechado
)
