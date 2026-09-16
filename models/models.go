package models

// Aqui ficam as structs de todas as entidades do sistema.

// Sala é uma sala de aula ou um laboratório da faculdade.
type Sala struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Capacidade int      `json:"capacidade"`
	Recursos   []string `json:"recursos"`
	Ativa      bool     `json:"ativa"`
}

// Aluno é o aluno cadastrado na instituição.
type Aluno struct {
	ID    string `json:"id"` // numero da matricula
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Alocacao guarda onde e quando a turma vai ter aula.
type Alocacao struct {
	SalaID     string `json:"sala_id"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}

// Turma é a turma academica, com os alunos e a sala dela.
type Turma struct {
	ID         string    `json:"id"`
	Nome       string    `json:"nome"`
	Disciplina string    `json:"disciplina"`
	Professor  string    `json:"professor"`
	Ativa      bool      `json:"ativa"`
	Alunos     []string  `json:"alunos"`   // guardo so o id dos alunos matriculados
	Alocacao   *Alocacao `json:"alocacao"` // fica nulo enquanto a turma nao tem sala
}
