package main

import (
	"net/http"
	"time"

	"api-gin/handlers"

	"github.com/gin-gonic/gin"
)

// versão da API, uso ela no /health
const versao = "1.0.0"

func main() {

	r := gin.New()

	// Uso dos Middlewares globais nativos e personalizados
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 4. Mapeamento de Rotas sob Grupo Versionado
	v1 := r.Group("/api/v1")
	{
		// Monitoramento da API
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   versao,
			})
		})

		// Domínio de Salas (Rooms)
		v1.POST("/salas", handlers.CriarSala)
		v1.GET("/salas", handlers.ListarSalas)
		v1.GET("/salas/:id/agenda", handlers.AgendaDaSala)

		// Domínio de Alunos (Students)
		v1.POST("/alunos", handlers.CriarAluno)
		v1.GET("/alunos", handlers.ListarAlunos)
		v1.GET("/alunos/:id", handlers.BuscarAlunoPorID)

		// Domínio de Turmas (Classes)
		v1.POST("/turmas", handlers.CriarTurma)
		v1.GET("/turmas", handlers.ListarTurmas)
		v1.POST("/turmas/:id/alunos", handlers.AdicionarAluno)
		v1.GET("/turmas/:id/alunos", handlers.ListarAlunosDaTurma)
		v1.POST("/turmas/:id/alocar", handlers.AlocarSala)
	}

	r.Run(":8080")
}
