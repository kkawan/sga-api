package handlers

import (
	"net/http"

	"api-gin/banco"
	"api-gin/models"

	"github.com/gin-gonic/gin"
)

// CriarSala cadastra uma sala nova no inventário.
// POST /api/v1/salas
func CriarSala(c *gin.Context) {
	var sala models.Sala

	if err := c.ShouldBindJSON(&sala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "não consegui ler o json: " + err.Error()})
		return
	}

	// validação na mão mesmo
	if sala.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo id é obrigatório"})
		return
	}
	if sala.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo nome é obrigatório"})
		return
	}
	if sala.Capacidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "a capacidade tem que ser maior que zero"})
		return
	}

	if banco.BuscarSala(sala.ID) != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": "já existe uma sala cadastrada com esse id"})
		return
	}

	// se não mandar recursos o json fica null, aí eu deixo uma lista vazia
	if sala.Recursos == nil {
		sala.Recursos = []string{}
	}
	sala.Ativa = true

	banco.Salas = append(banco.Salas, &sala)
	c.JSON(http.StatusCreated, sala)
}

// ListarSalas devolve todas as salas cadastradas.
// GET /api/v1/salas
func ListarSalas(c *gin.Context) {
	c.JSON(http.StatusOK, banco.Salas)
}

// AgendaDaSala mostra a grade de uso da sala, pra saber quando ela está livre.
// GET /api/v1/salas/:id/agenda
func AgendaDaSala(c *gin.Context) {
	id := c.Param("id")

	sala := banco.BuscarSala(id)
	if sala == nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "sala não encontrada"})
		return
	}

	// percorro todas as turmas e pego as que estão alocadas nessa sala
	agenda := []gin.H{}
	for _, turma := range banco.Turmas {
		if turma.Alocacao == nil {
			continue
		}
		if turma.Alocacao.SalaID != sala.ID {
			continue
		}

		agenda = append(agenda, gin.H{
			"turma_id":    turma.ID,
			"turma":       turma.Nome,
			"disciplina":  turma.Disciplina,
			"professor":   turma.Professor,
			"dia_semana":  turma.Alocacao.DiaSemana,
			"hora_inicio": turma.Alocacao.HoraInicio,
			"hora_fim":    turma.Alocacao.HoraFim,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"sala_id":    sala.ID,
		"sala":       sala.Nome,
		"capacidade": sala.Capacidade,
		"agenda":     agenda,
	})
}
