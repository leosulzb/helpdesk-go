package repository

import (
	"database/sql"
	"fmt"
	"helpdesk/entity"
)

type AtendimentoRepository interface {
	Save(atendimento *entity.ListaAtendimento) error
	FindOpenByBalcao(balcaoId int64) (int64, error)
}
type ListaAtendimentoRepositoryImpl struct {
	db *sql.DB
}

type AtendimentoService struct {
	atendimentoRepository AtendimentoRepository
	ListaAtendimento      *entity.ListaAtendimento
}

func NovoAtendimentoService(repo AtendimentoRepository) *AtendimentoService {
	return &AtendimentoService{
		atendimentoRepository: repo,
	}
}

func NovoListaAtendimentoRepository(db *sql.DB) *ListaAtendimentoRepositoryImpl {
	return &ListaAtendimentoRepositoryImpl{db: db}
}

func (repo *ListaAtendimentoRepositoryImpl) FindOpenByBalcao(balcaoID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) 
	          FROM lista_atendimento 
	          WHERE balcao_id = ? AND chamado_estado != ?`

	err := repo.db.QueryRow(query, balcaoID, "CONCLUIDO").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar atendimentos: %w", err)
	}

	return count, nil
}
