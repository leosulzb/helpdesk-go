package controller

import (
	"github.com/gin-gonic/gin"
	"helpdesk/Exception"
	"helpdesk/dto"
	"helpdesk/service"
	"net/http"
	"strconv"
)

type ChamadoController struct {
	ChamadoService *service.ChamadoService
}

func NovoChamadoController(service *service.ChamadoService) *ChamadoController {
	return &ChamadoController{ChamadoService: service}
}

func NovoChamadoServiceController() *service.ChamadoService {
	return &service.ChamadoService{}
}

func (cc *ChamadoController) CriarChamado(c *gin.Context) {
	var chamadoDTO dto.ChamadoDTO

	if err := c.ShouldBindJSON(&chamadoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	chamado, err := cc.ChamadoService.CriarChamado(&chamadoDTO)
	if err != nil {
		switch e := err.(type) {
		case *Exception.ConflictException:
			c.JSON(http.StatusConflict, gin.H{
				"message": e.Message,
				"uri":     e.Uri,
			})
		case *Exception.ForbiddenException:
			c.JSON(http.StatusForbidden, gin.H{
				"message": e.Message,
				"uri":     e.Uri,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Erro interno do servidor: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, chamado)
}

func (cc *ChamadoController) ListarChamados(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	chamados, err := cc.ChamadoService.ListarChamados(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chamados)
}

func (cc *ChamadoController) EditarChamados(c *gin.Context, chamadoService service.ChamadoService) {
	var chamadoDTO dto.ChamadoDTO
	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.ShouldBind(&chamadoDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chamadoAtualizado, err := chamadoService.EditarChamado(idInt64, &chamadoDTO)
	if err != nil {
		if err.Error() == "Chamado nao encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao encontrar o Chamado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, chamadoAtualizado)
}
