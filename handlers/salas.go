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
