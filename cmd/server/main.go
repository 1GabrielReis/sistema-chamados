package main

import (
	"context"
	"fmt"
	"os"

	"github.com/1GabrielReis/sistema-chamados/internal/database"
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

	_ = service.NovoPessoaService(pessoaRepository)
	_ = service.NovoChamadoService(chamadoRepository, pessoaRepository)

	fmt.Println("conexão com o PostgreSQL realizada com sucesso")
}
