package handlers

import (
	"net/http"
	"strings"

	"api-gin/banco"
	"api-gin/models"

	"github.com/gin-gonic/gin"
)

// CriarAluno cadastra um aluno novo.
// POST /api/v1/alunos
func CriarAluno(c *gin.Context) {
	var aluno models.Aluno

	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "não consegui ler o json: " + err.Error()})
		return
	}

	if aluno.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo id (matrícula) é obrigatório"})
		return
	}
	if aluno.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo nome é obrigatório"})
		return
	}
	// não fiz nada muito elaborado pra validar email, só olho se tem o @
	if !strings.Contains(aluno.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "email inválido"})
		return
	}

	if banco.BuscarAluno(aluno.ID) != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": "já existe um aluno com essa matrícula"})
		return
	}

	banco.Alunos = append(banco.Alunos, &aluno)
	c.JSON(http.StatusCreated, aluno)
}

// ListarAlunos devolve todos os alunos cadastrados.
// GET /api/v1/alunos
func ListarAlunos(c *gin.Context) {
	c.JSON(http.StatusOK, banco.Alunos)
}

// BuscarAlunoPorID devolve os dados de um aluno só.
// GET /api/v1/alunos/:id
func BuscarAlunoPorID(c *gin.Context) {
	id := c.Param("id")

	aluno := banco.BuscarAluno(id)
	if aluno == nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "aluno não encontrado"})
		return
	}

	// aproveito e já mostro em quais turmas esse aluno está matriculado
	turmas := []gin.H{}
	for _, turma := range banco.Turmas {
		for _, matriculado := range turma.Alunos {
			if matriculado == aluno.ID {
				turmas = append(turmas, gin.H{
					"turma_id":   turma.ID,
					"turma":      turma.Nome,
					"disciplina": turma.Disciplina,
					"alocacao":   turma.Alocacao,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     aluno.ID,
		"nome":   aluno.Nome,
		"email":  aluno.Email,
		"turmas": turmas,
	})
}
