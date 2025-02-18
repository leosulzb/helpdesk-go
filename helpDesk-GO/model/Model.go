package model

import (
	"helpdesk/entity"
	"time"
)

type Chamado struct {
	ID             int64     `json:"id"`
	CustomerID     int64     `json:"customer_id"`
	DataCreation   time.Time `json:"data_creation"`
	DataResolution time.Time `json:"data_resolution"`
	DeviceID       string    `json:"device_id"`
	SerialNumber   string    `json:"serial_number"`
	Chamado        string    `json:"chamado"`
	StatusChamado  string    `json:"status_chamado"`
	IDBalcao       int64     `json:"id_balcao"`
	Motivo         string    `json:"motivo"`
	Produto        string    `json:"produto"`
	UserClient     string    `json:"user_client"`
	UserAtendente  string    `json:"user_atendente"`
	Balcao         *Balcao   `json:"balcao"`
}

type Balcao struct {
	NomeAtendente   string `json:"nome_atendente"`
	FilaAtendimento int    `json:"fila_atendimento"`
	ID              int64  `json:"id"`
}

func ConvertBalcaoEntityToBalcao(balcaoEntity *entity.BalcaoEntity) Balcao {
	return Balcao{
		ID:              balcaoEntity.ID,
		NomeAtendente:   balcaoEntity.NomeAtendente,
		FilaAtendimento: balcaoEntity.FilaAtendimento,
	}
}
