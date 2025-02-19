package repository

import (
	"database/sql"
	"helpdesk/dto"
	"helpdesk/entity"
)

type ChamadoRepository interface {
	FindAll() ([]entity.ChamadoEntity, error)
	Save(chamado *entity.ChamadoEntity) (*entity.ChamadoEntity, error)
	FindById(id int64) (*entity.ChamadoEntity, error)
	FindByCustomerId(customerId int64) ([]entity.ChamadoEntity, error)
	FindByUsuarioAtendenteAndEstado(usuarioAtendente string, estado entity.StatusChamado) ([]entity.ChamadoEntity, error)
	FindByBalcaoAndStatus(balcao entity.BalcaoEntity, status dto.StatusChamado) ([]entity.ChamadoEntity, error)
	FindBySerial(serial string) (*entity.ChamadoEntity, error)
	FindAllPaginated(page int, size int) ([]entity.ChamadoEntity, error)
}

type ChamadoRepositoryImpl struct {
	db *sql.DB
}

func NewChamadoRepository(db *sql.DB) *ChamadoRepositoryImpl {
	return &ChamadoRepositoryImpl{db: db}
}

func (repo ChamadoRepositoryImpl) FindBySerial(serial string) (dto.ChamadoDTO, error) {
	var chamado dto.ChamadoDTO
	query := "SELECT id, serial_number FROM chamados WHERE serial_number = ?"

	row := repo.db.QueryRow(query, serial)
	err := row.Scan(&chamado.ID, &chamado.SerialNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.ChamadoDTO{}, nil
		}
		return dto.ChamadoDTO{}, err
	}
	return chamado, nil
}

func (repo *ChamadoRepositoryImpl) FindAllPaginated(page int, size int) ([]dto.ChamadoDTO, error) {
	offset := page * size
	query := "SELECT id, serial_number, customer_id FROM chamados LIMIT ? OFFSET ?"

	rows, err := repo.db.Query(query, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chamados []dto.ChamadoDTO
	for rows.Next() {
		var chamado dto.ChamadoDTO
		if err := rows.Scan(&chamado.ID, &chamado.SerialNumber, &chamado.CustomerID); err != nil {
			return nil, err
		}
		chamados = append(chamados, chamado)
	}

	return chamados, nil
}
