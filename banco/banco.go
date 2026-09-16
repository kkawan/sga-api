package banco

import "api-gin/models"

// Por enquanto não estou usando banco de dados de verdade, os dados ficam
// guardados em memória mesmo. Quando o servidor reinicia perde tudo.
var (
	Salas  = []*models.Sala{}
	Alunos = []*models.Aluno{}
	Turmas = []*models.Turma{}
)

// BuscarSala procura uma sala pelo id, se não achar devolve nil
func BuscarSala(id string) *models.Sala {
	for _, sala := range Salas {
		if sala.ID == id {
			return sala
		}
	}
	return nil
}

// BuscarAluno procura um aluno pelo id (matricula), se não achar devolve nil
func BuscarAluno(id string) *models.Aluno {
	for _, aluno := range Alunos {
		if aluno.ID == id {
			return aluno
		}
	}
	return nil
}

// BuscarTurma procura uma turma pelo id, se não achar devolve nil
func BuscarTurma(id string) *models.Turma {
	for _, turma := range Turmas {
		if turma.ID == id {
			return turma
		}
	}
	return nil
}
