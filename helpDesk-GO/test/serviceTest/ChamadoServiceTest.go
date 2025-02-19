package serviceTest

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"helpdesk/dto"
	"helpdesk/entity"
	"helpdesk/service"
	"testing"
)

type MockChamadoRepository struct {
	mock.Mock
}

func (m *MockChamadoRepository) Save(chamado entity.ChamadoEntity) (entity.ChamadoEntity, error) {
	args := m.Called(chamado)
	return args.Get(0).(entity.ChamadoEntity), args.Error(1)
}

func (m *MockChamadoRepository) FindBySerial(serial string) (*entity.ChamadoEntity, error) {
	args := m.Called(serial)
	return args.Get(0).(*entity.ChamadoEntity), args.Error(1)
}

func (m *MockChamadoRepository) FindByUsuarioAtendenteAndEstado(userAtendente string, status entity.StatusChamado) ([]entity.ChamadoEntity, error) {
	args := m.Called(userAtendente, status)
	return args.Get(0).([]entity.ChamadoEntity), args.Error(1)
}

func (m *MockBalcaoRepository) FindByID(id int) (*entity.BalcaoEntity, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.BalcaoEntity), args.Error(1)
}

func TestCriarChamado(t *testing.T) {
	tests := []struct {
		name          string
		chamadoDTO    *dto.ChamadoDTO
		mockSetup     func(*MockChamadoRepository, *MockBalcaoRepository, *MockAtendimentoRepository)
		expectedError string
	}{
		{
			name:       "ChamadoDTO é nulo",
			chamadoDTO: nil,
			mockSetup: func(chamadoRepo *MockChamadoRepository, balcaoRepo *MockBalcaoRepository, atendimentoRepo *MockAtendimentoRepository) {
			},
			expectedError: "Chamado não pode ser nulo.",
		},
		{
			name:       "Balcão não encontrado",
			chamadoDTO: &dto.ChamadoDTO{IDBalcao: 1},
			mockSetup: func(chamadoRepo *MockChamadoRepository, balcaoRepo *MockBalcaoRepository, atendimentoRepo *MockAtendimentoRepository) {
				balcaoRepo.On("FindById", 1).Return(nil, errors.New("Balcão não encontrado"))
			},
			expectedError: "Balcão não encontrado.",
		},
		{
			name:       "Balcão não pode atender",
			chamadoDTO: &dto.ChamadoDTO{IDBalcao: 1},
			mockSetup: func(chamadoRepo *MockChamadoRepository, balcaoRepo *MockBalcaoRepository, atendimentoRepo *MockAtendimentoRepository) {
				balcao := &entity.BalcaoEntity{}
				balcaoRepo.On("FindById", 1).Return(balcao, nil)
				chamadoRepo.On("FindBySerial", mock.Anything).Return(nil, nil) // Nenhum chamado com serial fornecido
				atendimentoRepo.On("FindOpenByBalcao", 1).Return(6, nil)       // Limite de atendimentos alcançado
			},
			expectedError: "Balcão cheio. O chamado será colocado na fila de espera.",
		},
		{
			name: "Chamado já existe para o serial",
			chamadoDTO: &dto.ChamadoDTO{
				SerialNumber: "123456",
				CustomerID:   1,
			},
			mockSetup: func(chamadoRepo *MockChamadoRepository, balcaoRepo *MockBalcaoRepository, atendimentoRepo *MockAtendimentoRepository) {
				chamadoRepo.On("FindBySerial", "123456").Return(&entity.ChamadoEntity{CustomerID: 1, StatusChamado: "ABERTO"}, nil)
			},
			expectedError: "Já existe um chamado aberto para este serial.",
		},
		{
			name: "Chamado criado com sucesso",
			chamadoDTO: &dto.ChamadoDTO{
				SerialNumber: "123456",
				CustomerID:   1,
				IDBalcao:     1,
			},
			mockSetup: func(chamadoRepo *MockChamadoRepository, balcaoRepo *MockBalcaoRepository, atendimentoRepo *MockAtendimentoRepository) {
				chamadoRepo.On("FindBySerial", "123456").Return(nil, nil) // Nenhum chamado com serial fornecido
				balcao := &entity.BalcaoEntity{ID: 1}
				balcaoRepo.On("FindById", 1).Return(balcao, nil)
				atendimentoRepo.On("FindOpenByBalcao", 1).Return(3, nil) // Aceita o atendimento
				chamadoRepo.On("Save", mock.Anything).Return(entity.ChamadoEntity{}, nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChamadoRepo := new(MockChamadoRepository)
			mockBalcaoRepo := new(MockBalcaoRepository)
			mockAtendimentoRepo := new(MockAtendimentoRepository)

			tt.mockSetup(mockChamadoRepo, mockBalcaoRepo, mockAtendimentoRepo)

			cs := service.NovoChamadoService(mockChamadoRepo, mockBalcaoRepo)

			result, err := cs.CriarChamado(tt.chamadoDTO)

			if tt.expectedError != "" {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NotNil(t, result)
				assert.NoError(t, err)
			}

			mockChamadoRepo.AssertExpectations(t)
			mockBalcaoRepo.AssertExpectations(t)
			mockAtendimentoRepo.AssertExpectations(t)
		})
	}
}
