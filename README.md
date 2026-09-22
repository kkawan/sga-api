# SGA - Sistema de Gestão de Alocação

Trabalho da matéria de Introdução ao Desenvolvimento Web com Go.

A API serve pra controlar as salas da faculdade, os alunos, as turmas e, principalmente,
pra não deixar duas turmas caírem na mesma sala no mesmo horário. Parti do exemplo
do professor (`exemplo_gin`) e fui incrementando as funcionalidades da especificação.

Feito com Go + Gin. Os dados ficam guardados em memória (uma lista dentro do pacote
`banco`), então **quando o servidor reinicia perde tudo**. Não deu tempo de colocar
banco de dados de verdade.

## Como rodar

```bash
go mod tidy
go run main.go
```

O servidor sobe em `http://localhost:8080`.

Pra testar rápido se está no ar:

```bash
curl http://localhost:8080/api/v1/health
```

## Organização das pastas

```
main.go              -> sobe o servidor e registra as rotas
models/models.go     -> as structs (Sala, Aluno, Turma, Alocacao)
banco/banco.go       -> as listas em memória e as funções de buscar por id
handlers/salas.go    -> endpoints de sala
handlers/alunos.go   -> endpoints de aluno
handlers/turmas.go   -> endpoints de turma, matrícula e alocação
handlers/horarios.go -> funções auxiliares pra tratar horário e dia da semana
```

## Endpoints

| Método | Rota | O que faz |
| --- | --- | --- |
| GET | `/api/v1/health` | status da API, hora do servidor e versão |
| POST | `/api/v1/salas` | cadastra uma sala |
| GET | `/api/v1/salas` | lista as salas |
| GET | `/api/v1/salas/:id/agenda` | mostra a grade de uso da sala |
| POST | `/api/v1/alunos` | cadastra um aluno |
| GET | `/api/v1/alunos` | lista os alunos |
| GET | `/api/v1/alunos/:id` | busca um aluno e as turmas dele |
| POST | `/api/v1/turmas` | cadastra uma turma |
| GET | `/api/v1/turmas` | lista as turmas com qtd de alunos e status da alocação |
| POST | `/api/v1/turmas/:id/alunos` | matricula um aluno na turma |
| GET | `/api/v1/turmas/:id/alunos` | lista os alunos da turma |
| POST | `/api/v1/turmas/:id/alocar` | aloca uma sala pra turma |

## Regras de negócio que eu implementei

**Na matrícula do aluno (`POST /turmas/:id/alunos`)**

- turma ou aluno que não existe -> `404`
- aluno que já está na turma -> `409`
- se a turma já tem sala e o aluno novo estourar a capacidade -> `422`
- se o aluno já tem aula em outra turma no mesmo dia e horário -> `409`

**Na alocação da sala (`POST /turmas/:id/alocar`)**

- turma ou sala que não existe -> `404`
- sala menor que a quantidade de alunos matriculados -> `422`
- sala já ocupada por outra turma no mesmo dia e horário -> `409`

A conta da sobreposição é essa aqui (está em `handlers/horarios.go`):

```go
inicioNovo < fimAntigo && fimNovo > inicioAntigo
```

Ou seja, das 08:00 às 10:00 **não** conflita com 10:00 às 12:00, porque uma
começa exatamente quando a outra acaba.

## Exemplos com curl

Cadastrando uma sala:

```bash
curl -X POST http://localhost:8080/api/v1/salas \
  -H "Content-Type: application/json" \
  -d '{"id":"S01","nome":"Laboratorio 1","capacidade":30,"recursos":["projetor","computadores","ar-condicionado"]}'
```

Cadastrando um aluno:

```bash
curl -X POST http://localhost:8080/api/v1/alunos \
  -H "Content-Type: application/json" \
  -d '{"id":"2024001","nome":"Maria Souza","email":"maria.souza@aluno.edu.br"}'
```

Cadastrando uma turma:

```bash
curl -X POST http://localhost:8080/api/v1/turmas \
  -H "Content-Type: application/json" \
  -d '{"id":"T01","nome":"Web com Go - Noturno","disciplina":"Introducao ao Desenvolvimento Web com Go","professor":"Tiago Ravache"}'
```

Matriculando o aluno na turma:

```bash
curl -X POST http://localhost:8080/api/v1/turmas/T01/alunos \
  -H "Content-Type: application/json" \
  -d '{"aluno_id":"2024001"}'
```

Alocando a sala pra turma:

```bash
curl -X POST http://localhost:8080/api/v1/turmas/T01/alocar \
  -H "Content-Type: application/json" \
  -d '{"sala_id":"S01","dia_semana":"segunda","hora_inicio":"19:00","hora_fim":"22:00"}'
```

Se tentar alocar outra turma na mesma sala das 20:00 às 21:00 na segunda, a API
responde `409` reclamando do conflito.

## Observações

- O campo `dia_semana` aceita: `segunda`, `terca`, `quarta`, `quinta`, `sexta`,
  `sabado` e `domingo`. Eu trato maiúscula/minúscula e o acento de "terça" e
  "sábado" antes de salvar, senão a comparação do conflito não funciona.
- O horário tem que vir no formato `HH:MM`.
- A validação do e-mail é bem simples, eu só verifico se tem `@`.
- Se a turma ainda não tem sala, a matrícula não valida capacidade nem conflito
  de horário, porque nesse momento não existe horário pra comparar.

README gerado por IA :)
