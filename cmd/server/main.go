package main

import (
	"context"
	"fmt"
	"os"

	"github.com/1GabrielReis/sistema-chamados/internal/database"
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
	fmt.Println("conexão com o PostgreSQL realizada com sucesso")
}
