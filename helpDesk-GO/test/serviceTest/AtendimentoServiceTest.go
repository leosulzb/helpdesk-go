package serviceTest

import (
	"github.com/stretchr/testify/mock"
	"helpdesk/entity"
)

type MockAtendimentoRepository struct {
	mock.Mock
}

func (m *MockAtendimentoRepository) Save(atendimento *entity.ListaAtendimento) error {
	args := m.Called(atendimento)
	return args.Error(0)
}

func (m *MockAtendimentoRepository) FindOpenByBalcao(balcaoID int) (int, error) {
	args := m.Called(balcaoID)
	return args.Int(0), args.Error(1)
}
