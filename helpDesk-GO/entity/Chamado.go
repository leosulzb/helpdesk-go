package entity

import (
	"helpdesk/dto"
	"helpdesk/model"
)

type ChamadoEntity struct {
	model.Chamado
}

type StatusChamado int

const (
	Aberto StatusChamado = iota
	EmAndamento
	Resolvido
	Fechado
)

func (c *ChamadoEntity) AlterarChamado(dto *dto.ChamadoDTO) {
	c.CustomerID = dto.CustomerID
	c.DataCreation = dto.DataCreation
	c.DataResolution = dto.DataResolution
	c.DeviceID = dto.DeviceID
	c.SerialNumber = dto.SerialNumber
	c.StatusChamado = dto.StatusChamado
	c.IDBalcao = dto.IDBalcao
	c.Motivo = dto.Motivo
	c.Produto = dto.Produto
	c.UserClient = dto.UserClient
	c.UserAtendente = dto.UserAtendente
}

type ChamadoEntity1 struct {
	ID               int64         `json:"id"`
	CustomerID       int64         `json:"customer_id"`
	SerialNumber     string        `json:"serial_number"`
	Produto          string        `json:"produto"`
	StatusChamado    StatusChamado `json:"estado"`
	UsuarioAtendente string        `json:"usuario_atendente"`
	Balcao           model.Balcao  `json:"balcao"`
}
