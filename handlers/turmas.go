package handlers

import (
	"net/http"

	"api-gin/banco"
	"api-gin/models"

	"github.com/gin-gonic/gin"
)

// struct só pra devolver a turma na listagem junto com a quantidade de alunos
// e se ela já tem sala ou não (a struct Turma não tem esses campos)
type turmaResumo struct {
	ID         string           `json:"id"`
	Nome       string           `json:"nome"`
	Disciplina string           `json:"disciplina"`
	Professor  string           `json:"professor"`
	Ativa      bool             `json:"ativa"`
	QtdAlunos  int              `json:"qtd_alunos"`
	Alocada    bool             `json:"alocada"`
	Alocacao   *models.Alocacao `json:"alocacao"`
}

// o corpo do POST de matrícula
type entradaMatricula struct {
	AlunoID string `json:"aluno_id"`
}

// CriarTurma cadastra uma turma nova.
// POST /api/v1/turmas
func CriarTurma(c *gin.Context) {
	var turma models.Turma

	if err := c.ShouldBindJSON(&turma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "não consegui ler o json: " + err.Error()})
		return
	}

	if turma.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo id é obrigatório"})
		return
	}
	if turma.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo nome é obrigatório"})
		return
	}
	if turma.Disciplina == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo disciplina é obrigatório"})
		return
	}
	if turma.Professor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo professor é obrigatório"})
		return
	}

	if banco.BuscarTurma(turma.ID) != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": "já existe uma turma cadastrada com esse id"})
		return
	}

	// turma sempre nasce vazia e sem sala, mesmo se mandarem esses campos no json
	turma.Alunos = []string{}
	turma.Alocacao = nil
	turma.Ativa = true

	banco.Turmas = append(banco.Turmas, &turma)
	c.JSON(http.StatusCreated, turma)
}

// ListarTurmas devolve as turmas com a quantidade de alunos e o status da alocação.
// GET /api/v1/turmas
func ListarTurmas(c *gin.Context) {
	lista := []turmaResumo{}

	for _, turma := range banco.Turmas {
		lista = append(lista, turmaResumo{
			ID:         turma.ID,
			Nome:       turma.Nome,
			Disciplina: turma.Disciplina,
			Professor:  turma.Professor,
			Ativa:      turma.Ativa,
			QtdAlunos:  len(turma.Alunos),
			Alocada:    turma.Alocacao != nil,
			Alocacao:   turma.Alocacao,
		})
	}

	c.JSON(http.StatusOK, lista)
}

// AdicionarAluno matricula um aluno na turma.
// POST /api/v1/turmas/:id/alunos
func AdicionarAluno(c *gin.Context) {
	turma := banco.BuscarTurma(c.Param("id"))
	if turma == nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "turma não encontrada"})
		return
	}

	var entrada entradaMatricula
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "não consegui ler o json: " + err.Error()})
		return
	}
	if entrada.AlunoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o campo aluno_id é obrigatório"})
		return
	}

	aluno := banco.BuscarAluno(entrada.AlunoID)
	if aluno == nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "aluno não encontrado"})
		return
	}

	// REGRA 1 - não pode matricular o mesmo aluno duas vezes na mesma turma
	for _, matriculado := range turma.Alunos {
		if matriculado == aluno.ID {
			c.JSON(http.StatusConflict, gin.H{"erro": "esse aluno já está matriculado nessa turma"})
			return
		}
	}

	turma.Alunos = append(turma.Alunos, aluno.ID)

	c.JSON(http.StatusCreated, gin.H{
		"mensagem":   "aluno matriculado na turma",
		"turma_id":   turma.ID,
		"aluno_id":   aluno.ID,
		"qtd_alunos": len(turma.Alunos),
	})
}

// ListarAlunosDaTurma mostra os alunos matriculados numa turma.
// GET /api/v1/turmas/:id/alunos
func ListarAlunosDaTurma(c *gin.Context) {
	turma := banco.BuscarTurma(c.Param("id"))
	if turma == nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "turma não encontrada"})
		return
	}

	// na turma eu guardo só o id do aluno, então preciso buscar os dados de cada um
	alunos := []*models.Aluno{}
	for _, id := range turma.Alunos {
		aluno := banco.BuscarAluno(id)
		if aluno != nil {
			alunos = append(alunos, aluno)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"turma_id":   turma.ID,
		"turma":      turma.Nome,
		"qtd_alunos": len(alunos),
		"alunos":     alunos,
	})
}
