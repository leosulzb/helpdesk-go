package service

import (
	"errors"
	"fmt"
	"helpdesk/dto"
	"helpdesk/entity"
	"helpdesk/model"
	"helpdesk/repository"
)

type BalcaoService struct {
	balcaoRepository repository.BalcaoRepository
}

func (cs *BalcaoService) CadastrarBalcao(balcaoDTO *dto.BalcaoDTO) (*entity.BalcaoEntity, error) {
	if balcaoDTO == nil {
		return nil, errors.New("O balcão nao pode ser nulo!")
	}
	if cs.AtentendePossuiBalcao(balcaoDTO.NomeAtendente) {
		return nil, errors.New("O atendente " + balcaoDTO.NomeAtendente + " já possui um balcão.")
	}

	balcao := entity.BalcaoEntity{
		Balcao: model.Balcao{
			NomeAtendente:   balcaoDTO.NomeAtendente,
			FilaAtendimento: balcaoDTO.FilaAtendimento,
		},
	}

	saveBalcao, err := cs.balcaoRepository.Save(balcao)
	if err != nil {
		return nil, err
	}
	return &saveBalcao, nil
}

func (bs *BalcaoService) AtentendePossuiBalcao(nomeAtendente string) bool {
	balcoes, err := bs.balcaoRepository.FindAll()
	if err != nil {
		return false
	}
	for _, balcao := range balcoes {
		if balcao.NomeAtendente == nomeAtendente {
			return true
		}
	}
	return false
}

func (bs *BalcaoService) ListarBalcoes() ([]entity.BalcaoEntity, error) {
	balcoes, err := bs.balcaoRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return balcoes, nil
}

func (bs *BalcaoService) EditarBalcao(balcaoDTO *dto.BalcaoDTO, id int64) (*dto.BalcaoDTO, error) {
	if balcaoDTO == nil {
		return nil, errors.New("Balcão ou ID não podem ser nulos")
	}

	if balcaoDTO.ID != id {
		return nil, errors.New("O ID do Balcão no DTO não corresponde ao ID fornecido.")
	}

	balcaoExistente, err := bs.balcaoRepository.FindById(id)
	if err != nil {
		return nil, &NotFoundError{ID: int(id)}
	}

	balcaoExistente.NomeAtendente = balcaoDTO.NomeAtendente
	balcaoExistente.FilaAtendimento = balcaoDTO.FilaAtendimento

	_, err = bs.balcaoRepository.Save(*balcaoExistente)
	if err != nil {
		return nil, fmt.Errorf("Erro ao salvar balcão: %w", err)
	}

	balcaoDTO.ID = balcaoExistente.ID
	return balcaoDTO, nil
}

type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("O recurso com ID %d não foi encontrado", e.ID)
}
