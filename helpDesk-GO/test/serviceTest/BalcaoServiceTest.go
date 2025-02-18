//package serviceTest
//
//import (
//	"errors"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"helpdesk/dto"
//	"helpdesk/entity"
//	"helpdesk/model"
//	"helpdesk/service"
//	"testing"
//)
//
//type MockBalcaoRepository struct {
//	mock.Mock
//}
//
//func (m *MockBalcaoRepository) Save(balcao entity.BalcaoEntity) (entity.BalcaoEntity, error) {
//	args := m.Called(balcao)
//	return args.Get(0).(entity.BalcaoEntity), args.Error(1)
//}
//
//func (m *MockBalcaoRepository) AtendentePossuiBalcao(nomeAtendente string) bool {
//	args := m.Called(nomeAtendente)
//	return args.Bool(0)
//}
//
//func (m *MockBalcaoRepository) FindById(id int) (entity.BalcaoEntity, error) {
//	args := m.Called(id)
//	return args.Get(0).(entity.BalcaoEntity), args.Error(1)
//}
//
//func TestCadastrarBalcao(t *testing.T) {
//	tests := []struct {
//		name          string
//		balcaoDTO     *dto.BalcaoDTO
//		attendExist   bool
//		expectedError string
//	}{
//		{
//			name:          "BalcaoDTO é nil",
//			balcaoDTO:     nil,
//			expectedError: "O balcão nao pode ser nulo!",
//		},
//		{
//			name: "Atendente já possui balcão",
//			balcaoDTO: &dto.BalcaoDTO{
//				NomeAtendente:   "João",
//				FilaAtendimento: 5,
//			},
//			attendExist:   true,
//			expectedError: "O atendente João já possui um balcão.",
//		},
//		{
//			name: "Cadastro bem-sucedido",
//			balcaoDTO: &dto.BalcaoDTO{
//				NomeAtendente:   "Maria",
//				FilaAtendimento: 10,
//			},
//			attendExist:   false,
//			expectedError: "",
//		},
//		{
//			name: "Erro ao salvar balcão",
//			balcaoDTO: &dto.BalcaoDTO{
//				NomeAtendente:   "Carlos",
//				FilaAtendimento: 3,
//			},
//			attendExist:   false,
//			expectedError: "Erro ao salvar o balcão",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mockRepo := new(MockBalcaoRepository)
//			cs := service.BalcaoService{
//				BalcaoRepository: mockRepo,
//			}
//
//			// Mock de Atendente Possui Balcão
//			if tt.balcaoDTO != nil {
//				mockRepo.On("FindAll").Return([]entity.BalcaoEntity{}, nil) // Simula retorno vazio no FindAll
//				mockRepo.On("FindById", mock.Anything).Return(nil, nil)     // Mock para FindById
//			}
//
//			if tt.balcaoDTO != nil && !tt.attendExist {
//				// Simula o comportamento do Save
//				if tt.expectedError == "" {
//					mockRepo.On("Save", mock.Anything).Return(entity.BalcaoEntity{}, nil)
//				} else {
//					mockRepo.On("Save", mock.Anything).Return(entity.BalcaoEntity{}, errors.New(tt.expectedError))
//				}
//			}
//
//			// Chama a função a ser testada
//			result, err := cs.CadastrarBalcao(tt.balcaoDTO)
//
//			// Verifica o resultado esperado
//			if tt.expectedError != "" {
//				assert.Nil(t, result)
//				assert.EqualError(t, err, tt.expectedError)
//			} else {
//				assert.NotNil(t, result)
//				assert.NoError(t, err)
//			}
//
//			// Verifica se as expectativas do mock foram atendidas
//			mockRepo.AssertExpectations(t)
//		})
//	}
//}
//func TestAtentendePossuiBalcao(t *testing.T) {
//	tests := []struct {
//		name          string
//		nomeAtendente string
//		mockReturn    []entity.BalcaoEntity
//		expected      bool
//	}{
//		{
//			name:          "Atendente com balcão",
//			nomeAtendente: "João",
//			mockReturn: []entity.BalcaoEntity{
//				{Balcao: model.Balcao{NomeAtendente: "João"}},
//			},
//			expected: true,
//		},
//		{
//			name:          "Atendente sem balcão",
//			nomeAtendente: "Maria",
//			mockReturn:    []entity.BalcaoEntity{},
//			expected:      false,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mockRepo := new(MockBalcaoRepository)
//			cs := service.BalcaoService{
//				BalcaoRepository: mockRepo,
//			}
//
//			mockRepo.On("FindAll").Return(tt.mockReturn, nil)
//
//			// Chama a função a ser testada
//			result := cs.AtentendePossuiBalcao(tt.nomeAtendente)
//
//			assert.Equal(t, tt.expected, result)
//
//			// Verifica se as expectativas do mock foram atendidas
//			mockRepo.AssertExpectations(t)
//		})
//	}
//}
//
//func TestListarBalcoes(t *testing.T) {
//	mockRepo := new(MockBalcaoRepository)
//	cs := service.BalcaoService{
//		BalcaoRepository: mockRepo,
//	}
//
//	mockRepo.On("FindAll").Return([]entity.BalcaoEntity{
//		{Balcao: model.Balcao{NomeAtendente: "João"}},
//	}, nil)
//
//	result, err := cs.ListarBalcoes()
//	assert.NoError(t, err)
//	assert.Len(t, result, 1)
//	assert.Equal(t, "João", result[0].NomeAtendente)
//
//	mockRepo.AssertExpectations(t)
//}
//
//func TestEditarBalcao(t *testing.T) {
//	tests := []struct {
//		name          string
//		balcaoDTO     *dto.BalcaoDTO
//		expectedError string
//	}{
//		{
//			name:          "DTO é nulo",
//			balcaoDTO:     nil,
//			expectedError: "Balcão ou ID não podem ser nulos",
//		},
//		{
//			name: "ID do DTO não corresponde ao ID fornecido",
//			balcaoDTO: &dto.BalcaoDTO{
//				ID:            1,
//				NomeAtendente: "Carlos",
//			},
//			expectedError: "O ID do Balcão no DTO não corresponde ao ID fornecido.",
//		},
//		{
//			name: "Erro ao encontrar o balcão",
//			balcaoDTO: &dto.BalcaoDTO{
//				ID:            2,
//				NomeAtendente: "Carlos",
//			},
//			expectedError: "O recurso com ID 2 não foi encontrado",
//		},
//		{
//			name: "Edição bem-sucedida",
//			balcaoDTO: &dto.BalcaoDTO{
//				ID:              1,
//				NomeAtendente:   "Carlos",
//				FilaAtendimento: 5,
//			},
//			expectedError: "",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mockRepo := new(MockBalcaoRepository)
//			cs := service.BalcaoService{
//				BalcaoRepository: mockRepo,
//			}
//
//			if tt.balcaoDTO != nil {
//				mockRepo.On("FindById", tt.balcaoDTO.ID).Return(&entity.BalcaoEntity{}, nil)
//			}
//
//			if tt.expectedError == "" {
//				mockRepo.On("Save", mock.Anything).Return(entity.BalcaoEntity{}, nil)
//			}
//
//			// Chama a função a ser testada
//			result, err := cs.EditarBalcao(tt.balcaoDTO, tt.balcaoDTO.ID)
//
//			if tt.expectedError != "" {
//				assert.Nil(t, result)
//				assert.EqualError(t, err, tt.expectedError)
//			} else {
//				assert.NotNil(t, result)
//				assert.NoError(t, err)
//			}
//
//			mockRepo.AssertExpectations(t)
//		})
//	}
//}
