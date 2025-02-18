package entity

type ListaAtendimento struct {
	ID      int64 `json:"ID"`
	Chamado *ChamadoEntity
	Balcao  *BalcaoEntity
}
