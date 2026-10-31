package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/1GabrielReis/sistema-chamados/internal/database"
	"github.com/1GabrielReis/sistema-chamados/internal/handler"
	"github.com/1GabrielReis/sistema-chamados/internal/repository/postgres"
	"github.com/1GabrielReis/sistema-chamados/internal/service"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("a variável DATABASE_URL não foi definida")
		return
	}

	pool, err := database.Conectar(ctx, databaseURL)
	if err != nil {
		fmt.Printf("erro ao conectar com o banco: %v\n", err)
		return
	}
	defer pool.Close()

	pessoaRepository := postgres.NovoPessoaRepository(pool)
	chamadoRepository := postgres.NovoChamadoRepository(pool)

	pessoaService := service.NovoPessoaService(pessoaRepository)
	chamadoService := service.NovoChamadoService(chamadoRepository, pessoaRepository)

	pessoaHandler := handler.NovoPessoaHandler(pessoaService)
	chamadoHandler := handler.NovoChamadoHandler(chamadoService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	mux.HandleFunc("POST /pessoas", pessoaHandler.Cadastrar)
	mux.HandleFunc("GET /pessoas/{codigo}", pessoaHandler.BuscarPorCodigoPublico)
	mux.HandleFunc("GET /pessoas/responsaveis", pessoaHandler.ListarResponsaveis)
	mux.HandleFunc("PATCH /pessoas/{codigo}", pessoaHandler.Atualizar)
	mux.HandleFunc("PATCH /pessoas/{codigo}/responsavel", pessoaHandler.HabilitarComoResponsavel)

	mux.HandleFunc("POST /chamados", chamadoHandler.Cadastrar)
	mux.HandleFunc("GET /chamados/{id}", chamadoHandler.BuscarPorID)
	mux.HandleFunc("GET /chamados", chamadoHandler.Listar)
	mux.HandleFunc("PUT /chamados/{id}", chamadoHandler.Atualizar)

	fmt.Println("servidor iniciado em http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("erro ao iniciar servidor HTTP: %v\n", err)
	}

}
