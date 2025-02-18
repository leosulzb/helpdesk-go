package controller

import (
	"github.com/gin-gonic/gin"
	"helpdesk/dto"
	"helpdesk/service"
	"net/http"
	"strconv"
)

type BalcaoController struct {
	BalcaoService service.BalcaoService
}

func NewBalcaoController(balcaoService *service.BalcaoService) *BalcaoController {
	return &BalcaoController{BalcaoService: *balcaoService}
}

func (bc *BalcaoController) CadastrarBalcao(c *gin.Context) {
	var balcaoDTO dto.BalcaoDTO

	if err := c.ShouldBindJSON(&balcaoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	saveBalcao, err := bc.BalcaoService.CadastrarBalcao(&balcaoDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar balcão: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Balcão criado com sucesso!", "data": saveBalcao})
}

func (bc *BalcaoController) ListarBalcao(c *gin.Context) {
	var balcaoDTO dto.BalcaoDTO

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido!"})
		return
	}
	if err := c.ShouldBindJSON(&balcaoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos!", "details": err.Error()})
		return
	}
	balcaoAtualizado, err := bc.BalcaoService.EditarBalcao(&balcaoDTO, id)
	if err != nil {
		if _, ok := err.(*service.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Balcão não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, balcaoAtualizado)
}
