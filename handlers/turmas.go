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

// o corpo do POST de alocação
type entradaAlocacao struct {
	SalaID     string `json:"sala_id"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
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

	// As duas regras de baixo só fazem sentido se a turma já tiver sala e horário,
	// porque antes disso não tem capacidade nem horário pra comparar.
	if turma.Alocacao != nil {

		// REGRA 2 - não pode passar da capacidade da sala que a turma está usando
		sala := banco.BuscarSala(turma.Alocacao.SalaID)
		if sala != nil && len(turma.Alunos)+1 > sala.Capacidade {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"erro":         "a sala alocada para essa turma não comporta mais alunos",
				"sala":         sala.Nome,
				"capacidade":   sala.Capacidade,
				"matriculados": len(turma.Alunos),
			})
			return
		}

		// REGRA 3 - o aluno não pode ter aula em outra turma no mesmo dia e horário
		outra := conflitoDeAgendaDoAluno(aluno.ID, turma)
		if outra != nil {
			c.JSON(http.StatusConflict, gin.H{
				"erro":           "o aluno já tem aula em outra turma nesse mesmo horário",
				"turma_conflito": outra.ID,
				"dia_semana":     outra.Alocacao.DiaSemana,
				"horario":        outra.Alocacao.HoraInicio + " às " + outra.Alocacao.HoraFim,
			})
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

// AlocarSala coloca a turma numa sala em um dia e horário.
// POST /api/v1/turmas/:id/alocar
func AlocarSala(c *gin.Context) {
	turma := banco.BuscarTurma(c.Param("id"))
	if turma == nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "turma não encontrada"})
		return
	}

	var entrada entradaAlocacao
	if err := c.ShouldBindJSON(&entrada); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "não consegui ler o json: " + err.Error()})
		return
	}

	sala := banco.BuscarSala(entrada.SalaID)
	if sala == nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "sala não encontrada"})
		return
	}

	if !sala.Ativa || !turma.Ativa {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erro": "a sala ou a turma não está ativa no sistema"})
		return
	}

	dia := arrumaDia(entrada.DiaSemana)
	if !diaValido(dia) {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dia da semana inválido", "aceitos": diasDaSemana})
		return
	}

	inicio, err := horaParaMinutos(entrada.HoraInicio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "hora_inicio: " + err.Error()})
		return
	}
	fim, err := horaParaMinutos(entrada.HoraFim)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "hora_fim: " + err.Error()})
		return
	}
	if inicio >= fim {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "o horário de início tem que ser antes do horário de término"})
		return
	}

	// REGRA 1 - a sala precisa caber todo mundo que já está matriculado
	if len(turma.Alunos) > sala.Capacidade {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro":         "a capacidade da sala é menor que a quantidade de alunos da turma",
			"capacidade":   sala.Capacidade,
			"matriculados": len(turma.Alunos),
		})
		return
	}

	// REGRA 2 - a sala não pode ter duas turmas no mesmo dia em horário sobreposto
	for _, outra := range banco.Turmas {
		if outra.ID == turma.ID {
			continue // se a turma já estava alocada, ignoro ela mesma pra poder remanejar
		}
		if outra.Alocacao == nil {
			continue
		}
		if outra.Alocacao.SalaID != sala.ID || outra.Alocacao.DiaSemana != dia {
			continue
		}

		inicioOcupado, _ := horaParaMinutos(outra.Alocacao.HoraInicio)
		fimOcupado, _ := horaParaMinutos(outra.Alocacao.HoraFim)

		if temSobreposicao(inicio, fim, inicioOcupado, fimOcupado) {
			c.JSON(http.StatusConflict, gin.H{
				"erro":           "a sala já está ocupada nesse dia e horário",
				"turma_conflito": outra.ID,
				"horario":        outra.Alocacao.HoraInicio + " às " + outra.Alocacao.HoraFim,
			})
			return
		}
	}

	turma.Alocacao = &models.Alocacao{
		SalaID:     sala.ID,
		DiaSemana:  dia,
		HoraInicio: entrada.HoraInicio,
		HoraFim:    entrada.HoraFim,
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "turma alocada na sala com sucesso",
		"turma_id": turma.ID,
		"sala":     sala.Nome,
		"alocacao": turma.Alocacao,
	})
}

// conflitoDeAgendaDoAluno olha as outras turmas em que o aluno está matriculado
// e devolve a primeira que bate no mesmo dia e horário da turma que ele quer
// entrar. Se não achar nenhuma devolve nil.
func conflitoDeAgendaDoAluno(alunoID string, turmaNova *models.Turma) *models.Turma {
	inicioNovo, err := horaParaMinutos(turmaNova.Alocacao.HoraInicio)
	if err != nil {
		return nil
	}
	fimNovo, err := horaParaMinutos(turmaNova.Alocacao.HoraFim)
	if err != nil {
		return nil
	}

	for _, outra := range banco.Turmas {
		if outra.ID == turmaNova.ID || outra.Alocacao == nil {
			continue
		}
		if outra.Alocacao.DiaSemana != turmaNova.Alocacao.DiaSemana {
			continue
		}

		// vejo se o aluno está matriculado nessa outra turma
		estaMatriculado := false
		for _, id := range outra.Alunos {
			if id == alunoID {
				estaMatriculado = true
				break
			}
		}
		if !estaMatriculado {
			continue
		}

		inicioOcupado, _ := horaParaMinutos(outra.Alocacao.HoraInicio)
		fimOcupado, _ := horaParaMinutos(outra.Alocacao.HoraFim)

		if temSobreposicao(inicioNovo, fimNovo, inicioOcupado, fimOcupado) {
			return outra
		}
	}

	return nil
}
