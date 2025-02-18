package service

import "C"
import (
	"errors"
	"fmt"
	"helpdesk/Exception"
	"helpdesk/dto"
	"helpdesk/entity"
	"helpdesk/model"
	"helpdesk/repository"
	"helpdesk/utils"
	"time"
)

type ChamadoService struct {
	chamadoRepository     repository.ChamadoRepository
	balcaoRepository      repository.BalcaoRepository
	atendimentoRepository repository.AtendimentoRepository
}

type AtendimentoService struct {
	atendimentoRepository repository.AtendimentoRepository
	balcaoRepository      repository.BalcaoRepository
}

type Chamado struct {
	chamado entity.ChamadoEntity
}

func NovoAtendimentoService(atendimentoRepo repository.AtendimentoRepository, balcaoRepo repository.BalcaoRepository) *AtendimentoService {
	return &AtendimentoService{
		atendimentoRepository: atendimentoRepo,
		balcaoRepository:      balcaoRepo,
	}
}

func NovoChamadoService(chamadoRepo repository.ChamadoRepository, balcaoRepo repository.BalcaoRepository) *ChamadoService {
	return &ChamadoService{
		chamadoRepository: chamadoRepo,
		balcaoRepository:  balcaoRepo,
	}
}

func (cs *ChamadoService) CriarChamado(chamadosDTO *dto.ChamadoDTO) (*entity.ChamadoEntity, error) {
	if chamadosDTO == nil {
		return nil, errors.New("Chamado não pode ser nulo.")
	}

	if _, err := cs.PegarChamado(chamadosDTO); err != nil {
		return nil, err
	}

	balcao, err := cs.balcaoRepository.FindById(chamadosDTO.IDBalcao)
	if err != nil || balcao == nil {
		return nil, errors.New("Balcão não encontrado.")
	}

	podeAtender, err := cs.BalcaoPodeAtender(balcao)
	if err != nil {
		return nil, err
	}
	if !podeAtender {
		return nil, errors.New("Balcão cheio. O chamado será colocado na fila de espera.")
	}

	chamadoExistente, err := cs.chamadoRepository.FindBySerial(chamadosDTO.SerialNumber)
	if err != nil {
		return nil, err
	}

	if chamadoExistente == nil {
		return nil, errors.New("Nenhum chamado encontrado para o serial fornecido.")
	}

	if chamadoExistente.CustomerID == chamadosDTO.CustomerID {
		if chamadoExistente.StatusChamado != "ABERTO" {
			return nil, &Exception.ConflictException{
				Message: "Já existe um chamado aberto para este serial.",
				Uri:     fmt.Sprintf("/api/chamados/%d", chamadoExistente.ID),
			}
		}
	} else {
		if chamadoExistente.StatusChamado != "RESOLVIDO" {
			return nil, &Exception.ForbiddenException{
				Message: "Este serial já está em atendimento por outro usuário.",
				Uri:     fmt.Sprintf("/api/chamados/%d", chamadoExistente.ID),
			}
		}
	}

	novoChamado := ConvertDTOToEntity(chamadosDTO)

	novoChamado.DataCreation = time.Now()
	novoChamado.DataResolution = time.Time{}
	novoChamado.StatusChamado = string(entity.Aberto)
	novoChamado.Balcao = utils.ConvertBalcaoEntityToBalcao(balcao)
	novoChamado.DeviceID = chamadosDTO.DeviceID
	novoChamado.Motivo = chamadosDTO.Motivo
	novoChamado.UserClient = chamadosDTO.UserClient
	novoChamado.IDBalcao = chamadosDTO.IDBalcao

	chamadoSalvo, err := cs.chamadoRepository.Save(novoChamado)
	if err != nil {
		return nil, fmt.Errorf("Erro ao salvar o chamado: %v", err)
	}

	if err := cs.AcrescentarFilaAtendimento(balcao, chamadoSalvo); err != nil {
		return nil, fmt.Errorf("Erro ao adicionar o chamado na fila de atendimento: %v", err)
	}

	return chamadoSalvo, nil
}

func (cs *ChamadoService) PegarChamado(chamadoDTO *dto.ChamadoDTO) (*dto.ChamadoDTO, error) {
	if chamadoDTO == nil {
		return nil, errors.New("chamado nao pode ser nulo!")
	}

	chamadoExistente, err := cs.chamadoRepository.FindByUsuarioAtendenteAndEstado(chamadoDTO.UserAtendente, entity.StatusChamado(entity.Aberto))
	if err != nil {
		return nil, err
	}
	if len(chamadoExistente) > 0 {
		return nil, errors.New("O atendente já possui um chamado ativo.")
	}
	return chamadoDTO, nil
}

func (cs *ChamadoService) BalcaoPodeAtender(balcao *entity.BalcaoEntity) (bool, error) {
	if balcao == nil {
		return false, errors.New("Balcão não encontrado.")
	}

	const limiteAtendimentos = 5

	qtdAbertos, err := cs.atendimentoRepository.FindOpenByBalcao(balcao.ID)
	if err != nil {
		return false, fmt.Errorf("erro ao buscar atendimentos abertos: %w", err)
	}

	if qtdAbertos < limiteAtendimentos {
		return true, nil
	}

	return false, nil
}

func (cs *ChamadoService) AcrescentarFilaAtendimento(balcao *entity.BalcaoEntity, chamado *entity.ChamadoEntity) error {
	if balcao == nil {
		return fmt.Errorf("Balcão não pode ser nulo")
	}

	if chamado == nil {
		return fmt.Errorf("Chamado não pode ser nulo")
	}

	atendimento := &entity.ListaAtendimento{
		Chamado: chamado,
		Balcao:  balcao,
	}

	if err := cs.atendimentoRepository.Save(atendimento); err != nil {
		return fmt.Errorf("erro ao salvar atendimento: %w", err)
	}

	return nil
}

func (cs *ChamadoService) DiminuirFilaAtendimento(idBalcao int64) error {
	balcao, err := cs.balcaoRepository.FindById(idBalcao)
	if err != nil {
		return fmt.Errorf("Balcão com ID %d não encontrado: %w", idBalcao, err)
	}

	if balcao.FilaAtendimento <= 0 {
		return fmt.Errorf("A fila de atendimento do balcão %d já está vazia ou é inválida (valor: %d).", idBalcao, balcao.FilaAtendimento)
	}

	balcao.FilaAtendimento -= 1

	if _, err := cs.balcaoRepository.Save(*balcao); err != nil {
		return fmt.Errorf("erro ao salvar balcão com ID %d: %w", idBalcao, err)
	}

	return nil
}

func (cs *ChamadoService) ChamadoDetalhado(id int64) (*entity.ChamadoEntity, error) {
	chamado, err := cs.chamadoRepository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar chamado com ID %d: %w", id, err)
	}
	if chamado == nil {
		return nil, fmt.Errorf("Chamado nao encontrado com ID %d", id)
	}
	return chamado, nil
}

func (cs *ChamadoService) ListarChamados(page, size int) ([]entity.ChamadoEntity, error) {
	return cs.chamadoRepository.FindAllPaginated(page, size)
}

func (cs *ChamadoService) ListaChamadosCustomerId(customerId int64) ([]entity.ChamadoEntity, error) {
	chamados, err := cs.chamadoRepository.FindByCustomerId(customerId)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar chamados com ID %d: %w", customerId, err)
	}
	if len(chamados) == 0 {
		return nil, fmt.Errorf("Nenhum chamado encontrado para o customer com ID %d", customerId)
	}
	return chamados, nil
}

func (cs *ChamadoService) EditarChamado(id int64, chamadoDTO *dto.ChamadoDTO) (*entity.ChamadoEntity, error) {
	if chamadoDTO == nil {
		return nil, fmt.Errorf("Chamado nao pode ser nulo!")
	}

	chamadoExistente, err := cs.chamadoRepository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar chamado com ID %d: %w", id, err)
	}
	if chamadoExistente == nil {
		return nil, fmt.Errorf("Chamado não encontrado com ID %d", id)
	}

	if chamadoExistente.StatusChamado == "ABERTO" {
		err = cs.DiminuirFilaAtendimento(chamadoExistente.Balcao.ID)
		if err != nil {
			return nil, fmt.Errorf("Erro ao diminuir fila de atendimento: %w", err)
		}
	}

	chamadoExistente.AlterarChamado(chamadoDTO)

	updatedChamado, err := cs.chamadoRepository.Save(chamadoExistente)
	if err != nil {
		return nil, fmt.Errorf("Erro ao salvar o chamado atualizado: %w", err)
	}

	return updatedChamado, nil
}

func ConvertDTOToEntity(dto *dto.ChamadoDTO) *entity.ChamadoEntity {
	return &entity.ChamadoEntity{
		Chamado: model.Chamado{
			ID:             dto.ID,
			CustomerID:     dto.CustomerID,
			DataCreation:   dto.DataCreation,
			DataResolution: dto.DataResolution,
			DeviceID:       dto.DeviceID,
			SerialNumber:   dto.SerialNumber,
			StatusChamado:  dto.StatusChamado,
			IDBalcao:       dto.IDBalcao,
			Motivo:         dto.Motivo,
			Produto:        dto.Produto,
			UserClient:     dto.UserClient,
			UserAtendente:  dto.UserAtendente,
		},
	}
}

func ConvertEntityToDTO(entity *entity.ChamadoEntity) *dto.ChamadoDTO {
	return &dto.ChamadoDTO{
		Chamado: model.Chamado{
			ID:             entity.ID,
			CustomerID:     entity.CustomerID,
			DataCreation:   entity.DataCreation,
			DataResolution: entity.DataResolution,
			DeviceID:       entity.DeviceID,
			SerialNumber:   entity.SerialNumber,
			StatusChamado:  entity.StatusChamado,
			IDBalcao:       entity.IDBalcao,
			Motivo:         entity.Motivo,
			Produto:        entity.Produto,
			UserClient:     entity.UserClient,
			UserAtendente:  entity.UserAtendente,
		},
	}
}
